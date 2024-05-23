package services

import (
	"context"
	"eatwo/models"
	"eatwo/shared"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.UserRecord) error
	GetByEmail(ctx context.Context, email string) (*models.UserRecord, error)
}

type AuthService struct {
	userRepository UserRepository
}

func NewAuthService(userRepository UserRepository) *AuthService {
	return &AuthService{
		userRepository: userRepository,
	}
}

func (a AuthService) Validate(ctx context.Context, logInData models.UserLogIn) (models.User, error) {
	user, err := a.userRepository.GetByEmail(ctx, logInData.Email)
	if err != nil {
		if errors.Is(err, shared.ErrNotExists) {
			return models.User{}, shared.ErrUserWrongEmailOrPassword
		}
		return models.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(logInData.Password))
	if err != nil {
		return models.User{}, shared.ErrUserWrongEmailOrPassword
	}

	return user.User, nil
}

func (a AuthService) Create(ctx context.Context, signInData models.UserSignIn) (models.User, error) {
	_, err := a.userRepository.GetByEmail(ctx, signInData.Email)
	if err == nil {
		return models.User{}, shared.ErrUserWithEmailExist
	}
	if !errors.Is(err, shared.ErrNotExists) {
		return models.User{}, shared.ErrDefaultInternal
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signInData.Password), bcrypt.MinCost)
	if err != nil {
		return models.User{}, err
	}

	userRecord := &models.UserRecord{
		User: models.User{
			Name:  signInData.Name,
			Email: signInData.Email,
		},
		HashPassword: string(hashedPassword),
	}

	if err = a.userRepository.Create(ctx, userRecord); err != nil {
		return models.User{}, err
	}

	return userRecord.User, nil
}
