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
	"github.com/yash-codes/students-api/internal/http/handlers/student"
	"github.com/yash-codes/students-api/internal/storage/sqlite"
)

func main() {
	slog.Info("Welcome to students-api")

	// load config
	cfg := config.MustLoad()
	slog.Info("Configurations parsing completed", "configurations", *cfg)

	// TODO: database setup
	//storage, err := sqlite.New(cfg)
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Storage initialized", "env", cfg.Env, "version", "1.0.0")

	// setup router
	router := http.NewServeMux()

	// declare and defile handler funnction corrosponding to the url path "/"
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))
	router.HandleFunc("DELETE /api/students/{id}", student.DeleteStudent(storage))

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
