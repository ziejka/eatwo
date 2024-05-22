package models

import "github.com/golang-jwt/jwt"

type User struct {
	Email string `form:"email"`
	Name  string `form:"name"`
}

type UserLogIn struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type UserSignIn struct {
	UserLogIn
	Name string `form:"name"`
}

type UserRecord struct {
	User
	HashPassword string `form:"hash_password"`
}

type Claims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}
