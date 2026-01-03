package user

import (
	"net/http"
)

// RegisterRoutes registers user routes.
func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("POST /api/v1/user/register", h.Register)
	mux.HandleFunc("GET /api/v1/user/profile", h.Profile)
}
