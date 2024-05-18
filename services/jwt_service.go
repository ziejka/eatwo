package services

import (
	"eatwo/models"
	"time"

	"github.com/golang-jwt/jwt"
)

// The string "my_secret_key" is just an example and should be replaced with a secret key of sufficient length and complexity in a real-world scenario.
var jwtKey = []byte("my_secret_key")

func generateToken(user models.UserRecord) (string, error) {
	expirationTime := time.Now().Add(24 * 30 * time.Hour)

	claims := &models.Claims{
		Role: "user",
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Email,
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}
