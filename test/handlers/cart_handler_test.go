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

func TestCartHandler_GetCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ts := setupTestServer(ctrl)
	router := ts.Server.SetupRoutes()
	userID := uint(1)
	token := createTestToken(userID)

	t.Run("Success", func(t *testing.T) {
		ts.CartService.EXPECT().GetCart(userID).Return(&dto.CartResponse{UserID: userID}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/cart/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		ts.CartService.EXPECT().GetCart(userID).Return(nil, errors.New("not found"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/cart/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})
}

func TestCartHandler_AddToCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ts := setupTestServer(ctrl)
	router := ts.Server.SetupRoutes()
	userID := uint(1)
	token := createTestToken(userID)

	t.Run("Success", func(t *testing.T) {
		reqBody := dto.AddToCartRequest{ProductID: 1, Quantity: 2}
		body, _ := json.Marshal(reqBody)

		ts.CartService.EXPECT().AddToCart(userID, gomock.Any()).Return(&dto.CartResponse{}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/cart/items", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})
}

func TestCartHandler_UpdateCartItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ts := setupTestServer(ctrl)
	router := ts.Server.SetupRoutes()
	userID := uint(1)
	token := createTestToken(userID)

	t.Run("Success", func(t *testing.T) {
		reqBody := dto.UpdateCartItemRequest{Quantity: 5}
		body, _ := json.Marshal(reqBody)

		ts.CartService.EXPECT().UpdateCartItem(userID, uint(10), gomock.Any()).Return(&dto.CartResponse{}, nil)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/cart/items/10", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("InvalidID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/cart/items/abc", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestCartHandler_RemoveFromCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ts := setupTestServer(ctrl)
	router := ts.Server.SetupRoutes()
	userID := uint(1)
	token := createTestToken(userID)

	t.Run("Success", func(t *testing.T) {
		ts.CartService.EXPECT().RemoveFromCart(userID, uint(10)).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/cart/items/10", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}
