package main

import (
	"fmt"
	"time"

	"github.com/AntonZatsepilin/mephi-database-homework/app/internal/service"
	"github.com/AntonZatsepilin/mephi-database-homework/app/internal/workerpool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var urls = []string{
	"https://www.google.com",
	"https://www.github.com",
	"https://www.stackoverflow.com",
	"https://www.microsoft.com",
	"https://www.apple.com",
	"https://www.amazon.com",
	"https://www.reddit.com",
	"https://www.youtube.com",
	"https://www.linkedin.com",
	"https://www.twitter.com",
	"https://www.instagram.com",
	"https://www.facebook.com",
	"https://www.wikipedia.org",
	"https://www.twitch.tv",
	"https://www.netflix.com",
	"https://www.spotify.com",
	"https://www.dropbox.com",
	"https://www.heroku.com",
	"https://www.medium.com",
	"https://www.cloudflare.com",
	"https://www.digitalocean.com",
	"https://www.nytimes.com",
	"https://www.bbc.com",
	"https://www.cnn.com",
	"https://www.weather.com",
}

func MeasureExecutionTime1(name string, f func([]string) map[string]string, log *logrus.Logger) {
	start := time.Now()
	ressult := f(urls)
	for url, status := range ressult {
		fmt.Println(url, " ", status)
	}
	duration := time.Since(start)
	log.Info("Function ", name, " took ", duration, " to execute\n")
}

func MeasureExecutionTime(name string, f func(int, []string) map[string]string, log *logrus.Logger) {
	start := time.Now()

	workerCount := viper.GetInt("worker_count")
	ressult := f(workerCount, urls)
	for url, status := range ressult {
		fmt.Println(url, " ", status)
	}
	duration := time.Since(start)
	log.Info("Function ", name, " took ", duration, " to execute\n")
}

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{})
	log.SetLevel(logrus.InfoLevel)

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("../configs")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file: ", err)
	}

	//workerCount := viper.GetInt("worker_count")

	MeasureExecutionTime1("Without workerpool", service.CheckURLs, log)
	MeasureExecutionTime("With workerpool", workerpool.WorkerPool, log)
}
