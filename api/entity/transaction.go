package entity

import (
	"time"
)

type Transaction struct {
	IdTransaction int       `json:"id_transaction"`
	ProductName   string    `json:"product_name"`
	CustomerName  string    `json:"customer_name"`
	NoTransaction string    `json:"no_transaction"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
