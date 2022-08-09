package entity

import (
	"time"
)

type Category struct {
	IdCategory int
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
