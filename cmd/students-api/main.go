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

	"github.com/royalthapa/students-api/internal/config"
	"github.com/royalthapa/students-api/internal/http/handlers/student"
)

func main() {
	// load config
	cfg := config.MustLoad()
	// database setup

	// setup router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New())

	// setup server
	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	slog.Info("server started", slog.String("address", cfg.HTTPServer.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}
	}()

	<-done

	slog.Info("shitting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("Failed to shutdown server", slog.String("Error:", err.Error()))
	}

	slog.Info("server shutdown successfully")

}
