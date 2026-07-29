package repository

import (
	"context"
	"sync"
	apperror "todo-app/internal/core/app-error"
	"todo-app/internal/core/domains/task"
)

type Repo struct {
	Tasks map[string]task.Task
	mtx   sync.Mutex
	// pool *pgx.Pool
}

func NewRepo() *Repo {
	return &Repo{
		Tasks: make(map[string]task.Task),
	}
}

func (r *Repo) AddTask(ctx context.Context, task task.Task) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if _, ok := r.Tasks[task.Title]; ok {
		return apperror.ErrTaskAlreadyExists
	}

	r.Tasks[task.Title] = task

	return nil
}

func (r *Repo) GetTask(ctx context.Context, title string, completed *bool) ([]task.Task, error) {
	if completed != nil {
		var tasks []task.Task
		for _, value := range r.Tasks {
			if value.Completed == *completed {
				tasks = append(tasks, value)
			}
		}

		if len(tasks) == 0 {
			if *completed {
				return nil, apperror.ErrNotFoundCompletedTask
			}
			return nil, apperror.ErrNotFoundUncompletedTask
		}

		return tasks, nil
	}

	var value []task.Task
	if title == "" {
		for _, v := range r.Tasks {
			value = append(value, v)
		}

		return value, nil
	}

	v, ok := r.Tasks[title]
	if !ok {
		return []task.Task{}, apperror.ErrTaskNotFound
	}
	value = append(value, v)

	return value, nil
}
