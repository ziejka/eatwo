package models

import "time"

type Dream struct {
	ID          string
	UserID      string
	Description string
	Explanation string
	Date        time.Time
}

func (d *Dream) GetDateOnly() string {
	return d.Date.Format(time.DateOnly)
}

type DreamsByDate [][]*Dream
