package service

import (
	"goPingRobot/pkg/repository"
	"goPingRobot/pkg/workerpool"
	"log"
)

type GeneratorService struct {
	repo repository.Generate
}

func NewGeneratorService(repo repository.Generate) *GeneratorService {
	return &GeneratorService{
		repo: repo,
	}
}

func (s *GeneratorService) GenerateJobs(wp *workerpool.Pool) error {
    urls, err := s.repo.GetUrls()
    if err != nil {
        log.Printf("Error fetching URLs from MongoDB: %v", err)
        return err
    }

    for _, url := range urls {
        wp.Push(workerpool.Job{URL: url})
    }
    return nil
}