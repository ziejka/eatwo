package handlers

import (
	"context"
	"eatwo/models"
	"eatwo/services"
	"eatwo/views/pages"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AIService interface {
	Call(ctx context.Context, prompt string) (*models.AIResponse, error)
}

type DreamHandler struct {
	aiService AIService
}

func NewDreamHandler(aiService AIService) *DreamHandler {
	return &DreamHandler{
		aiService: aiService,
	}
}

// /api/v1/dream
func (l *DreamHandler) PostDream(c echo.Context) error {
	claims := c.Get("claims")
	if claims == nil {
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), nil)
	}

	jwtClaims, ok := claims.(*services.CustomClaims)
	if !ok {
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), nil)
	}

	prompt := c.QueryParam("prompt")

	claude, err := l.aiService.Call(c.Request().Context(), prompt)
	println(claude)
	if err != nil {
		c.Logger().Error(err)
	}

	return renderHTMX(c, http.StatusOK, pages.HomePage(jwtClaims.Subject), jwtClaims)
}
