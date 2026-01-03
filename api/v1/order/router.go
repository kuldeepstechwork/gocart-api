package order

import (
	"net/http"
)

// RegisterRoutes registers order routes.
func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("POST /api/v1/order", h.CreateOrder)
	mux.HandleFunc("GET /api/v1/order/status", h.GetOrder)
}
