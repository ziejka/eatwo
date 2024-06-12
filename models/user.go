package models

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
	ID           string `form:"id"`
	HashPassword string `form:"hash_password"`
}

type Claims struct {
	Name string
}
