package repositories

import "gorm.io/gorm"

// UserRepository is a placeholder for a user repository.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}
