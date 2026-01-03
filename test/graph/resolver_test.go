package graph

import (
	"context"
	"errors"
	"testing"

	"github.com/kuldeepstechwork/gocart-api/graph/resolver"
	"github.com/kuldeepstechwork/gocart-api/internal/dto"
	"github.com/kuldeepstechwork/gocart-api/test/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestQueryResolver_Categories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductService := mocks.NewMockProductServiceInterface(ctrl)
	r := resolver.NewResolver(nil, nil, mockProductService, nil, nil)
	query := r.Query()

	t.Run("success", func(t *testing.T) {
		expectedCategories := []dto.CategoryResponse{
			{ID: 1, Name: "Electronics"},
			{ID: 2, Name: "Books"},
		}
		mockProductService.EXPECT().GetCategories().Return(expectedCategories, nil)

		res, err := query.Categories(context.Background())

		assert.NoError(t, err)
		assert.Len(t, res, 2)
		assert.Equal(t, "Electronics", res[0].Name)
	})

	t.Run("error", func(t *testing.T) {
		mockProductService.EXPECT().GetCategories().Return(nil, errors.New("db error"))

		res, err := query.Categories(context.Background())

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestMutationResolver_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthServiceInterface(ctrl)
	r := resolver.NewResolver(mockAuthService, nil, nil, nil, nil)
	mutation := r.Mutation()

	t.Run("success", func(t *testing.T) {
		input := dto.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		expectedResponse := &dto.AuthResponse{
			AccessToken: "token123",
		}
		mockAuthService.EXPECT().Login(&input).Return(expectedResponse, nil)

		res, err := mutation.Login(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, "token123", res.AccessToken)
	})
}

func TestQueryResolver_Me(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserServiceInterface(ctrl)
	r := resolver.NewResolver(nil, mockUserService, nil, nil, nil)
	query := r.Query()

	t.Run("success", func(t *testing.T) {
		userID := uint(1)
		ctx := createAuthContext(userID, "user")
		expectedUser := &dto.UserResponse{
			ID:    userID,
			Email: "test@example.com",
		}
		mockUserService.EXPECT().GetProfile(userID).Return(expectedUser, nil)

		res, err := query.Me(ctx)

		assert.NoError(t, err)
		assert.Equal(t, userID, res.ID)
	})

	t.Run("unauthorized", func(t *testing.T) {
		res, err := query.Me(context.Background())

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
