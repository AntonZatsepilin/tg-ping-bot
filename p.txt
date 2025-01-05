package main

import (
	"fmt"
	"goPingRobot/telegram"
	"goPingRobot/workerpool"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"gopkg.in/yaml.v2"
)

func main() {

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

	results := make(chan workerpool.Result)

	WORKERS_COUNT, err := strconv.Atoi(os.Getenv("WORKERS_COUNT"))
	if err != nil {
		log.Fatal("WORKERS_COUNT is not set")
	}

	REQUEST_TIMEOUT, err := strconv.Atoi(os.Getenv("REQUEST_TIMEOUT"))
	if err != nil {
		log.Fatal("REQUEST_TIMEOUT is not set")
	}

	workerPool := workerpool.New(WORKERS_COUNT, time.Duration(REQUEST_TIMEOUT)*time.Second, results)

	workerPool.Init()

	go generateJobs(workerPool)
	go proccessResults(results)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	workerPool.Stop()
}

func proccessResults(results chan workerpool.Result) {
	go func() {
		for result := range results {
			info := result.Info()
			fmt.Println(info)

			if result.Error != nil {
				log.Printf("Error result: %v", result.Error)
				telegram.SendMessage(info)
			}
		}
	}()
}

func generateJobs(wp *workerpool.Pool) {

	links := make(map[string][]string)

	yamlFile, err := os.ReadFile("links.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &links)
	if err != nil {
		panic(err)
	}

	for {
		for _, url := range links["links"] {
			wp.Push(workerpool.Job{URL: url})
		}

		INTERVAL, err := strconv.Atoi(os.Getenv("INTERVAL"))
		if err != nil {
			log.Fatal("INTERVAL is not set")
		}

		time.Sleep(time.Duration(INTERVAL) * time.Second)
	}

}
