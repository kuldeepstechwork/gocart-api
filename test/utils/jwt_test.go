package utils_test

import (
	"testing"
	"time"

	"github.com/kuldeepstechwork/gocart-api/internal/config"
	"github.com/kuldeepstechwork/gocart-api/internal/utils"
)

func TestJWT(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:              "test_secret_key",
		ExpiresIn:           time.Minute * 15,
		RefreshTokenExpires: time.Hour * 24,
	}
	userID := uint(1)
	email := "test@example.com"
	role := "admin"

	// Test GenerateTokenPair
	accessToken, refreshToken, err := utils.GenerateTokenPair(cfg, userID, email, role)
	if err != nil {
		t.Fatalf("unexpected error generating token pair: %v", err)
	}

	if accessToken == "" || refreshToken == "" {
		t.Fatal("expected non-empty tokens")
	}

	// Test ValidateToken - Success (Access Token)
	claims, err := utils.ValidateToken(accessToken, cfg.Secret)
	if err != nil {
		t.Fatalf("unexpected error validating access token: %v", err)
	}
	if claims.UserID != userID || claims.Email != email || claims.Role != role {
		t.Errorf("unexpected claims: %+v", claims)
	}

	// Test ValidateToken - Success (Refresh Token)
	claims, err = utils.ValidateToken(refreshToken, cfg.Secret)
	if err != nil {
		t.Fatalf("unexpected error validating refresh token: %v", err)
	}
	if claims.UserID != userID || claims.Email != email || claims.Role != role {
		t.Errorf("unexpected claims: %+v", claims)
	}

	// Test ValidateToken - Invalid Secret
	_, err = utils.ValidateToken(accessToken, "wrong_secret")
	if err == nil {
		t.Error("expected error with wrong secret but got nil")
	}

	// Test ValidateToken - Malformed Token
	_, err = utils.ValidateToken("not.a.token", cfg.Secret)
	if err == nil {
		t.Error("expected error with malformed token but got nil")
	}

	// Test ValidateToken - Expired Token
	expiredCfg := &config.JWTConfig{
		Secret:    "test_secret",
		ExpiresIn: -time.Minute, // Already expired
	}
	expiredToken, _, _ := utils.GenerateTokenPair(expiredCfg, userID, email, role)
	_, err = utils.ValidateToken(expiredToken, expiredCfg.Secret)
	if err == nil {
		t.Error("expected error with expired token but got nil")
	}
}
