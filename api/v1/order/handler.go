package order

import (
    "encoding/json"
    "net/http"
)

type Handler struct{}

func NewHandler() *Handler {
    return &Handler{}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    var o struct {
        CartID string `json:"cart_id"`
    }
    if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Order created v1"))
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Get Order v1"))
}
