package models

import "time"

type Dream struct {
	ID             string
	UserID         string
	Description    string
	Explanation    string
	Date           time.Time
}
