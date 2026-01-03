package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kuldeepstechwork/gocart-api/internal/dto"
	"github.com/kuldeepstechwork/gocart-api/internal/utils"
	"go.uber.org/mock/gomock"
)

func TestProductHandler_Public(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ts := setupTestServer(ctrl)
	router := ts.Server.SetupRoutes()

	t.Run("GetProducts", func(t *testing.T) {
		ts.ProductService.EXPECT().GetProducts(gomock.Any(), gomock.Any()).Return([]dto.ProductResponse{}, &utils.PaginationMeta{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("GetProduct", func(t *testing.T) {
		ts.ProductService.EXPECT().GetProduct(uint(1)).Return(&dto.ProductResponse{ID: 1}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/products/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("SearchProducts", func(t *testing.T) {
		ts.ProductService.EXPECT().SearchProducts(gomock.Any()).Return([]dto.ProductSearchResult{}, &utils.PaginationMeta{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/search?q=test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("GetCategories", func(t *testing.T) {
		ts.ProductService.EXPECT().GetCategories().Return([]dto.CategoryResponse{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestProductHandler_Admin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ts := setupTestServer(ctrl)
	router := ts.Server.SetupRoutes()
	adminToken := createAdminToken(1)

	t.Run("CreateProduct_Success", func(t *testing.T) {
		reqBody := dto.CreateProductRequest{
			CategoryID: 1,
			Name:       "New Product",
			Price:      100.0,
			Stock:      10,
			SKU:        "TEST-SKU-1",
		}
		body, _ := json.Marshal(reqBody)

		ts.ProductService.EXPECT().CreateProduct(gomock.Any()).Return(&dto.ProductResponse{}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/products/", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("UpdateProduct_Success", func(t *testing.T) {
		reqBody := dto.UpdateProductRequest{
			CategoryID: 1,
			Name:       "Updated Product",
			Price:      150.0,
			Stock:      5,
		}
		body, _ := json.Marshal(reqBody)

		ts.ProductService.EXPECT().UpdateProduct(uint(1), gomock.Any()).Return(&dto.ProductResponse{}, nil)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/products/1", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("DeleteProduct_Success", func(t *testing.T) {
		ts.ProductService.EXPECT().DeleteProduct(uint(1)).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/products/1", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("CreateCategory_Success", func(t *testing.T) {
		reqBody := dto.CreateCategoryRequest{Name: "New Category"}
		body, _ := json.Marshal(reqBody)

		ts.ProductService.EXPECT().CreateCategory(gomock.Any()).Return(&dto.CategoryResponse{}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/categories/", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}
	})

	t.Run("UpdateCategory_Success", func(t *testing.T) {
		reqBody := dto.UpdateCategoryRequest{Name: "Updated Category"}
		body, _ := json.Marshal(reqBody)

		ts.ProductService.EXPECT().UpdateCategory(uint(1), gomock.Any()).Return(&dto.CategoryResponse{}, nil)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/categories/1", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("DeleteCategory_Success", func(t *testing.T) {
		ts.ProductService.EXPECT().DeleteCategory(uint(1)).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/categories/1", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("Forbidden_For_Customer", func(t *testing.T) {
		customerToken := createTestToken(2)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/products/", nil)
		req.Header.Set("Authorization", "Bearer "+customerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("expected status 403, got %d", w.Code)
		}
	})
}
