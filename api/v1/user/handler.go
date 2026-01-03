// Package user handles user related operations.
package user

import (
	"encoding/json"
	"net/http"
)

// Handler manages user requests.
type Handler struct{}

// NewHandler creates a new user handler.
func NewHandler() *Handler {
	return &Handler{}
}

// Register registers a new user.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var u struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("User registered v1"))
}

// Profile returns the user profile.
func (h *Handler) Profile(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("User Profile v1"))
}
