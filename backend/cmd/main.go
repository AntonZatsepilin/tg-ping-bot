package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"goPingRobot/pkg/repository"
	"goPingRobot/pkg/service"
	"goPingRobot/pkg/telegram"
	"goPingRobot/pkg/workerpool"

	"github.com/gin-gonic/gin"
)

type App struct {
    performanceData []PerformanceEntry
    mu              sync.Mutex
    allTestsDone    bool // Флаг для отслеживания завершения всех тестов
}

type PerformanceEntry struct {
    Workers int     `json:"workers"`
    Time    float64 `json:"time"`
}

func main() {
    // Конфигурация MongoDB
    cfg := repository.Config{
        Host:     "mongodb",
        Port:     "27017",
        Username: "admin",
        Password: "sslowmm",
        DBname:   "mongodb",
    }

    // Подключение к MongoDB
    client, err := repository.NewMongoDB(cfg)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    defer client.Disconnect(context.TODO())

    // Инициализация репозитория
    db := client.Database(cfg.DBname)
    repo := repository.NewRepository(db)

    // Генерация 1000 ссылок
    repo.Generator.GenerateUrls(10)

    // Инициализация сервиса
    svc := service.NewGeneratorService(repo.Generator)

    // Настройка Telegram бота
    token := os.Getenv("TELEGRAM_BOT_TOKEN")
    if token == "" {
        log.Fatal("TELEGRAM_BOT_TOKEN is not set")
    }
    chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
    if chatIDStr == "" {
        log.Fatal("TELEGRAM_CHAT_ID is not set")
    }
    chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
    if err != nil {
        log.Fatalf("Error converting TELEGRAM_CHAT_ID to int64: %v", err)
    }
    telegram.Init(token, chatID)

    // Инициализация приложения
    app := App{}

    // Настройка HTTP-сервера
    r := gin.Default()

    // Роут для получения данных о производительности
    r.GET("/performance", func(c *gin.Context) {
        app.mu.Lock()
        defer app.mu.Unlock()

        if !app.allTestsDone {
            c.JSON(200, gin.H{"message": "Tests are still running. Please wait."})
            return
        }

        if len(app.performanceData) == 0 {
            c.JSON(200, gin.H{"message": "No data available"})
            return
        }

        c.JSON(200, app.performanceData)
    })

    // Запуск HTTP-сервера
    go func() {
        if err := r.Run(":8080"); err != nil {
            log.Fatalf("Failed to start HTTP server: %v", err)
        }
    }()

    // Тестирование производительности
    app.testPerformance(svc)

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
    <-quit
}

func (app *App) testPerformance(svc *service.GeneratorService) {
    workerCounts := []int{1, 2, 4, 8, 16} // Количество воркеров для тестирования

    for _, workersCount := range workerCounts {
        REQUEST_TIMEOUT, _ := strconv.Atoi(os.Getenv("REQUEST_TIMEOUT"))
        results := make(chan workerpool.Result)
        workerPool := workerpool.New(workersCount, time.Duration(REQUEST_TIMEOUT)*time.Second, results)
        workerPool.Init()

        // Запуск обработки результатов в отдельной горутине
        go processResults(results)

        // Замер времени выполнения
        startTime := time.Now()
        if err := svc.GenerateJobs(workerPool); err != nil {
            log.Printf("Error generating jobs: %v", err)
        }
        workerPool.Stop()
        duration := time.Since(startTime).Seconds()

        // Сохраняем данные о производительности
        app.mu.Lock()
        app.performanceData = append(app.performanceData, PerformanceEntry{
            Workers: workersCount,
            Time:    duration,
        })
        app.mu.Unlock()

        log.Printf("Workers: %d, Time taken: %.2f seconds\n", workersCount, duration)

        // Ждем завершения обработки результатов
        close(results)
    }

    // Устанавливаем флаг завершения всех тестов
    app.mu.Lock()
    app.allTestsDone = true
    app.mu.Unlock()

    log.Println("All performance tests completed.")
}

func processResults(results chan workerpool.Result) {
    for result := range results {
        info := result.Info()
        log.Println(info)
        if result.Error != nil {
            telegram.SendMessage(info) // Отправка ошибок в Telegram
        }
    }
}