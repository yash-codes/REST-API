package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yash-codes/students-api/internal/config"
)

func main() {
	fmt.Println("Welcome to students-api")

	// load config
	cfg := config.MustLoad()
	fmt.Println("Configurations:", *cfg)

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

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server:%s\n", err.Error())
	}

	fmt.Println("Server Started at", cfg.Addr)
}
