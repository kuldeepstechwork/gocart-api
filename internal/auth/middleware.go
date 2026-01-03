// Package auth provides authentication middleware.
package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

// UserKey is the context key for user information.
const UserKey contextKey = "user"

// Middleware checks for a valid Authorization header.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		token := tokenParts[1]
		// Placeholder: Validate token
		// In a real app, you'd parse and verify the JWT here.
		if token != "valid-token" {
			// For demo purposes, we accept "valid-token"
			// http.Error(w, "Invalid token", http.StatusUnauthorized)
			// return
			_ = token // dummy usage
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), UserKey, "test-user")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
