package entity

import "time"

type Transaction struct {
	IdTransaction uint64
	IdProduct     uint64
	CustomerName  string
	NoTransaction string
	Quantity      int32
	CreatedAt     time.Time
}
