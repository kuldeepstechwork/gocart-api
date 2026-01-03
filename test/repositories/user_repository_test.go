package repositories_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kuldeepstechwork/gocart-api/internal/models"
	"github.com/kuldeepstechwork/gocart-api/internal/repositories"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupUserRepositoryTest() (*repositories.UserRepository, sqlmock.Sqlmock, error) {
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

	return repositories.NewUserRepository(gormDB), mock, nil
}

func TestUserRepository_GetByEmail(t *testing.T) {
	repo, mock, err := setupUserRepositoryTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	email := "test@example.com"

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(1, email)
		mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 AND "users"\."deleted_at" IS NULL .* LIMIT .*`).
			WithArgs(email, 1).
			WillReturnRows(rows)

		user, err := repo.GetByEmail(email)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.Email != email {
			t.Errorf("expected email %s, got %s", email, user.Email)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 .*`).
			WithArgs(email, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := repo.GetByEmail(email)
		if err != gorm.ErrRecordNotFound {
			t.Errorf("expected ErrRecordNotFound, got %v", err)
		}
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	repo, mock, err := setupUserRepositoryTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	userID := uint(1)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(userID, "test@example.com")
		mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 AND "users"\."deleted_at" IS NULL .* LIMIT .*`).
			WithArgs(userID, 1).
			WillReturnRows(rows)

		user, err := repo.GetByID(userID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.ID != userID {
			t.Errorf("expected ID %d, got %d", userID, user.ID)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 .*`).
			WithArgs(userID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := repo.GetByID(userID)
		if err != gorm.ErrRecordNotFound {
			t.Errorf("expected ErrRecordNotFound, got %v", err)
		}
	})
}

func TestUserRepository_GetByEmailAndActive(t *testing.T) {
	repo, mock, err := setupUserRepositoryTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	email := "test@example.com"

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "is_active"}).AddRow(1, email, true)
		mock.ExpectQuery(`SELECT \* FROM "users" WHERE \(email = \$1 AND is_active = \$2\) AND "users"\."deleted_at" IS NULL .* LIMIT .*`).
			WithArgs(email, true, 1).
			WillReturnRows(rows)

		user, err := repo.GetByEmailAndActive(email, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.Email != email {
			t.Errorf("expected email %s, got %s", email, user.Email)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "users" WHERE \(email = \$1 AND is_active = \$2\) .*`).
			WithArgs(email, true, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := repo.GetByEmailAndActive(email, true)
		if err != gorm.ErrRecordNotFound {
			t.Errorf("expected ErrRecordNotFound, got %v", err)
		}
	})
}

func TestUserRepository_Create(t *testing.T) {
	repo, mock, err := setupUserRepositoryTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	user := &models.User{
		Email: "new@example.com",
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "users" .+ VALUES .+(RETURNING "id")?`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := repo.Create(user)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestUserRepository_Update(t *testing.T) {
	repo, mock, err := setupUserRepositoryTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	user := &models.User{ID: 1, Email: "updated@example.com"}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET .+ WHERE .* "id" = \$[0-9]+`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Update(user)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestUserRepository_Delete(t *testing.T) {
	repo, mock, err := setupUserRepositoryTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	userID := uint(1)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET "deleted_at"=\$1 WHERE "users"\."id" = \$2 AND "users"\."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), userID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Delete(userID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestUserRepository_CreateRefreshToken(t *testing.T) {
	repo, mock, err := setupUserRepositoryTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	token := &models.RefreshToken{UserID: 1, Token: "abc"}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "refresh_tokens" .+ VALUES .+(RETURNING "id")?`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := repo.CreateRefreshToken(token)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestUserRepository_GetValidRefreshToken(t *testing.T) {
	repo, mock, err := setupUserRepositoryTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	tokenStr := "some-token"

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "token"}).AddRow(1, tokenStr)
		mock.ExpectQuery(`SELECT \* FROM "refresh_tokens" WHERE \(token = \$1 AND expires_at > \$2\) AND "refresh_tokens"\."deleted_at" IS NULL .* LIMIT .*`).
			WithArgs(tokenStr, sqlmock.AnyArg(), 1).
			WillReturnRows(rows)

		token, err := repo.GetValidRefreshToken(tokenStr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token.Token != tokenStr {
			t.Errorf("expected token %s, got %s", tokenStr, token.Token)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "refresh_tokens" WHERE \(token = \$1 AND expires_at > \$2\) .*`).
			WithArgs(tokenStr, sqlmock.AnyArg(), 1).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := repo.GetValidRefreshToken(tokenStr)
		if err != gorm.ErrRecordNotFound {
			t.Errorf("expected ErrRecordNotFound, got %v", err)
		}
	})
}

func TestUserRepository_DeleteRefreshToken(t *testing.T) {
	repo, mock, err := setupUserRepositoryTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	tokenStr := "some-token"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "refresh_tokens" SET "deleted_at"=\$1 WHERE token = \$2 AND "refresh_tokens"\."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), tokenStr).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.DeleteRefreshToken(tokenStr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestUserRepository_DeleteRefreshTokenByID(t *testing.T) {
	repo, mock, err := setupUserRepositoryTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	tokenID := uint(100)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "refresh_tokens" SET "deleted_at"=\$1 WHERE "refresh_tokens"\."id" = \$2 AND "refresh_tokens"\."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), tokenID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.DeleteRefreshTokenByID(tokenID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
