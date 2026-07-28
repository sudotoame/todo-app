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
