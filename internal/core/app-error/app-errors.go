package apperror

import "errors"

var ErrEmptyTitle = errors.New("empty title")
var ErrEmptyDescription = errors.New("empty description")
var ErrTaskAlreadyExists = errors.New("task already exists")
var ErrTaskNotFound = errors.New("task not found")
var ErrNotFoundCompletedTask = errors.New("completed task not found")
var ErrNotFoundUncompletedTask = errors.New("uncompleted task not found")
