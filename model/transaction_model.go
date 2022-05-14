package model

import "time"

type CreateTransactionRequest struct {
	IdProduct     uint64
	CustomerName  string
	NoTransaction string
	Quantity      int32
	CreatedAt     time.Time
}

type UpdateTransactionRequest struct {
	IdTransaction uint64
	IdProduct     uint64
	CustomerName  string
	Quantity      int32
}

type GetTransactionResponse struct {
	IdTransaction uint64
	IdProduct     uint64
	CustomerName  string
	NoTransaction string
	Quantity      int32
	CreatedAt     time.Time
}
