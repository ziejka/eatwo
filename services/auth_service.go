package services

import (
	"context"
	"database/sql"
	"eatwo/db"
	"eatwo/models"
	"eatwo/shared"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUser(ctx context.Context, email string) (db.User, error)
}

type AuthService struct {
	userRepository UserRepository
}

func NewAuthService(userRepository UserRepository) *AuthService {
	return &AuthService{
		userRepository: userRepository,
	}
}

func (a *AuthService) Validate(ctx context.Context, logInData models.UserLogIn) (*models.User, error) {
	user, err := a.userRepository.GetUser(ctx, logInData.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrUserWrongEmailOrPassword
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(logInData.Password))
	if err != nil {
		return nil, shared.ErrUserWrongEmailOrPassword
	}

	return user.ToModel(), nil
}

func (a *AuthService) Create(ctx context.Context, signUpData models.UserSignUp) (*models.User, error) {
	_, err := a.userRepository.GetUser(ctx, signUpData.Email)
	if err == nil {
		return nil, shared.ErrUserWithEmailExist
	}
	println(err.Error())

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, shared.ErrDefaultInternal
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signUpData.Password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}

	createUserParams := db.CreateUserParams{
		ID:           uuid.NewString(),
		Name:         signUpData.Name,
		Email:        signUpData.Email,
		HashPassword: string(hashedPassword),
	}
	user, err := a.userRepository.CreateUser(ctx, createUserParams)

	if err != nil {
		return nil, err
	}

	return user.ToModel(), nil
}
