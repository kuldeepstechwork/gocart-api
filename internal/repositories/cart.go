package repositories

import "gorm.io/gorm"

// CartRepository is a placeholder for a cart repository.
type CartRepository struct {
	db *gorm.DB
}

// NewCartRepository creates a new cart repository.
func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}
