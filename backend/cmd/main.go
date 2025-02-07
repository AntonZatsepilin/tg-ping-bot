package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"goPingRobot/pkg/repository"
	"goPingRobot/pkg/service"
	"goPingRobot/pkg/telegram"
	"goPingRobot/pkg/workerpool"

	"github.com/wcharczuk/go-chart/v2"
)

type App struct {
    performanceData []PerformanceEntry
    mu              sync.Mutex
    allTestsDone    bool
}

type PerformanceEntry struct {
    Workers int     `json:"workers"`
    Time    float64 `json:"time"`
}

func main() {
    cfg := repository.Config{
        Host:     "mongodb",
        Port:     "27017",
        Username: "admin",
        Password: "sslowmm",
        DBname:   "mongodb",
    }

    client, err := repository.NewMongoDB(cfg)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    defer client.Disconnect(context.TODO())

    db := client.Database(cfg.DBname)
    repo := repository.NewRepository(db)

    repo.Generator.GenerateUrls(10)

    svc := service.NewGeneratorService(repo.Generator)

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

    app := App{}

    app.testPerformance(svc)

    app.openChart()

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
    <-quit
}

func (app *App) testPerformance(svc *service.GeneratorService) {
    workerCounts := []int{1, 2, 4, 8, 16}

    for _, workersCount := range workerCounts {
        REQUEST_TIMEOUT, _ := strconv.Atoi(os.Getenv("REQUEST_TIMEOUT"))
        results := make(chan workerpool.Result)
        workerPool := workerpool.New(workersCount, time.Duration(REQUEST_TIMEOUT)*time.Second, results)
        workerPool.Init()

        go processResults(results)

        startTime := time.Now()
        if err := svc.GenerateJobs(workerPool); err != nil {
            log.Printf("Error generating jobs: %v", err)
        }
        workerPool.Stop()
        duration := time.Since(startTime).Seconds()

        app.mu.Lock()
        app.performanceData = append(app.performanceData, PerformanceEntry{
            Workers: workersCount,
            Time:    duration,
        })
        app.mu.Unlock()

        log.Printf("Workers: %d, Time taken: %.2f seconds\n", workersCount, duration)

        close(results)
    }

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
            telegram.SendMessage(info)
        }
    }
}


func generateChart(data []PerformanceEntry) ([]byte, error) {
    var xValues []float64
    var yValues []float64

    for _, entry := range data {
        xValues = append(xValues, float64(entry.Workers))
        yValues = append(yValues, entry.Time)
    }

    graph := chart.Chart{
        Title: "Test",
        XAxis: chart.XAxis{
            Name: "Number of Workers",
            Style: chart.Style{
                StrokeColor: chart.ColorBlack,
                FontSize:    12,           
            },
        },
        YAxis: chart.YAxis{
            Name: "Time Taken (seconds)",
            Style: chart.Style{
                StrokeColor: chart.ColorBlack, 
                FontSize:    12,            
            },
        },
        Series: []chart.Series{
            chart.ContinuousSeries{
                Name:    "Time vs Workers",
                XValues: xValues,
                YValues: yValues,
                Style: chart.Style{
                    StrokeColor: chart.ColorBlue, 
                    FillColor:   chart.ColorBlue.WithAlpha(64), 
                },
            },
        },
    }

    buffer := bytes.NewBuffer([]byte{})
    err := graph.Render(chart.PNG, buffer)
    if err != nil {
        return nil, err
    }

    return buffer.Bytes(), nil
}

func (app *App) openChart() {
    app.mu.Lock()
    defer app.mu.Unlock()

    if len(app.performanceData) == 0 {
        log.Println("No data available for chart")
        return
    }

    graphBytes, err := generateChart(app.performanceData)
    if err != nil {
        log.Printf("Failed to generate chart: %v", err)
        return
    }

    filePath := "/app/data/performance_chart.png" 
    err = os.WriteFile(filePath, graphBytes, 0644)
    if err != nil {
        log.Printf("Failed to save chart to file: %v", err)
        return
    }

    log.Printf("Chart saved to %s", filePath)

    if _, err := exec.LookPath("xdg-open"); err == nil {
        err = openFileInBrowser(filePath)
        if err != nil {
            log.Printf("Failed to open chart in browser: %v", err)
        }
    } else {
        log.Printf("To view the chart, open the file manually: ./app_data/performance_chart.png")
    }
}

func openFileInBrowser(filePath string) error {
    var cmd *exec.Cmd

    switch os := os.Getenv("GOOS"); os {
    case "darwin": // macOS
        cmd = exec.Command("open", filePath)
    case "windows":
        cmd = exec.Command("cmd", "/c", "start", filePath)
    default: 
        cmd = exec.Command("xdg-open", filePath)
    }

    return cmd.Run()
}