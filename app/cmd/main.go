package main

import (
	"fmt"
	"time"

	"github.com/AntonZatsepilin/mephi-database-homework/app/internal/service"
	"github.com/sirupsen/logrus"
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

func MeasureExecutionTime(name string, f func(), log *logrus.Logger) {
	start := time.Now()
	f()
	duration := time.Since(start)
	log.Info("Function ", name, " took ", duration, " to execute\n")
}

func main() {
	log := logrus.New()

	log.SetFormatter(&logrus.TextFormatter{})

	ressult := service.CheckURLs(urls)
	for _, n := range ressult {
		fmt.Println(n)
	}
}
