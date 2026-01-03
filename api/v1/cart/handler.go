package cart

import (
    "encoding/json"
    "net/http"
)

type Handler struct{}

func NewHandler() *Handler {
    return &Handler{}
}

func (h *Handler) GetCart(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Get Cart v1"))
}

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
    w.Write([]byte("Item added to cart v1"))
}
