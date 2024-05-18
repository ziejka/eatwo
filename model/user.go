package model

type User struct {
	Email string `form:"email"`
	Name  string `form:"name"`
}

type UserLogIn struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type UserSignIn struct {
	Email    string `form:"email"`
	Name     string `form:"name"`
	Password string `form:"password"`
}

type UserRecord struct {
	Email        string `form:"email"`
	Name         string `form:"name"`
	HashPassword string `form:"hash_password"`
	Salt         string `form:"salt"`
}
