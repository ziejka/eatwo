package services

import (
	"context"
	"eatwo/models"
	"time"

	"github.com/google/uuid"
)

type DreamRepository interface {
	Create(ctx context.Context, dream *models.DreamRecord) error
  UpdateExplanation(ctx context.Context, dreamID, explanation string, userID string) error 
	GetByUserID(ctx context.Context, userID string) ([]*models.DreamRecord, error)
}

type DreamService struct {
	dreamRepository DreamRepository
}

func NewDreamService(dreamRepository DreamRepository) *DreamService {
	return &DreamService{
		dreamRepository: dreamRepository,
	}
}

func (s *DreamService) Create(ctx context.Context, prompt string, userID string) (*models.DreamRecord, error) {
	dream := &models.DreamRecord{
		ID:          uuid.NewString(),
		UserID:      userID,
		Description: prompt,
		Explanation: "",
		Date:        time.Now(),
	}
	err := s.dreamRepository.Create(ctx, dream)
	if err != nil {
		return nil, err
	}
	return dream, nil
}

func (s *DreamService) UpdateExplanation(ctx context.Context, dreamID, explanation string, userID string) error {
  return s.dreamRepository.UpdateExplanation(ctx, dreamID, explanation, userID)
}

func (s *DreamService) GetByUserID(ctx context.Context, userID string) ([]*models.DreamRecord, error) {
	return s.dreamRepository.GetByUserID(ctx, userID)
}
