package services

import "database/sql"

type Auth struct {
	db *sql.DB
}

func ValidateUser() {
	println("ValidateUser")
}
