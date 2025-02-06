package service

import (
	"goPingRobot/pkg/repository"
	"goPingRobot/pkg/workerpool"
)

type Generate interface {
	GenerateJobs(wp *workerpool.Pool) error
}

type Service struct {
    Generate
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Generate: NewGeneratorService(repos.Generator),
	}
}