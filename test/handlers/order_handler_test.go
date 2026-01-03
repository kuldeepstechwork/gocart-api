package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kuldeepstechwork/gocart-api/internal/dto"
	"github.com/kuldeepstechwork/gocart-api/internal/utils"
	"go.uber.org/mock/gomock"
)

func TestOrderHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ts := setupTestServer(ctrl)
	router := ts.Server.SetupRoutes()
	userID := uint(1)
	token := createTestToken(userID)

	t.Run("CreateOrder_Success", func(t *testing.T) {
		ts.OrderService.EXPECT().CreateOrder(userID).Return(&dto.OrderResponse{ID: 1}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}
	})

	t.Run("GetOrders_Success", func(t *testing.T) {
		ts.OrderService.EXPECT().GetOrders(userID, 1, 10).Return([]dto.OrderResponse{}, &utils.PaginationMeta{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/?page=1&limit=10", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("GetOrder_Success", func(t *testing.T) {
		ts.OrderService.EXPECT().GetOrder(userID, uint(100)).Return(&dto.OrderResponse{ID: 100}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/100", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("GetOrder_NotFound", func(t *testing.T) {
		ts.OrderService.EXPECT().GetOrder(userID, uint(999)).Return(nil, errors.New("not found"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/999", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})
}
