package handlers_test

import (
	"context"
	"eatwo/handlers"
	"eatwo/models"
	"eatwo/shared"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

type mockUserAuthService struct{}

func (m *mockUserAuthService) Create(ctx context.Context, signInData models.UserSignIn) error {
	if signInData.Email == "existing@example.com" {
		return shared.ErrUserWithEmailExist
	}
	if signInData.Password == "InternalError" {
		return shared.ErrDefaultInternal
	}
	return nil
}

func (m *mockUserAuthService) LogIn(ctx context.Context, signInData models.UserLogIn) (string, error) {
	if signInData.Email == "existing@example.com" && signInData.Password == "password" {
		return "token", nil
	}
	if signInData.Password == "InternalError" {
		return "", shared.ErrDefaultInternal
	}
	return "", shared.ErrUserWrongEmailOrPassword
}

func TestAuthHandler_LogInPostHandler(t *testing.T) {
	e := echo.New()
	f := make(url.Values)
	f.Set("email", "existing@example.com")
	f.Set("password", "password")
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
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
	f := make(url.Values)
	f.Set("email", "nonexisting@example.com")
	f.Set("password", "password")
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := handlers.NewAuthHandler(&mockUserAuthService{})
	err := handler.LogInPostHandler(c)

	if err.Error() != "code=401, message=wrong email or password" {
		t.Errorf("Expected %+v, but got: %+v", echo.NewHTTPError(http.StatusUnauthorized, "wrong email or password"), err)
	}
}

func TestAuthHandler_LogInPostHandler_InternalServerError(t *testing.T) {
	e := echo.New()
	f := make(url.Values)
	f.Set("email", "existing@example.com")
	f.Set("password", "InternalError")
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := handlers.NewAuthHandler(&mockUserAuthService{})
	err := handler.LogInPostHandler(c)

	if err.Error() != "code=500, message=something went wrong please try again" {
		t.Errorf("Expected no error, but got: %v", err)
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

	if err.Error() != "code=409, message=User with such email already exist" {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestAuthHandler_SignInPostHandler_InternalServerError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/signin", strings.NewReader(`{"email": "new@example.com", "password": "InternalError", "name": "New User"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := handlers.NewAuthHandler(&mockUserAuthService{})
	err := handler.SignInPostHandler(c)

	if err.Error() != "code=500, message=Could not create user" {
		t.Errorf("Expected no error, but got: %v", err)
	}
}