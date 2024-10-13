package models

import "time"

type DreamRecord struct {
	ID             string
	UserID         string
	Description    string
	Explanation    string
	Date           time.Time
}
