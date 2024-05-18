package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"eatwo/models"
	"eatwo/shared"
	"encoding/hex"
	"errors"
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

	hashedPassword := createHashedPassword(user.Salt, logInData.Password)
	if hashedPassword != user.HashPassword {
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

	salt, err := generateSalt()
	if err != nil {
		return err
	}
	hashedPassword := createHashedPassword(salt, signInData.Password)

	userRecord := &models.UserRecord{
		Email:        signInData.Email,
		Name:         signInData.Name,
		HashPassword: hashedPassword,
		Salt:         salt,
	}

	return a.userRepository.Create(ctx, userRecord)
}

func createHashedPassword(salt, password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password + salt))
	return hex.EncodeToString(hasher.Sum(nil))
}

func generateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}
