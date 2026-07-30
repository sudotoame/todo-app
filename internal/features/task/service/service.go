package service

import (
	"context"
	"todo-app/internal/core/domains/task"
)

type Repository interface {
	AddTask(ctx context.Context, task task.Task) error
	GetTask(ctx context.Context, title string, completed *bool) ([]task.Task, error)
	SetCompletedTask(ctx context.Context, title string, completed bool) error
	DeleteTask(ctx context.Context, title string) error
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

func (s *Service) FoundTask(ctx context.Context, title string, complete *bool) ([]task.Task, error) {
	t, err := s.repo.GetTask(ctx, title, complete)
	if err != nil {
		return []task.Task{}, err
	}
	return t, nil
}

func (s *Service) SetTask(ctx context.Context, title string, completed bool) error {
	err := s.repo.SetCompletedTask(ctx, title, completed)

	return err
}

func (s *Service) DeleteTask(ctx context.Context, title string) error {
	err := s.repo.DeleteTask(ctx, title)

	return err
}
