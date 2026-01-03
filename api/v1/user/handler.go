package user

import (
    "encoding/json"
    "net/http"
)

type Handler struct{}

func NewHandler() *Handler {
    return &Handler{}
}

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
    w.Write([]byte("User registered v1"))
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("User Profile v1"))
}
