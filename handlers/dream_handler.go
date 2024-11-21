package handlers

import (
	"context"
	"eatwo/models"
	"eatwo/services"
	"eatwo/views/pages"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AIService interface {
	Call(ctx context.Context, prompt string) (*models.AIResponse, error)
}

type DreamUpdater interface {
	Create(ctx context.Context, prompt string, userID string) (*models.Dream, error)
	GetByUserID(ctx context.Context, userID string) ([]models.Dream, error)
	UpdateExplanation(ctx context.Context, dreamID, explanation string, userID string) (*models.Dream, error)
}

type DreamHandler struct {
	aiService    AIService
	dreamUpdater DreamUpdater
}

func NewDreamHandler(aiService AIService, du DreamUpdater) *DreamHandler {
	return &DreamHandler{
		aiService:    aiService,
		dreamUpdater: du,
	}
}

type DreamRequestBody struct {
	Prompt string `form:"prompt"`
}

func (l *DreamHandler) PostDream(c echo.Context) error {
	claims := c.Get("claims")
	if claims == nil {
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), false)
	}

	customClaims, ok := claims.(*services.CustomClaims)
	if !ok {
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), false)
	}

	var dreamRequestBody DreamRequestBody
	if err := c.Bind(&dreamRequestBody); err != nil {
		return renderError(c, http.StatusBadRequest, "Could not parse the request")
	}

	if dreamRequestBody.Prompt == "" {
		return renderError(c, http.StatusBadRequest, "Prompt cannot be empty")
	}

	dream, err := l.dreamUpdater.Create(c.Request().Context(), dreamRequestBody.Prompt, customClaims.UserID)
	if err != nil {
		return renderError(c, http.StatusBadRequest, fmt.Sprint("Could not create dream: ", err))
	}

	resp, err := l.aiService.Call(c.Request().Context(), dream.Description)
	if err != nil {
		c.Logger().Error(err)
		return renderError(c, http.StatusBadRequest, fmt.Sprint("Could not decode a dream: ", err))
	}

	dream, err = l.dreamUpdater.UpdateExplanation(c.Request().Context(), dream.ID, resp.Content, customClaims.UserID)
	if err != nil {
		return renderError(c, http.StatusBadRequest, fmt.Sprint("Could not update dream: ", err))
	}

	return render(c, http.StatusOK, pages.DreamPromptResponse(dream))
}
