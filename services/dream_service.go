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

	return dream.ToModel(), nil
}

func (s *DreamService) UpdateExplanation(ctx context.Context, dreamID, explanation string, userID string) (*models.Dream, error) {
	dream, err := s.dreamRepository.UpdateDreamExplanation(ctx, db.UpdateDreamExplanationParams{
		ID:          dreamID,
		Explanation: explanation,
		UserID:      userID,
	})

	if err != nil {
		return nil, err
	}

	return dream.ToModel(), nil
}

func (s *DreamService) GetByUserID(ctx context.Context, userID string) (models.DreamsByDate, error) {
	dreamsRecords, err := s.dreamRepository.GetDreams(ctx, userID)

	if err != nil {
		return nil, err
	}

	dreams := make(models.DreamsByDate, 0)
	for _, dream := range dreamsRecords {
		d := dream.ToModel()
		if len(dreams) == 0 {
			dreams = append(dreams, []*models.Dream{d})
			continue
		}

		lastDay := dreams[len(dreams)-1]
		if dremInDay := lastDay[0]; len(lastDay) > 0 && dremInDay.GetDateOnly() == d.GetDateOnly() {
			dreams[len(dreams)-1] = append(lastDay, d)
		} else {
			dreams = append(dreams, []*models.Dream{d})
		}
	}

	return dreams, nil
}
