// Package main запускает HTTP-сервис назначения ревьюеров для Pull Request'ов.
package main

import (
	"log"
	"net/http"
)

// main - точка входа в сервис назначения ревьюеров.
func main() {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("starting pr-reviewer-service on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to start http server: %v", err)
	}
}
