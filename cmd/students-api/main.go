package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/z1shivam/learning-go/internal/config"
	"github.com/z1shivam/learning-go/internal/http/handlers/student"
	"github.com/z1shivam/learning-go/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal("Database not connected")
	}

	slog.Info("Storage Initialized.", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage))

	// setup server
	server := http.Server{
		Addr:    cfg.HttpServer.Addr,
		Handler: router,
	}

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("Server started at: ", slog.String("address", cfg.HttpServer.Addr))
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server.")
		}
	}()

	<-done

	slog.Info("Shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown the server.", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully!")
}
