package services

import (
	"github.com/kuldeepstechwork/gocart-api/internal/config"
	"github.com/kuldeepstechwork/gocart-api/internal/events"
	"github.com/kuldeepstechwork/gocart-api/internal/interfaces"
	"github.com/kuldeepstechwork/gocart-api/internal/repositories"
	"gorm.io/gorm"
)

// AuthService handles authentication.
type AuthService struct{}

// NewAuthService creates a new auth service.
func NewAuthService(cfg *config.Config, publisher *events.EventPublisher, userRepo *repositories.UserRepository, cartRepo *repositories.CartRepository) *AuthService {
	return &AuthService{}
}

// ProductService handles products.
type ProductService struct{}

// NewProductService creates a new product service.
func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{}
}

// UserService handles users.
type UserService struct{}

// NewUserService creates a new user service.
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{}
}

// CartService handles carts.
type CartService struct{}

// NewCartService creates a new cart service.
func NewCartService(db *gorm.DB) *CartService {
	return &CartService{}
}

// OrderService handles orders.
type OrderService struct{}

// NewOrderService creates a new order service.
func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{}
}

// UploadService handles uploads.
type UploadService struct{}

// NewUploadService creates a new upload service.
func NewUploadService(provider interfaces.UploadProvider) *UploadService {
	return &UploadService{}
}
