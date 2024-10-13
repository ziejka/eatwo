package services_test

import (
	"context"
	"database/sql"
	"eatwo/db"
	"eatwo/models"
	"eatwo/services"
	"eatwo/shared"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const dbFileName = "test_sqlite.db"

func getAuthUserService(t *testing.T) (*services.AuthService, func()) {
	sqlDB, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	userRepository := db.NewUserRepository(sqlDB)
	err = userRepository.Migrate(context.Background())
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	authService := services.NewAuthService(userRepository)

	_, err = authService.Create(context.Background(), models.UserSignUp{
		UserLogIn: models.UserLogIn{
			Email:    "existing@example.com",
			Password: "password",
		},
		Name: "Test User",
	})
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	close := func() {
		defer sqlDB.Close()
		os.Remove(dbFileName)
	}

	return authService, close
}

func TestAuthService_Validate(t *testing.T) {
	authService, close := getAuthUserService(t)
	defer close()

	// Test case 1: Valid login credentials
	_, err := authService.Validate(context.Background(), models.UserLogIn{
		Email:    "existing@example.com",
		Password: "password",
	})
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Test case 2: Invalid email credentials
	_, err = authService.Validate(context.Background(), models.UserLogIn{
		Email:    "nonexisting@example.com",
		Password: "password",
	})
	if err != shared.ErrUserWrongEmailOrPassword {
		t.Errorf("Expected ErrUserWrongEmailOrPassword, but got: %v", err)
	}

	// Test case 2: Invalid login credentials
	_, err = authService.Validate(context.Background(), models.UserLogIn{
		Email:    "existing@example.com",
		Password: "wrong",
	})
	if err != shared.ErrUserWrongEmailOrPassword {
		t.Errorf("Expected ErrUserWrongEmailOrPassword, but got: %v", err)
	}
}

func TestAuthService_Create(t *testing.T) {
	authService, close := getAuthUserService(t)
	defer close()

	// Test case 1: User with the same email already exists
	_, err := authService.Create(context.Background(), models.UserSignUp{
		UserLogIn: models.UserLogIn{
			Email:    "existing@example.com",
			Password: "password",
		},
		Name: "New User",
	})
	if err != shared.ErrUserWithEmailExist {
		t.Errorf("Expected ErrUserWithEmailExist, but got: %v", err)
	}

	// Test case 2: User with the same email does not exist
	_, err = authService.Create(context.Background(), models.UserSignUp{
		UserLogIn: models.UserLogIn{
			Email:    "new@example.com",
			Password: "password",
		},
		Name: "New User",
	})
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}
