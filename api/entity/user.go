package entity

import (
	"time"
)

type User struct {
	IdUser       int             `json:"id_user"`
	Role         int             `json:"role"`
	Name         string          `json:"name"`
	Email        string          `json:"email"`
	Address      string          `json:"address"`
	Phone        string          `json:"phone"`
	Password     string          `json:"password"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	Notification []*Notification `json:"notification"`
	Payment      []*Payment      `json:"payment"`
}
