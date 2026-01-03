// Package cart handles cart related operations.
package cart

import (
	"encoding/json"
	"net/http"
)

// Handler manages cart requests.
type Handler struct{}

// NewHandler creates a new cart handler.
func NewHandler() *Handler {
	return &Handler{}
}

// GetCart returns the current cart.
func (h *Handler) GetCart(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Get Cart v1"))
}

// AddItem adds an item to the cart.
func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	var item struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("Item added to cart v1"))
}
