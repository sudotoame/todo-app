package service

import (
	"context"
	"todo-app/internal/core/domains/task"
)

type Repository interface {
	AddTask(ctx context.Context, task task.Task) error
	GetTask(ctx context.Context, title string) (task.Task, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateTask(ctx context.Context, title, description string) error {
	task, err := task.NewTask(title, description)
	if err != nil {
		return err
	}

	if err := s.repo.AddTask(ctx, task); err != nil {
		return err
	}

	return nil
}

func (s *Service) FoundTask(ctx context.Context, title string) (task.Task, error) {
	t, err := s.repo.GetTask(ctx, title)
	if err != nil {
		return task.Task{}, err
	}
	return t, nil
}
