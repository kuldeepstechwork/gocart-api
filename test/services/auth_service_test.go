package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/kuldeepstechwork/gocart-api/internal/config"
	"github.com/kuldeepstechwork/gocart-api/internal/dto"
	"github.com/kuldeepstechwork/gocart-api/internal/models"
	"github.com/kuldeepstechwork/gocart-api/internal/services"
	"github.com/kuldeepstechwork/gocart-api/internal/utils"
	"github.com/kuldeepstechwork/gocart-api/test/mocks"
	"go.uber.org/mock/gomock"
)

func TestAuthService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockCartRepo := mocks.NewMockCartRepositoryInterface(ctrl)
	mockPublisher := mocks.NewMockPublisher(ctrl)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:              "test_secret",
			ExpiresIn:           time.Hour,
			RefreshTokenExpires: time.Hour * 24,
		},
	}

	authService := services.NewAuthService(cfg, mockPublisher, mockUserRepo, mockCartRepo)

	req := &dto.RegisterRequest{
		Email:     "new@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	t.Run("Success", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmail(req.Email).Return(nil, errors.New("not found"))
		mockUserRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(u *models.User) error {
			u.ID = 1
			return nil
		})
		mockCartRepo.EXPECT().Create(gomock.Any()).Return(nil)
		mockUserRepo.EXPECT().CreateRefreshToken(gomock.Any()).Return(nil)
		mockPublisher.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		resp, err := authService.Register(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.User.Email != req.Email {
			t.Errorf("expected email %s, got %s", req.Email, resp.User.Email)
		}
	})

	t.Run("UserAlreadyExists", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmail(req.Email).Return(&models.User{}, nil)
		_, err := authService.Register(req)
		if err == nil || err.Error() != "you cannot register with this email" {
			t.Errorf("expected 'you cannot register with this email' error, got %v", err)
		}
	})

	t.Run("UserCreationFailure", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmail(req.Email).Return(nil, errors.New("not found"))
		mockUserRepo.EXPECT().Create(gomock.Any()).Return(errors.New("db error"))

		_, err := authService.Register(req)
		if err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("CartCreationFailure", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmail(req.Email).Return(nil, errors.New("not found"))
		mockUserRepo.EXPECT().Create(gomock.Any()).Return(nil)
		mockCartRepo.EXPECT().Create(gomock.Any()).Return(errors.New("cart error"))
		mockUserRepo.EXPECT().CreateRefreshToken(gomock.Any()).Return(nil)
		mockPublisher.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		_, err := authService.Register(req)
		if err != nil {
			t.Errorf("expected no error even if cart fails (per implementation), got %v", err)
		}
	})
}

func TestAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockCartRepo := mocks.NewMockCartRepositoryInterface(ctrl)
	mockPublisher := mocks.NewMockPublisher(ctrl)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:              "test_secret",
			ExpiresIn:           time.Hour,
			RefreshTokenExpires: time.Hour * 24,
		},
	}

	authService := services.NewAuthService(cfg, mockPublisher, mockUserRepo, mockCartRepo)

	hashedPassword, _ := utils.HashPassword("password123")
	user := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: hashedPassword,
		IsActive: true,
	}

	t.Run("Success", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmailAndActive(user.Email, true).Return(user, nil)
		mockUserRepo.EXPECT().CreateRefreshToken(gomock.Any()).Return(nil)
		mockPublisher.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		resp, err := authService.Login(&dto.LoginRequest{Email: user.Email, Password: "password123"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.User.ID != user.ID {
			t.Errorf("expected user ID %d, got %d", user.ID, resp.User.ID)
		}
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmailAndActive("wrong@example.com", true).Return(nil, errors.New("not found"))
		_, err := authService.Login(&dto.LoginRequest{Email: "wrong@example.com", Password: "any"})
		if err == nil || err.Error() != "invalid credentials" {
			t.Errorf("expected 'invalid credentials' error, got %v", err)
		}
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmailAndActive(user.Email, true).Return(user, nil)
		_, err := authService.Login(&dto.LoginRequest{Email: user.Email, Password: "wrong_password"})
		if err == nil || err.Error() != "invalid credentials" {
			t.Errorf("expected 'invalid credentials' error, got %v", err)
		}
	})
}

