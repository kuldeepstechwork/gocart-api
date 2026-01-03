package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kuldeepstechwork/gocart-api/internal/dto"
	"go.uber.org/mock/gomock"
)

func TestAuthHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ts := setupTestServer(ctrl)
	router := ts.Server.SetupRoutes()

	t.Run("Success", func(t *testing.T) {
		reqBody := dto.RegisterRequest{
			Email:     "test@example.com",
			Password:  "password123",
			FirstName: "John",
			LastName:  "Doe",
		}
		body, _ := json.Marshal(reqBody)

		ts.AuthService.EXPECT().Register(gomock.Any()).Return(&dto.AuthResponse{
			AccessToken: "access",
		}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("InvalidData", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString("invalid"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("ServiceError", func(t *testing.T) {
		reqBody := dto.RegisterRequest{
			Email:     "error@example.com",
			Password:  "password123",
			FirstName: "John",
			LastName:  "Doe",
		}
		body, _ := json.Marshal(reqBody)

		ts.AuthService.EXPECT().Register(gomock.Any()).Return(nil, errors.New("fail"))

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestAuthHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ts := setupTestServer(ctrl)
	router := ts.Server.SetupRoutes()

	t.Run("Success", func(t *testing.T) {
		reqBody := dto.LoginRequest{Email: "test@example.com", Password: "password123"}
		body, _ := json.Marshal(reqBody)

		ts.AuthService.EXPECT().Login(gomock.Any()).Return(&dto.AuthResponse{AccessToken: "access"}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("Unauthorized", func(t *testing.T) {
		reqBody := dto.LoginRequest{Email: "wrong@example.com", Password: "password123"}
		body, _ := json.Marshal(reqBody)

		ts.AuthService.EXPECT().Login(gomock.Any()).Return(nil, errors.New("unauthorized"))

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ts := setupTestServer(ctrl)
	router := ts.Server.SetupRoutes()

	t.Run("Success", func(t *testing.T) {
		reqBody := dto.RefreshTokenRequest{RefreshToken: "valid"}
		body, _ := json.Marshal(reqBody)

		ts.AuthService.EXPECT().RefreshToken(gomock.Any()).Return(&dto.AuthResponse{AccessToken: "new"}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ts := setupTestServer(ctrl)
	router := ts.Server.SetupRoutes()

	t.Run("Success", func(t *testing.T) {
		reqBody := dto.RefreshTokenRequest{RefreshToken: "valid"}
		body, _ := json.Marshal(reqBody)

		ts.AuthService.EXPECT().Logout("valid").Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}
