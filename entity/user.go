package entity

import (
	"time"
)

type User struct {
	IdUser    uint64
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
