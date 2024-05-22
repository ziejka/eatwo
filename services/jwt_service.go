package services

import (
	"eatwo/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func generateToken(user models.UserRecord) (string, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	expirationTime := time.Now().Add(24 * 30 * time.Hour)

	claims := &models.Claims{
		Role: "user",
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Name,
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}
