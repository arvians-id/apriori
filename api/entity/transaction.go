package entity

import (
	"apriori/model"
	"time"
)

type Transaction struct {
	IdTransaction int
	ProductName   string
	CustomerName  string
	NoTransaction string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (transaction *Transaction) ToTransactionResponse() *model.GetTransactionResponse {
	return &model.GetTransactionResponse{
		IdTransaction: transaction.IdTransaction,
		ProductName:   transaction.ProductName,
		CustomerName:  transaction.CustomerName,
		NoTransaction: transaction.NoTransaction,
		CreatedAt:     transaction.CreatedAt.String(),
		UpdatedAt:     transaction.UpdatedAt.String(),
	}
}
