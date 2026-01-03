package server

import (
	"github.com/gin-gonic/gin"
	"github.com/kuldeepstechwork/gocart-api/internal/config"
	"github.com/kuldeepstechwork/gocart-api/internal/services"
	"github.com/rs/zerolog"
)

// Server represents the HTTP server.
type Server struct {
	cfg            *config.Config
	log            *zerolog.Logger
	authService    *services.AuthService
	productService *services.ProductService
	userService    *services.UserService
	uploadService  *services.UploadService
	cartService    *services.CartService
	orderService   *services.OrderService
}

// New creates a new server instance.
func New(
	cfg *config.Config,
	log *zerolog.Logger,
	authService *services.AuthService,
	productService *services.ProductService,
	userService *services.UserService,
	uploadService *services.UploadService,
	cartService *services.CartService,
	orderService *services.OrderService,
) *Server {
	return &Server{
		cfg:            cfg,
		log:            log,
		authService:    authService,
		productService: productService,
		userService:    userService,
		uploadService:  uploadService,
		cartService:    cartService,
		orderService:   orderService,
	}
}

// SetupRoutes configures the server routes.
func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	return router
}
