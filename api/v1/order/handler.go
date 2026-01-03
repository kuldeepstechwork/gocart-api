// Package order handles order related operations.
package order

import (
	"encoding/json"
	"net/http"
)

// Handler manages order requests.
type Handler struct{}

// NewHandler creates a new order handler.
func NewHandler() *Handler {
	return &Handler{}
}

// CreateOrder creates a new order.
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var o struct {
		CartID string `json:"cart_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("Order created v1"))
}

// GetOrder returns an order status.
func (h *Handler) GetOrder(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Get Order v1"))
}
