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

func (a AuthService) LogIn(ctx context.Context, logInData models.UserLogIn) (string, error) {
	user, err := a.userRepository.GetByEmail(ctx, logInData.Email)
	if err != nil {
		if errors.Is(err, shared.ErrNotExists) {
			return "", shared.ErrUserWrongEmailOrPassword
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(logInData.Password))
	if err != nil {
		return "", shared.ErrUserWrongEmailOrPassword
	}

	return generateToken(*user)
}

func (a AuthService) Create(ctx context.Context, signInData models.UserSignIn) error {
	_, err := a.userRepository.GetByEmail(ctx, signInData.Email)
	if err == nil {
		return shared.ErrUserWithEmailExist
	}
	if !errors.Is(err, shared.ErrNotExists) {
		return shared.ErrDefaultInternal
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signInData.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	userRecord := &models.UserRecord{
		User: models.User{
			Name:  signInData.Name,
			Email: signInData.Email,
		},
		HashPassword: string(hashedPassword),
	}

	return a.userRepository.Create(ctx, userRecord)
}
