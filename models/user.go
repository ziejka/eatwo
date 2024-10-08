package models

type User struct {
	ID    string `form:"id"`
	Email string `form:"email"`
	Name  string `form:"name"`
}

type UserLogIn struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type UserSignUp struct {
	UserLogIn
	Name string `form:"name"`
}

type UserRecord struct {
	User
	HashPassword string `form:"hash_password"`
}

type Claims struct {
	Name string
}
