package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"eatwo/model"
	"eatwo/shared"
	"encoding/hex"
	"errors"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.UserRecord) error
	GetByEmail(ctx context.Context, email string) (*model.UserRecord, error)
}

type AuthService struct {
	userRepository UserRepository
}

func NewAuthService(userRepository UserRepository) *AuthService {
	return &AuthService{
		userRepository: userRepository,
	}
}

func (a AuthService) Validate(ctx context.Context, signInData *model.UserSignIn) error {
	return nil
}

func (a AuthService) Create(ctx context.Context, signInData *model.UserSignIn) error {
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
	hasher := sha256.New()
	hasher.Write([]byte(signInData.Password + salt))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	userRecord := &model.UserRecord{
		Email:        signInData.Email,
		Name:         signInData.Name,
		HashPassword: hashedPassword,
		Salt:         salt,
	}

	return a.userRepository.Create(ctx, userRecord)
}

func generateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}
