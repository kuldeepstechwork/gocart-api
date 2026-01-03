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

func TestUserHandler_Profile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ts := setupTestServer(ctrl)
	router := ts.Server.SetupRoutes()
	userID := uint(1)
	token := createTestToken(userID)

	t.Run("GetProfile_Success", func(t *testing.T) {
		ts.UserService.EXPECT().GetProfile(userID).Return(&dto.UserResponse{ID: userID}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("UpdateProfile_Success", func(t *testing.T) {
		reqBody := dto.UpdateProfileRequest{
			FirstName: "Jane",
			LastName:  "Smith",
		}
		body, _ := json.Marshal(reqBody)

		ts.UserService.EXPECT().UpdateProfile(userID, gomock.Any()).Return(&dto.UserResponse{}, nil)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("GetProfile_NotFound", func(t *testing.T) {
		ts.UserService.EXPECT().GetProfile(userID).Return(nil, errors.New("not found"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})
}
