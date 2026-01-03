package graph

import (
	"context"

	"github.com/kuldeepstechwork/gocart-api/internal/utils"
)

func createAuthContext(userID uint, role string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, utils.UserIDKey, userID)
	ctx = context.WithValue(ctx, utils.UserRoleKey, role)
	return ctx
}