func TestAuthService_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockCartRepo := mocks.NewMockCartRepositoryInterface(ctrl)
	mockPublisher := mocks.NewMockPublisher(ctrl)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:              "test_secret",
			ExpiresIn:           time.Hour,
			RefreshTokenExpires: time.Hour * 24,
		},
	}

	authService := services.NewAuthService(cfg, mockPublisher, mockUserRepo, mockCartRepo)

	userID := uint(1)
	_, refreshToken, _ := utils.GenerateTokenPair(&cfg.JWT, userID, "test@example.com", "customer")

	t.Run("Success", func(t *testing.T) {
		mockUserRepo.EXPECT().GetValidRefreshToken(refreshToken).Return(&models.RefreshToken{ID: 10, Token: refreshToken}, nil)
		mockUserRepo.EXPECT().GetByID(userID).Return(&models.User{ID: userID}, nil)
		mockUserRepo.EXPECT().DeleteRefreshTokenByID(uint(10)).Return(nil)
		mockUserRepo.EXPECT().CreateRefreshToken(gomock.Any()).Return(nil)
		mockPublisher.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		resp, err := authService.RefreshToken(&dto.RefreshTokenRequest{RefreshToken: refreshToken})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.User.ID != userID {
			t.Errorf("expected user ID %d, got %d", userID, resp.User.ID)
		}
	})

	t.Run("InvalidTokenFormat", func(t *testing.T) {
		_, err := authService.RefreshToken(&dto.RefreshTokenRequest{RefreshToken: "invalid"})
		if err == nil || err.Error() != "invalid refresh token" {
			t.Errorf("expected 'invalid refresh token' error, got %v", err)
		}
	})

	t.Run("TokenNotFound", func(t *testing.T) {
		mockUserRepo.EXPECT().GetValidRefreshToken(refreshToken).Return(nil, errors.New("not found"))
		_, err := authService.RefreshToken(&dto.RefreshTokenRequest{RefreshToken: refreshToken})
		if err == nil || err.Error() != "refresh token not found or expired" {
			t.Errorf("expected 'refresh token not found or expired' error, got %v", err)
		}
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockUserRepo.EXPECT().GetValidRefreshToken(refreshToken).Return(&models.RefreshToken{ID: 10, Token: refreshToken}, nil)
		mockUserRepo.EXPECT().GetByID(userID).Return(nil, errors.New("user not found"))

		_, err := authService.RefreshToken(&dto.RefreshTokenRequest{RefreshToken: refreshToken})
		if err == nil || err.Error() != "user not found" {
			t.Errorf("expected 'user not found' error, got %v", err)
		}
	})
}

func TestAuthService_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	authService := services.NewAuthService(nil, nil, mockUserRepo, nil)

	token := "some_token"
	mockUserRepo.EXPECT().DeleteRefreshToken(token).Return(nil)

	err := authService.Logout(token)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAuthService_GenerateAuthResponse_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockPublisher := mocks.NewMockPublisher(ctrl)

	// Generate an error in Publish to trigger error path in generateAuthResponse
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:              "test_secret",
			ExpiresIn:           time.Hour,
			RefreshTokenExpires: time.Hour * 24,
		},
	}

	authService := services.NewAuthService(cfg, mockPublisher, mockUserRepo, nil)
	user := &models.User{ID: 1, Email: "test@example.com"}

	mockUserRepo.EXPECT().CreateRefreshToken(gomock.Any()).Return(nil)
	mockPublisher.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("publisher error"))

	// We can't call generateAuthResponse directly since it's unexported, but we can call Login or Register which calls it.
	mockUserRepo.EXPECT().GetByEmailAndActive(user.Email, true).Return(user, nil)
	// CheckPassword will also be called, but we don't have a password set in user for this test, which is fine if we mock it or set it.
	user.Password, _ = utils.HashPassword("password")

	_, err := authService.Login(&dto.LoginRequest{Email: user.Email, Password: "password"})
	if err == nil {
		t.Fatal("expected error from publisher but got nil")
	}
}

func TestAuthService_GenerateAuthResponse_CreateRefreshTokenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockPublisher := mocks.NewMockPublisher(ctrl)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:              "test_secret",
			ExpiresIn:           time.Hour,
			RefreshTokenExpires: time.Hour * 24,
		},
	}

	authService := services.NewAuthService(cfg, mockPublisher, mockUserRepo, nil)
	user := &models.User{ID: 1, Email: "test@example.com"}

	mockUserRepo.EXPECT().CreateRefreshToken(gomock.Any()).Return(errors.New("db error"))
	// generateAuthResponse ignores CreateRefreshToken error but logs it.
	// It should still proceed to publish event.
	mockPublisher.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	user.Password, _ = utils.HashPassword("password")
	mockUserRepo.EXPECT().GetByEmailAndActive(user.Email, true).Return(user, nil)

	_, err := authService.Login(&dto.LoginRequest{Email: user.Email, Password: "password"})
	if err != nil {
		t.Fatalf("expected no error even if CreateRefreshToken fails, got %v", err)
	}
}
