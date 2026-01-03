package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kuldeepstechwork/gocart-api/internal/config"
	"github.com/kuldeepstechwork/gocart-api/internal/models"
	"github.com/kuldeepstechwork/gocart-api/internal/server"
	"github.com/kuldeepstechwork/gocart-api/internal/utils"
	mocks "github.com/kuldeepstechwork/gocart-api/test/mocks/services"
	"go.uber.org/mock/gomock"
)

const TestJWTSecret = "test-secret"

type TestServer struct {
	Server         *server.Server
	AuthService    *mocks.MockAuthServiceInterface
	UserService    *mocks.MockUserServiceInterface
	ProductService *mocks.MockProductServiceInterface
	CartService    *mocks.MockCartServiceInterface
	OrderService   *mocks.MockOrderServiceInterface
	UploadService  *mocks.MockUploadServiceInterface
	Config         *config.Config
}

func setupTestServer(ctrl *gomock.Controller) *TestServer {
	authService := mocks.NewMockAuthServiceInterface(ctrl)
	userService := mocks.NewMockUserServiceInterface(ctrl)
	productService := mocks.NewMockProductServiceInterface(ctrl)
	cartService := mocks.NewMockCartServiceInterface(ctrl)
	orderService := mocks.NewMockOrderServiceInterface(ctrl)
	uploadService := mocks.NewMockUploadServiceInterface(ctrl)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:              TestJWTSecret,
			ExpiresIn:           time.Hour,
			RefreshTokenExpires: time.Hour * 24,
		},
	}
	gin.SetMode(gin.TestMode)

	srv := server.New(
		cfg,
		nil, // Logger
		authService,
		productService,
		userService,
		uploadService,
		cartService,
		orderService,
	)

	return &TestServer{
		Server:         srv,
		AuthService:    authService,
		UserService:    userService,
		ProductService: productService,
		CartService:    cartService,
		OrderService:   orderService,
		UploadService:  uploadService,
		Config:         cfg,
	}
}

func createTestToken(userID uint) string {
	jwtCfg := &config.JWTConfig{
		Secret:              TestJWTSecret,
		ExpiresIn:           time.Hour,
		RefreshTokenExpires: time.Hour * 24,
	}
	accessToken, _, _ := utils.GenerateTokenPair(jwtCfg, userID, "test@example.com", string(models.UserRoleCustomer))
	return accessToken
}

func createAdminToken(userID uint) string {
	jwtCfg := &config.JWTConfig{
		Secret:              TestJWTSecret,
		ExpiresIn:           time.Hour,
		RefreshTokenExpires: time.Hour * 24,
	}
	accessToken, _, _ := utils.GenerateTokenPair(jwtCfg, userID, "admin@example.com", string(models.UserRoleAdmin))
	return accessToken
}
