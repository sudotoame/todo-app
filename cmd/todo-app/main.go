package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo-app/internal/features/task/repository"
	"todo-app/internal/features/task/service"
	"todo-app/internal/features/task/transport"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	pool, err := pgxpool.New(context.Background(), os.Getenv("PGX_CONN_DOCKER"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("pgx connected")

	repo := repository.NewPostgresRepo(pool)
	service := service.NewService(repo)
	handler := transport.NewHandler(service)
	server := transport.NewServer(*handler)

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Run()
	}()

	select {
	case <-ctx.Done():
		fmt.Println("ctx cancelled")
	case err := <-errCh:
		if err != nil {
			fmt.Println(err)
		}
	}

	shutdownCtx, shutdownCancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancle()
	if err := server.Server.Shutdown(shutdownCtx); err != nil {
		fmt.Println("Graceful shutdown failed")
	} else {
		fmt.Println("Graceful shutdown!")
	}
}
