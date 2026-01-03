// Package handlers provides HTTP handlers.
package handlers

import (
	"encoding/json"
	"net/http"
)

// CartItem represents an item in the cart.
type CartItem struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// AddItem handles adding an item to the cart.
func AddItem(w http.ResponseWriter, r *http.Request) {
	var item CartItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	// Placeholder: In a real app, you'd store the item.
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(item)
}
