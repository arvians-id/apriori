package entity

import (
	"time"
)

type User struct {
	IdUser    uint64
	Role      int
	Name      string
	Email     string
	Address   string
	Phone     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
