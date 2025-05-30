package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/z1shivam/learning-go/internal/config"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// database setup
	// setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to students api"))
	})

	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	fmt.Printf("Server started at port: %s", cfg.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to start server.")
	}
}
