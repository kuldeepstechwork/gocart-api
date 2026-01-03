package services_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kuldeepstechwork/gocart-api/internal/dto"
	"github.com/kuldeepstechwork/gocart-api/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupCartServiceTest() (*services.CartService, sqlmock.Sqlmock, error) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	return services.NewCartService(gormDB), mock, nil
}

func TestCartService_GetCart(t *testing.T) {
	s, mock, err := setupCartServiceTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	userID := uint(1)
	cartID := uint(10)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .* FROM "carts"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(cartID, userID))

		mock.ExpectQuery(`SELECT .* FROM "cart_items"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "cart_id", "product_id", "quantity"}).
				AddRow(100, cartID, 1000, 2))

		mock.ExpectQuery(`SELECT .* FROM "products"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "category_id", "name"}).
				AddRow(1000, 50, "Test Product"))

		mock.ExpectQuery(`SELECT .* FROM "categories"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
				AddRow(50, "Test Category"))

		resp, err := s.GetCart(userID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.ID != cartID {
			t.Errorf("expected cart ID %d, got %d", cartID, resp.ID)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .* FROM "carts"`).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := s.GetCart(userID)
		if err == nil {
			t.Error("expected error but got nil")
		}
	})
}

func TestCartService_AddToCart(t *testing.T) {
	s, mock, err := setupCartServiceTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	userID := uint(1)
	productID := uint(1000)
	req := &dto.AddToCartRequest{
		ProductID: productID,
		Quantity:  2,
	}

	t.Run("Success_NewItem", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .* FROM "products"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "stock", "price"}).AddRow(productID, 10, 50.0))

		mock.ExpectQuery(`SELECT .* FROM "carts"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(10, userID))

		mock.ExpectQuery(`SELECT .* FROM "cart_items"`).
			WillReturnError(gorm.ErrRecordNotFound)

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "cart_items"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(100))
		mock.ExpectCommit()

		mock.ExpectQuery(`SELECT .* FROM "carts"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))
		mock.ExpectQuery(`SELECT .* FROM "cart_items"`).WillReturnRows(sqlmock.NewRows([]string{"id"}))

		_, err := s.AddToCart(userID, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("Success_ExistingItem", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .* FROM "products"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "stock", "price"}).AddRow(productID, 10, 50.0))

		mock.ExpectQuery(`SELECT .* FROM "carts"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(10, userID))

		mock.ExpectQuery(`SELECT .* FROM "cart_items"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "quantity"}).AddRow(100, 1))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "cart_items" SET .*quantity.*`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		mock.ExpectQuery(`SELECT .* FROM "carts"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))
		mock.ExpectQuery(`SELECT .* FROM "cart_items"`).WillReturnRows(sqlmock.NewRows([]string{"id"}))

		_, err := s.AddToCart(userID, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("InsufficientStock_ExistingItem", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .* FROM "products"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "stock", "price"}).AddRow(productID, 2, 50.0))

		mock.ExpectQuery(`SELECT .* FROM "carts"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(10, userID))

		mock.ExpectQuery(`SELECT .* FROM "cart_items"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "quantity"}).AddRow(100, 2))

		_, err := s.AddToCart(userID, req)
		if err == nil || err.Error() != "insufficient stock" {
			t.Errorf("expected 'insufficient stock' error, got %v", err)
		}
	})
}

func TestCartService_UpdateCartItem(t *testing.T) {
	s, mock, err := setupCartServiceTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	userID := uint(1)
	itemID := uint(100)
	req := &dto.UpdateCartItemRequest{Quantity: 5}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .* FROM "cart_items" JOIN carts`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "product_id"}).AddRow(itemID, 1000))

		mock.ExpectQuery(`SELECT .* FROM "products"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "stock"}).AddRow(1000, 10))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "cart_items" SET .*quantity.*`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		mock.ExpectQuery(`SELECT .* FROM "carts"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))
		mock.ExpectQuery(`SELECT .* FROM "cart_items"`).WillReturnRows(sqlmock.NewRows([]string{"id"}))

		_, err := s.UpdateCartItem(userID, itemID, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestCartService_RemoveFromCart(t *testing.T) {
	s, mock, err := setupCartServiceTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	userID := uint(1)
	itemID := uint(100)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "cart_items" SET "deleted_at"`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := s.RemoveFromCart(userID, itemID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
