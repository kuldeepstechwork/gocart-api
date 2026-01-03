package services_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kuldeepstechwork/gocart-api/internal/dto"
	"github.com/kuldeepstechwork/gocart-api/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupUserServiceTest() (*services.UserService, sqlmock.Sqlmock, error) {
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

	return services.NewUserService(gormDB), mock, nil
}

func TestUserService_GetProfile(t *testing.T) {
	s, mock, err := setupUserServiceTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	userID := uint(1)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(userID, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(userID, "test@example.com"))

		resp, err := s.GetProfile(userID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.ID != userID {
			t.Errorf("expected user ID %d, got %d", userID, resp.ID)
		}
	})
}

func TestUserService_UpdateProfile(t *testing.T) {
	s, mock, err := setupUserServiceTest()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}

	userID := uint(1)
	req := &dto.UpdateProfileRequest{
		FirstName: "Jane",
		LastName:  "Doe",
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(userID, "test@example.com"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "first_name", "last_name"}).
				AddRow(userID, "test@example.com", "Jane", "Doe"))

		resp, err := s.UpdateProfile(userID, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.FirstName != req.FirstName {
			t.Errorf("expected First Name %s, got %s", req.FirstName, resp.FirstName)
		}
	})
}
