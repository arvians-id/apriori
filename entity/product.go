package entity

import "time"

type Product struct {
	IdProduct   uint64
	Code        string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
