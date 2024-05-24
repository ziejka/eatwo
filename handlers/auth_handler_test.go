package handlers_test

import (
	"bytes"
	"context"
	"eatwo/handlers"
	"eatwo/models"
	"eatwo/shared"
	"eatwo/views/pages"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

type mockUserAuthService struct{}

func (m *mockUserAuthService) Create(ctx context.Context, signInData models.UserSignIn) (models.User, error) {
	if signInData.Email == "existing@example.com" {
		return models.User{}, shared.ErrUserWithEmailExist
	}
	if signInData.Password == "InternalError" {
		return models.User{}, shared.ErrDefaultInternal
	}
	return models.User{
		Email: signInData.Email,
		Name:  signInData.Name,
	}, nil
}

func (m *mockUserAuthService) Validate(ctx context.Context, logInData models.UserLogIn) (models.User, error) {
	if logInData.Email == "existing@example.com" && logInData.Password == "password" {
		return models.User{
			Email: logInData.Email,
		}, nil
	}
	if logInData.Password == "InternalError" {
		return models.User{}, shared.ErrDefaultInternal
	}
	return models.User{}, shared.ErrUserWrongEmailOrPassword
}

func generateTokenMock(user models.User) (string, error) {
	if user.Email == "invalid" {
		return "", fmt.Errorf("Error")
	}
	return "token", nil
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

	handler := handlers.NewAuthHandler(&mockUserAuthService{}, generateTokenMock)
	err := handler.LogInPostHandler(c)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got: %d", http.StatusOK, rec.Code)
	}

	expectedCookieName := "token"
	cookie := rec.Header().Get("Set-Cookie")
	if !strings.Contains(cookie, expectedCookieName) {
		t.Errorf("Expected Set-Cookie header to contain %s, but got: %s", expectedCookieName, cookie)
	}

	var buf bytes.Buffer
	pages.HomePage("existing@example.com").Render(req.Context(), &buf)
	expectedBody := buf.String()
	body := rec.Body.String()
	if !strings.Contains(body, expectedBody) {
		t.Errorf("Expected response body to contain %s, but got: %s", expectedBody, body)
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

	handler := handlers.NewAuthHandler(&mockUserAuthService{}, generateTokenMock)
	err := handler.LogInPostHandler(c)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, but got: %d", http.StatusOK, rec.Code)
	}

	var buf bytes.Buffer
	pages.AuthError("Wrong email or password").Render(req.Context(), &buf)
	expectedBody := buf.String()
	body := rec.Body.String()
	if !strings.Contains(body, expectedBody) {
		t.Errorf("Expected response body to contain %s, but got: %s", expectedBody, body)
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

	handler := handlers.NewAuthHandler(&mockUserAuthService{}, generateTokenMock)
	err := handler.LogInPostHandler(c)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, but got: %d", http.StatusOK, rec.Code)
	}

	var buf bytes.Buffer
	pages.AuthError("something went wrong please try again").Render(req.Context(), &buf)
	expectedBody := buf.String()
	body := rec.Body.String()
	if !strings.Contains(body, expectedBody) {
		t.Errorf("Expected response body to contain %s, but got: %s", expectedBody, body)
	}
}

func TestAuthHandler_SignInPostHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/signin", strings.NewReader(`{"email": "new@example.com", "password": "password", "name": "New User"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := handlers.NewAuthHandler(&mockUserAuthService{}, generateTokenMock)
	err := handler.SignInPostHandler(c)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, but got: %d", http.StatusCreated, rec.Code)
	}

	var buf bytes.Buffer
	pages.HomePage("new@example.com").Render(req.Context(), &buf)
	expectedBody := buf.String()
	body := rec.Body.String()
	if !strings.Contains(body, expectedBody) {
		t.Errorf("Expected response body to contain %s, but got: %s", expectedBody, body)
	}
}

func TestAuthHandler_SignInPostHandler_UserWithEmailExist(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/signin", strings.NewReader(`{"email": "existing@example.com", "password": "password", "name": "New User"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := handlers.NewAuthHandler(&mockUserAuthService{}, generateTokenMock)
	err := handler.SignInPostHandler(c)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, but got: %d", http.StatusUnauthorized, rec.Code)
	}

	var buf bytes.Buffer
	pages.AuthError("User with that email already exist").Render(req.Context(), &buf)
	expectedBody := buf.String()
	body := rec.Body.String()
	if !strings.Contains(body, expectedBody) {
		t.Errorf("Expected response body to contain %s, but got: %s", expectedBody, body)
	}
}

func TestAuthHandler_SignInPostHandler_InternalServerError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/signin", strings.NewReader(`{"email": "new@example.com", "password": "InternalError", "name": "New User"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := handlers.NewAuthHandler(&mockUserAuthService{}, generateTokenMock)
	err := handler.SignInPostHandler(c)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, but got: %d", http.StatusUnauthorized, rec.Code)
	}

	var buf bytes.Buffer
	pages.AuthError("Could not create user").Render(req.Context(), &buf)
	expectedBody := buf.String()
	body := rec.Body.String()
	if !strings.Contains(body, expectedBody) {
		t.Errorf("Expected response body to contain %s, but got: %s", expectedBody, body)
	}
}
