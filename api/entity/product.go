package entity

import "time"

type Product struct {
	IdProduct   int
	Code        string
	Name        string
	Description string
	Price       int
	Category    string
	IsEmpty     int
	Mass        int
	Image       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
