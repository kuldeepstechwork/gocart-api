package services_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kuldeepstechwork/gocart-api/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupOrderServiceTest() (*services.OrderService, sqlmock.Sqlmock, error) {
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

	return services.NewOrderService(gormDB), mock, nil
}

func TestOrderService_CreateOrder(t *testing.T) {
	s, mock, err := setupOrderServiceTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	userID := uint(1)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()

		// 1. Get cart
		mock.ExpectQuery(`SELECT .* FROM "carts"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(10, userID))

		// 2. Preload CartItems.Product
		mock.ExpectQuery(`SELECT .* FROM "cart_items"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "cart_id", "product_id", "quantity"}).AddRow(100, 10, 1000, 1))
		mock.ExpectQuery(`SELECT .* FROM "products"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "price", "stock", "name"}).AddRow(1000, 100.0, 10, "Prod 1"))

		// 3. Update Product Stock (tx.Save)
		mock.ExpectExec(`UPDATE "products" SET`).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// 4. Create Order
		mock.ExpectQuery(`INSERT INTO "orders"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(500))

		// 5. Create OrderItems (Save Order handles this because of association)
		mock.ExpectQuery(`INSERT INTO "order_items"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(600))

		// 6. Clear Cart (Unscoped)
		mock.ExpectExec(`DELETE FROM "cart_items"`).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// 7. getOrderResponse (Preload Category)
		mock.ExpectQuery(`SELECT .* FROM "orders"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(500, userID))
		mock.ExpectQuery(`SELECT .* FROM "order_items"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "product_id"}).AddRow(600, 500, 1000))
		mock.ExpectQuery(`SELECT .* FROM "products"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "category_id"}).AddRow(1000, 50))
		mock.ExpectQuery(`SELECT .* FROM "categories"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(50))

		mock.ExpectCommit()

		resp, err := s.CreateOrder(userID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.ID != 500 {
			t.Errorf("expected order ID 500, got %d", resp.ID)
		}
	})
}

func TestOrderService_GetOrders(t *testing.T) {
	s, mock, err := setupOrderServiceTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	userID := uint(1)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM "orders"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT .* FROM "orders"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(500, userID))

		mock.ExpectQuery(`SELECT .* FROM "order_items"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "product_id"}).AddRow(600, 500, 1000))

		mock.ExpectQuery(`SELECT .* FROM "products"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "category_id"}).AddRow(1000, 50))

		mock.ExpectQuery(`SELECT .* FROM "categories"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(50))

		resp, meta, err := s.GetOrders(userID, 1, 10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp) != 1 {
			t.Errorf("expected 1 order, got %d", len(resp))
		}
		if meta.Total != 1 {
			t.Errorf("expected total 1, got %d", meta.Total)
		}
	})
}

func TestOrderService_GetOrder(t *testing.T) {
	s, mock, err := setupOrderServiceTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	userID := uint(1)
	orderID := uint(500)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .* FROM "orders"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(orderID, userID))

		mock.ExpectQuery(`SELECT .* FROM "order_items"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "product_id"}).AddRow(600, orderID, 1000))

		mock.ExpectQuery(`SELECT .* FROM "products"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "category_id"}).AddRow(1000, 50))

		mock.ExpectQuery(`SELECT .* FROM "categories"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(50))

		resp, err := s.GetOrder(userID, orderID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.ID != orderID {
			t.Errorf("expected order ID %d, got %d", orderID, resp.ID)
		}
	})
}
