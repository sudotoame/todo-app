package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"todo-app/internal/core/domains/task"
)

type Service interface {
	CreateTask(ctx context.Context, title, description string) error
	FoundTask(ctx context.Context, title string, complete *bool) ([]task.Task, error)
	SetTask(ctx context.Context, title string, completed bool) error
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

	t, err := h.Service.FoundTask(r.Context(), req.Title, nil)
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

func (h *Handlers) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	completed := query.Get("completed")
	var req CreateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	var completedPtr *bool // nil
	boolVal, err := strconv.ParseBool(completed)
	if err == nil { // if err != nil {completedPtr == nil}
		completedPtr = &boolVal // else {true or false}
	}

	tasks, err := h.Service.FoundTask(r.Context(), req.Title, completedPtr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		fmt.Println("error to encoding response")

		return
	}
}

func (h *Handlers) HandleSetTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	id := r.PathValue("id")

	if err := h.Service.SetTask(r.Context(), id, req.Completed); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	task, _ := h.Service.FoundTask(r.Context(), id, nil)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(task); err != nil {
		fmt.Println("Ошибка энкодинга при response")

		return
	}
}
