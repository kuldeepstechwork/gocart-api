package services_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kuldeepstechwork/gocart-api/internal/dto"
	"github.com/kuldeepstechwork/gocart-api/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupProductServiceTest() (*services.ProductService, sqlmock.Sqlmock, error) {
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

	return services.NewProductService(gormDB), mock, nil
}

func TestProductService_CreateCategory(t *testing.T) {
	s, mock, err := setupProductServiceTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	req := &dto.CreateCategoryRequest{
		Name:        "Test Category",
		Description: "Test Description",
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "categories"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		resp, err := s.CreateCategory(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Name != req.Name {
			t.Errorf("expected name %s, got %s", req.Name, resp.Name)
		}
	})
}

func TestProductService_GetCategories(t *testing.T) {
	s, mock, err := setupProductServiceTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .* FROM "categories"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Cat 1"))

		resp, err := s.GetCategories()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp) != 1 {
			t.Errorf("expected 1 category, got %d", len(resp))
		}
	})
}

func TestProductService_UpdateCategory(t *testing.T) {
	s, mock, err := setupProductServiceTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	id := uint(1)
	req := &dto.UpdateCategoryRequest{Name: "Updated"}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .* FROM "categories"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(id, "Old"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "categories" SET`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		resp, err := s.UpdateCategory(id, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Name != "Updated" {
			t.Errorf("expected name Updated, got %s", resp.Name)
		}
	})
}

func TestProductService_DeleteCategory(t *testing.T) {
	s, mock, err := setupProductServiceTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	id := uint(1)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "categories" SET "deleted_at"`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := s.DeleteCategory(id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestProductService_GetProducts(t *testing.T) {
	s, mock, err := setupProductServiceTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM "products"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))

		mock.ExpectQuery(`SELECT .* FROM "products"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Prod 1"))

		mock.ExpectQuery(`SELECT .* FROM "product_images"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "product_id"}).AddRow(100, 1))

		mock.ExpectQuery(`SELECT .* FROM "categories"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Cat 1"))

		resp, meta, err := s.GetProducts(1, 10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp) != 1 {
			t.Errorf("expected 1 product, got %d", len(resp))
		}
		if meta.Total != 10 {
			t.Errorf("expected 10 total items, got %d", meta.Total)
		}
	})
}

func TestProductService_GetProduct(t *testing.T) {
	s, mock, err := setupProductServiceTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	id := uint(1)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .* FROM "products"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(id, "Prod 1"))

		mock.ExpectQuery(`SELECT .* FROM "product_images"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "product_id"}).AddRow(100, id))

		mock.ExpectQuery(`SELECT .* FROM "categories"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Cat 1"))

		resp, err := s.GetProduct(id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.ID != id {
			t.Errorf("expected product ID %d, got %d", id, resp.ID)
		}
	})
}

func TestProductService_CreateProduct(t *testing.T) {
	s, mock, err := setupProductServiceTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	req := &dto.CreateProductRequest{
		CategoryID: 1,
		Name:       "New Prod",
		Price:      10.0,
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "products"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		mock.ExpectQuery(`SELECT .* FROM "products"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "New Prod"))
		mock.ExpectQuery(`SELECT .* FROM "product_images"`).WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectQuery(`SELECT .* FROM "categories"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Cat 1"))

		resp, err := s.CreateProduct(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Name != "New Prod" {
			t.Errorf("expected name New Prod, got %s", resp.Name)
		}
	})
}

func TestProductService_UpdateProduct(t *testing.T) {
	s, mock, err := setupProductServiceTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	id := uint(1)
	req := &dto.UpdateProductRequest{Name: "Updated"}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .* FROM "products"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(id, "Old"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "products" SET`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		mock.ExpectQuery(`SELECT .* FROM "products"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(id, "Updated"))
		mock.ExpectQuery(`SELECT .* FROM "product_images"`).WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectQuery(`SELECT .* FROM "categories"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Cat 1"))

		resp, err := s.UpdateProduct(id, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Name != "Updated" {
			t.Errorf("expected name Updated, got %s", resp.Name)
		}
	})
}

func TestProductService_DeleteProduct(t *testing.T) {
	s, mock, err := setupProductServiceTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	id := uint(1)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "products" SET "deleted_at"`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := s.DeleteProduct(id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestProductService_AddProductImage(t *testing.T) {
	s, mock, err := setupProductServiceTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM "product_images"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "product_images"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := s.AddProductImage(1, "url", "alt")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestProductService_SearchProducts(t *testing.T) {
	s, mock, err := setupProductServiceTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	req := &dto.SearchProductsRequest{
		Query: "test",
		Page:  1,
		Limit: 10,
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM "products"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		mock.ExpectQuery(`SELECT products\.\*, ts_rank\(search_vector, plainto_tsquery\('english', \$1\)\) as rank FROM "products"`).
			WithArgs("test", "test", true, 10).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "rank"}).AddRow(1, "Test Prod", 0.5))

		mock.ExpectQuery(`SELECT .* FROM "product_images"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "product_id"}).AddRow(100, 1))

		mock.ExpectQuery(`SELECT .* FROM "categories"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Cat 1"))

		resp, _, err := s.SearchProducts(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp) != 1 {
			t.Errorf("expected 1 search result, got %d", len(resp))
		}
	})
}
