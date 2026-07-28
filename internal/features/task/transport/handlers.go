package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"todo-app/internal/core/domains/task"
)

type Service interface {
	CreateTask(ctx context.Context, title, description string) error
	FoundTask(ctx context.Context, title string) (task.Task, error)
}

type Handlers struct {
	Service Service
}

func NewHandler(service Service) *Handlers {
	return &Handlers{
		Service: service,
	}
}

func (h *Handlers) HandleNewTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error to decoding request", http.StatusBadRequest)
		// TODO: добавить логгер
		return
	}

	if err := h.Service.CreateTask(r.Context(), req.Title, req.Description); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		// TODO: добавить логгер
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	t, err := h.Service.FoundTask(r.Context(), req.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		// TODO: добавить логгер
		return
	}

	if err := json.NewEncoder(w).Encode(t); err != nil {
		fmt.Println("Encoding error!")

		return
	}
}
