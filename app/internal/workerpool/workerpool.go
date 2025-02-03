package workerpool

import "github.com/AntonZatsepilin/mephi-database-homework/app/internal/service"

type URL_response struct {
	URL    string
	Status string
}

func worker(urlChan chan string, resultChan chan URL_response) {
	for url := range urlChan {
		resultChan <- URL_response{url, service.CheckURL(url)}
	}
}

func workerPool(workerCount int, tasks []string) map[string]string {
	taskChan := make(chan string, len(tasks))
	resultChan := make(chan URL_response, len(tasks))

	for i := 0; i < workerCount; i++ {
		go worker(taskChan, resultChan)
	}

	for i := 0; i < len(tasks); i++ {
		taskChan <- tasks[i]
	}
	close(taskChan)
	statuses := make(map[string]string)
	for resp := range resultChan {
		statuses[resp.URL] = resp.Status
	}
	return statuses
}
