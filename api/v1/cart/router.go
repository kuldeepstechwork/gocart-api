package cart

import (
    "net/http"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
    mux.HandleFunc("GET /api/v1/cart", h.GetCart)
    mux.HandleFunc("POST /api/v1/cart/items", h.AddItem)
}
