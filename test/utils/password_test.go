package utils_test

import (
	"testing"

	"github.com/kuldeepstechwork/gocart-api/internal/utils"
)

func TestPasswordHashing(t *testing.T) {
	password := "my_secret_password"

	// Test HashPassword
	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	if hash == "" {
		t.Fatal("expected non-empty hash")
	}

	if hash == password {
		t.Fatal("hash should not be the same as password")
	}

	// Test CheckPassword - Success case
	if !utils.CheckPassword(password, hash) {
		t.Error("expected password to match hash")
	}

	// Test CheckPassword - Failure case
	if utils.CheckPassword("wrong_password", hash) {
		t.Error("expected password not to match hash")
	}
}

func TestHashPassword_Error(t *testing.T) {
	// bcrypt has a limit of 72 bytes, but it doesn't normally return an error for long strings, it just truncates.
	// However, we can test with a very long string just in case.
	longPassword := make([]byte, 100)
	for i := range longPassword {
		longPassword[i] = 'a'
	}

	_, err := utils.HashPassword(string(longPassword))
	if err == nil {
		t.Error("expected error with long password (>72 bytes) but got nil")
	}
}
