// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

type Dream struct {
	ID          string
	UserID      string
	Description string
	Explanation string
	Date        string
}

type User struct {
	ID           string
	Email        string
	Name         string
	HashPassword string
}
