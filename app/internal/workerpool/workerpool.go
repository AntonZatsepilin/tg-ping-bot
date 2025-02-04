package workerpool

import (
	"github.com/AntonZatsepilin/mephi-database-homework/app/internal/service"
	"github.com/sirupsen/logrus"
)

type url_response struct {
	URL    string
	Status string
}

func worker(id int, urlChan chan string, resultChan chan url_response) {
	for url := range urlChan {
		logrus.Infoln("Worker ", id, " took ", url, " as work")
		result := url_response{url, service.CheckURL(url)}
		resultChan <- result
	}
}

func WorkerPool(workerCount int, tasks []string) map[string]string {
	taskChan := make(chan string, len(tasks))
	resultChan := make(chan url_response, len(tasks))

	for i := 0; i < workerCount; i++ {
		go worker(i, taskChan, resultChan)
		logrus.Infoln("Made ", i, " worker")
	}

	for i := 0; i < len(tasks); i++ {
		taskChan <- tasks[i]
		logrus.Infoln("Put ", i, " task")
	}
	close(taskChan)
	logrus.Infoln("Task chan closed")

	statuses := make(map[string]string)
	for i := 0; i < len(tasks); i++ {
		resp := <-resultChan
		statuses[resp.URL] = resp.Status
	}
	close(resultChan)
	logrus.Infoln("Result chan closed")

	return statuses
}
