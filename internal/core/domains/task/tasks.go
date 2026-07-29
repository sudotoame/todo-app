package task

import (
	"time"
	apperror "todo-app/internal/core/app-error"
)

type Task struct {
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}

func NewTask(title string, description string) (Task, error) {
	if title == "" {
		return Task{}, apperror.ErrEmptyTitle
	}
	if description == "" {
		return Task{}, apperror.ErrEmptyDescription
	}

	return Task{
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}, nil
}

func (t *Task) Complete() {
	completeTime := time.Now()

	t.Completed = true
	t.CompletedAt = &completeTime
}

func (t *Task) Uncomplete() {
	t.Completed = false
	t.CompletedAt = nil
}
