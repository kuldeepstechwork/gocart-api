package repositories_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kuldeepstechwork/gocart-api/internal/models"
	"github.com/kuldeepstechwork/gocart-api/internal/repositories"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupCartRepositoryTest() (*repositories.CartRepository, sqlmock.Sqlmock, error) {
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

	return repositories.NewCartRepository(gormDB), mock, nil
}

func TestCartRepository_GetByUserID(t *testing.T) {
	repo, mock, err := setupCartRepositoryTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	userID := uint(1)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id"}).AddRow(1, userID)
		mock.ExpectQuery(`SELECT \* FROM "carts" WHERE user_id = \$1 AND "carts"\."deleted_at" IS NULL .* LIMIT .*`).
			WithArgs(userID, 1).
			WillReturnRows(rows)

		cart, err := repo.GetByUserID(userID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cart.UserID != userID {
			t.Errorf("expected userID %d, got %d", userID, cart.UserID)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "carts" WHERE user_id = \$1 .*`).
			WithArgs(userID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := repo.GetByUserID(userID)
		if err != gorm.ErrRecordNotFound {
			t.Errorf("expected ErrRecordNotFound, got %v", err)
		}
	})
}

func TestCartRepository_Create(t *testing.T) {
	repo, mock, err := setupCartRepositoryTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	cart := &models.Cart{UserID: 1}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "carts" .+ VALUES .+(RETURNING "id")?`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := repo.Create(cart)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestCartRepository_Update(t *testing.T) {
	repo, mock, err := setupCartRepositoryTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	cart := &models.Cart{ID: 1, UserID: 1}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "carts" SET .+ WHERE .* "id" = \$[0-9]+`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Update(cart)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestCartRepository_Delete(t *testing.T) {
	repo, mock, err := setupCartRepositoryTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	cartID := uint(1)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "carts" SET "deleted_at"=\$1 WHERE "carts"\."id" = \$2 AND "carts"\."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), cartID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Delete(cartID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
