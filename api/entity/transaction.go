package entity

import "time"

type Transaction struct {
	IdTransaction int
	ProductName   string
	CustomerName  string
	NoTransaction string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
