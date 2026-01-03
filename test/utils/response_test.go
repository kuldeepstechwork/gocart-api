package utils_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kuldeepstechwork/gocart-api/internal/utils"
)

func TestResponseUtils(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("SuccessResponse", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		utils.SuccessResponse(c, "success", "data")

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var res utils.Response
		json.Unmarshal(w.Body.Bytes(), &res)
		if !res.Success || res.Message != "success" || res.Data != "data" {
			t.Errorf("unexpected response: %+v", res)
		}
	})

	t.Run("CreatedResponse", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		utils.CreatedResponse(c, "created", 123)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}
	})

	t.Run("ErrorResponse", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		utils.ErrorResponse(c, http.StatusInternalServerError, "error msg", errors.New("actual err"))

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}

		var res utils.Response
		json.Unmarshal(w.Body.Bytes(), &res)
		if res.Success || res.Message != "error msg" || res.Error != "actual err" {
			t.Errorf("unexpected response: %+v", res)
		}
	})

	t.Run("ErrorResponse_NilErr", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		utils.ErrorResponse(c, http.StatusInternalServerError, "error msg", nil)

		var res utils.Response
		json.Unmarshal(w.Body.Bytes(), &res)
		if res.Error != "" {
			t.Errorf("expected empty error string, got %s", res.Error)
		}
	})

	t.Run("BadRequestResponse", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		utils.BadRequestResponse(c, "bad data", nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("UnauthorizedResponse", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		utils.UnauthorizedResponse(c, "unauthorized")
		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("ForbiddenResponse", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		utils.ForbiddenResponse(c, "forbidden")
		if w.Code != http.StatusForbidden {
			t.Errorf("expected status 403, got %d", w.Code)
		}
	})

	t.Run("NotFoundResponse", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		utils.NotFoundResponse(c, "not found")
		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})

	t.Run("InternalServerErrorResponse", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		utils.InternalServerErrorResponse(c, "server error", nil)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})

	t.Run("PaginatedSuccessResponse", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		meta := utils.PaginationMeta{Page: 1, Limit: 10, Total: 100, TotalPages: 10}
		utils.PaginatedSuccessResponse(c, "success", []string{"a", "b"}, meta)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var res utils.PaginatedResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		if res.Meta.Total != 100 {
			t.Errorf("unexpected meta: %+v", res.Meta)
		}
	})
}
