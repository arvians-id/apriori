package entity

import "time"

type Product struct {
	IdProduct   uint64
	Code        string
	Name        string
	Description string
	Price       int
	Image       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
