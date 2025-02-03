package service

import (
	"net/http"
	"time"
)

func CheckURLs(urls []string) map[string]string {
	statuses := make(map[string]string)
	results := make(chan struct {
		URL    string
		Status string
	})

	for _, url := range urls {
		go func(url string) {
			client := http.Client{Timeout: 5 * time.Second}
			resp, err := client.Get(url)
			if err != nil {
				results <- struct {
					URL    string
					Status string
				}{url, "unreachable"}
				return
			}
			results <- struct {
				URL    string
				Status string
			}{url, resp.Status}
		}(url)
	}

	for range urls {
		result := <-results
		statuses[result.URL] = result.Status
	}

	return statuses
}

func CheckURL(url string) string {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "unreachable"
	}
	return resp.Status
}
