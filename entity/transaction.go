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

type TransactionProduct struct {
	TransactionId           uint64
	TransactionCustomerName string
	TransactionNo           string
	TransactionQuantity     int32
	TransactionCreatedAt    time.Time
	ProductId               uint64
	ProductCode             string
	ProductName             string
	ProductDescription      string
}
