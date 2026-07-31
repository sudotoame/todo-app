package transport

import (
	"errors"
	"net/http"
)

type Server struct {
	Handler Handlers
	Server  *http.Server
}

func NewServer(handler Handlers) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks", handler.HandleNewTask)
	mux.HandleFunc("PATCH /tasks/{id}", handler.HandleSetCompleteTask)
	mux.HandleFunc("PUT /tasks/{id}", handler.HandleUpdateTask)
	mux.HandleFunc("GET /tasks", handler.HandleListTasks) // ?completed=true/false
	mux.HandleFunc("GET /tasks/{id}", handler.HandleGetTaskByID)
	mux.HandleFunc("DELETE /tasks/{id}", handler.HandleDeleteTask)

	return &Server{
		Server: &http.Server{
			Addr:    ":9091",
			Handler: mux,
		},
	}
}

func (s *Server) Run() error {
	if err := s.Server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		} else {
			return err
		}
	}

	return nil
}
