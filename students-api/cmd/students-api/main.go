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

	"github.com/yash-codes/students-api/internal/config"
)

func main() {
	slog.Info("Welcome to students-api")

	// load config
	cfg := config.MustLoad()
	slog.Info("Configurations parsing completed", "configurations", *cfg)

	// TODO: database setup

	// setup router
	router := http.NewServeMux()

	// declare and defile handler funnction corrosponding to the url path "/"
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to students-api"))
	})

	// setup server
	// config the server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// Start the server
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start server: %s\n", err.Error())
		}
		slog.Info("Server Started at", "Address", cfg.Addr)
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server glacefully", "error", err)
	}

	slog.Info("Server shutdown successfully!")
}
