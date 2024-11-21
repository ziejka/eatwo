package db

import (
	"eatwo/models"
	"time"
)

func (d *Dream) ToModel() *models.Dream {
  date, _ := time.Parse(time.RFC822, d.Date)
  return &models.Dream{
    ID:          d.ID,
    UserID:      d.UserID,
    Description: d.Description,
    Explanation: d.Explanation,
    Date:        date,
  }
}

func (u *User) ToModel() *models.User {
  return &models.User{
    ID:    u.ID,
    Name:  u.Name,
    Email: u.Email,
  }
}
