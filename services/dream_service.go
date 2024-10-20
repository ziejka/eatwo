package services

import (
	"context"
	"eatwo/db"
	"eatwo/models"
	"time"

	"github.com/google/uuid"
)

type DreamRepository interface {
	CreateDream(ctx context.Context, arg db.CreateDreamParams) (db.Dream, error)
	UpdateDreamExplanation(ctx context.Context, arg db.UpdateDreamExplanationParams) (db.Dream, error)
	GetDreams(ctx context.Context, userID string) ([]db.Dream, error)
}

type DreamService struct {
	dreamRepository DreamRepository
}

func NewDreamService(dreamRepository DreamRepository) *DreamService {
	return &DreamService{
		dreamRepository: dreamRepository,
	}
}

func (s *DreamService) Create(ctx context.Context, prompt string, userID string) (*models.Dream, error) {
	now := time.Now()
	dreamParams := db.CreateDreamParams{
		ID:          uuid.NewString(),
		UserID:      userID,
		Description: prompt,
		Explanation: "",
		Date:        now.Format(time.RFC822),
	}

	dream, err := s.dreamRepository.CreateDream(ctx, dreamParams)
	if err != nil {
		return nil, err
	}

	return &models.Dream{
		ID:          dream.ID,
		UserID:      dream.UserID,
		Description: dream.Description,
		Explanation: dream.Explanation,
		Date:        now,
	}, nil
}

func (s *DreamService) UpdateExplanation(ctx context.Context, dreamID, explanation string, userID string) error {
	_, err := s.dreamRepository.UpdateDreamExplanation(ctx, db.UpdateDreamExplanationParams{
		ID:          dreamID,
		Explanation: explanation,
		UserID:      userID,
	})

	return err
}

func (s *DreamService) GetByUserID(ctx context.Context, userID string) ([]models.Dream, error) {
	dreamsRecords, err := s.dreamRepository.GetDreams(ctx, userID)
	if err != nil {
		return nil, err
	}

	var dreams []models.Dream
	for _, dream := range dreamsRecords {
		date, err := time.Parse(time.RFC822, dream.Date)
		if err != nil {
			return nil, err
		}

		dreams = append(dreams, models.Dream{
			ID:          dream.ID,
			UserID:      dream.UserID,
			Description: dream.Description,
			Explanation: dream.Explanation,
			Date:        date,
		})
	}

	return dreams, nil
}
