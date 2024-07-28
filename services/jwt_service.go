package services

import (
	"eatwo/models"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Email  string
	UserID string
}

func GenerateToken(user models.User) (string, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Email,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		Email:  user.Email,
		UserID: user.ID,
	})

	return token.SignedString(jwtKey)
}

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookieToken, err := c.Cookie("token")
		if err != nil {
			return next(c)
		}

		claims, err := parseToken(cookieToken.Value)
		if err != nil {
			return next(c)
		}

		c.Set("claims", claims)
		return next(c)
	}
}

func parseToken(tokenString string) (*CustomClaims, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("JWT Token parsing failed: %v", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("JWT Token: Fail to cast claims")
	}
	return claims, nil
}
