package service

import (
	"context"
	"todo-app/internal/core/domains/task"
)

type Repository interface {
	AddTask(ctx context.Context, t task.Task) (*task.Task, error)
	GetTaskByID(ctx context.Context, id int) (*task.Task, error)
	ListTasks(ctx context.Context) ([]task.Task, error)
	ListTasksByComplete(ctx context.Context, completed *bool) ([]task.Task, error)
	SetCompleteTask(ctx context.Context, id int, t task.Task) error
	UpdateTask(ctx context.Context, id int, t task.Task) error
	DeleteTask(ctx context.Context, id int) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateTask(ctx context.Context, title, description string) (*task.Task, error) {
	t, err := task.NewTask(title, description)
	if err != nil {
		return nil, err
	}

	tr, err := s.repo.AddTask(ctx, t)
	if err != nil {
		return nil, err
	}

	return tr, nil
}

func (s *Service) GetTaskByID(ctx context.Context, id int) (*task.Task, error) {
	t, err := s.repo.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) ListTasks(ctx context.Context, completed *bool) ([]task.Task, error) {
	if completed != nil {
		task, err := s.repo.ListTasksByComplete(ctx, completed)
		if err != nil {
			return nil, err
		}
		return task, nil
	}

	task, err := s.repo.ListTasks(ctx)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Service) SetCompleteTask(ctx context.Context, id int, t task.Task) error {
	return s.repo.SetCompleteTask(ctx, id, t)
}

func (s *Service) UpdateTask(ctx context.Context, id int, t task.Task) error {
	return s.repo.UpdateTask(ctx, id, t)
}

func (s *Service) DeleteTask(ctx context.Context, id int) error {
	err := s.repo.DeleteTask(ctx, id)

	return err
}
