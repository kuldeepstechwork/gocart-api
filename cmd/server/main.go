// Package main is the entry point for the server.
package main

import (
	"log"
	"net/http"

	"github.com/kuldeepstechwork/gocart-api/api/v1/cart"
	"github.com/kuldeepstechwork/gocart-api/api/v1/order"
	"github.com/kuldeepstechwork/gocart-api/api/v1/user"
	"github.com/kuldeepstechwork/gocart-api/internal/auth"
)

func main() {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Initialize handlers
	cartHandler := cart.NewHandler()
	userHandler := user.NewHandler()
	orderHandler := order.NewHandler()

	// Register routes
	cart.RegisterRoutes(mux, cartHandler)
	user.RegisterRoutes(mux, userHandler)
	order.RegisterRoutes(mux, orderHandler)

	// Wrap with middleware
	handler := auth.Middleware(mux)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
