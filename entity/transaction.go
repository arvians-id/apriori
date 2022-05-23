package entity

import "time"

type Transaction struct {
	IdTransaction uint64
	ProductName   string
	CustomerName  string
	NoTransaction string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
