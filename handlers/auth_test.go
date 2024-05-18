package handlers_test

import (
	"context"
	"eatwo/handlers"
	"eatwo/models"
	"eatwo/shared"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

type mockUserAuthService struct{}

func (m *mockUserAuthService) Create(ctx context.Context, signInData models.UserSignIn) error {
	if signInData.Email == "existing@example.com" {
		return shared.ErrUserWithEmailExist
	}
	return nil
}

func (m *mockUserAuthService) LogIn(ctx context.Context, signInData models.UserLogIn) (string, error) {
	if signInData.Email == "existing@example.com" && signInData.Password == "password" {
		return "token", nil
	}
	return "", shared.ErrUserWrongEmailOrPassword
}

func TestAuthHandler_LogInPostHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"email": "existing@example.com", "password": "password"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := handlers.NewAuthHandler(&mockUserAuthService{})
	err := handler.LogInPostHandler(c)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got: %d", http.StatusOK, rec.Code)
	}

	expectedToken := "Bearer token"
	if rec.Header().Get("Authorization") != expectedToken {
		t.Errorf("Expected Authorization header to be %s, but got: %s", expectedToken, rec.Header().Get("Authorization"))
	}

	expectedBody := "Logged in successfully"
	if rec.Body.String() != expectedBody {
		t.Errorf("Expected response body to be %s, but got: %s", expectedBody, rec.Body.String())
	}
}

func TestAuthHandler_LogInPostHandler_WrongEmailOrPassword(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"email": "nonexisting@example.com", "password": "password"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := handlers.NewAuthHandler(&mockUserAuthService{})
	err := handler.LogInPostHandler(c)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, but got: %d", http.StatusUnauthorized, rec.Code)
	}

	expectedBody := "wrong email or password"
	if rec.Body.String() != expectedBody {
		t.Errorf("Expected response body to be %s, but got: %s", expectedBody, rec.Body.String())
	}
}

func TestAuthHandler_LogInPostHandler_InternalServerError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"email": "existing@example.com", "password": "wrong"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := handlers.NewAuthHandler(&mockUserAuthService{})
	err := handler.LogInPostHandler(c)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, but got: %d", http.StatusInternalServerError, rec.Code)
	}

	expectedBody := "something went wrong please try again"
	if rec.Body.String() != expectedBody {
		t.Errorf("Expected response body to be %s, but got: %s", expectedBody, rec.Body.String())
	}
}

func TestAuthHandler_SignInPostHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/signin", strings.NewReader(`{"email": "new@example.com", "password": "password", "name": "New User"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := handlers.NewAuthHandler(&mockUserAuthService{})
	err := handler.SignInPostHandler(c)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, but got: %d", http.StatusCreated, rec.Code)
	}

	expectedBody := "user created"
	if rec.Body.String() != expectedBody {
		t.Errorf("Expected response body to be %s, but got: %s", expectedBody, rec.Body.String())
	}
}

func TestAuthHandler_SignInPostHandler_UserWithEmailExist(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/signin", strings.NewReader(`{"email": "existing@example.com", "password": "password", "name": "New User"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := handlers.NewAuthHandler(&mockUserAuthService{})
	err := handler.SignInPostHandler(c)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if rec.Code != http.StatusConflict {
		t.Errorf("Expected status code %d, but got: %d", http.StatusConflict, rec.Code)
	}

	expectedBody := "User with such email already exist"
	if rec.Body.String() != expectedBody {
		t.Errorf("Expected response body to be %s, but got: %s", expectedBody, rec.Body.String())
	}
}

func TestAuthHandler_SignInPostHandler_InternalServerError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/signin", strings.NewReader(`{"email": "existing@example.com", "password": "password", "name": "New User"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := handlers.NewAuthHandler(&mockUserAuthService{})
	err := handler.SignInPostHandler(c)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, but got: %d", http.StatusInternalServerError, rec.Code)
	}

	expectedBody := "Could not create user"
	if rec.Body.String() != expectedBody {
		t.Errorf("Expected response body to be %s, but got: %s", expectedBody, rec.Body.String())
	}
}
