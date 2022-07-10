package repository

import (
	"apriori/entity"
	"context"
	"database/sql"
)

type PaymentRepository interface {
	Create(ctx context.Context, tx *sql.Tx, product entity.Product) (entity.Product, error)
}

type paymentRepository struct {
}
