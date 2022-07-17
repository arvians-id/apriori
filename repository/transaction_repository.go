package repository

import (
	"apriori/entity"
	"context"
	"database/sql"
	"errors"
	"strings"
)

type TransactionRepository interface {
	FindItemSet(ctx context.Context, tx *sql.Tx, startDate string, endDate string) ([]entity.Transaction, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]entity.Transaction, error)
	FindByTransaction(ctx context.Context, tx *sql.Tx, noTransaction string) (entity.Transaction, error)
	Create(ctx context.Context, tx *sql.Tx, transaction entity.Transaction) (entity.Transaction, error)
	CreateFromCsv(ctx context.Context, tx *sql.Tx, transaction []entity.Transaction) error
	Update(ctx context.Context, tx *sql.Tx, transaction entity.Transaction) (entity.Transaction, error)
	Delete(ctx context.Context, tx *sql.Tx, noTransaction string) error
	Truncate(ctx context.Context, tx *sql.Tx) error
}

type transactionRepository struct {
}

func NewTransactionRepository() TransactionRepository {
	return &transactionRepository{}
}

func (repository *transactionRepository) FindItemSet(ctx context.Context, tx *sql.Tx, startDate string, endDate string) ([]entity.Transaction, error) {
	query := `SELECT * FROM transactions WHERE DATE(created_at) >= $1 AND DATE(created_at) <= $2`

	rows, err := tx.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return []entity.Transaction{}, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var transactions []entity.Transaction

	for rows.Next() {
		var transaction entity.Transaction
		err := rows.Scan(
			&transaction.IdTransaction,
			&transaction.ProductName,
			&transaction.CustomerName,
			&transaction.NoTransaction,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)

		if err != nil {
			return []entity.Transaction{}, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (repository *transactionRepository) FindAll(ctx context.Context, tx *sql.Tx) ([]entity.Transaction, error) {
	query := `SELECT * FROM transactions t ORDER BY t.id_transaction DESC`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return []entity.Transaction{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var transactions []entity.Transaction
	for rows.Next() {
		var transaction entity.Transaction
		err := rows.Scan(
			&transaction.IdTransaction,
			&transaction.ProductName,
			&transaction.CustomerName,
			&transaction.NoTransaction,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return []entity.Transaction{}, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (repository *transactionRepository) FindByTransaction(ctx context.Context, tx *sql.Tx, noTransaction string) (entity.Transaction, error) {
	query := `SELECT * FROM transactions WHERE no_transaction = $1 LIMIT 1`

	rows, err := tx.QueryContext(ctx, query, noTransaction)
	if err != nil {
		return entity.Transaction{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	if rows.Next() {
		var transaction entity.Transaction
		err := rows.Scan(
			&transaction.IdTransaction,
			&transaction.ProductName,
			&transaction.CustomerName,
			&transaction.NoTransaction,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)

		if err != nil {
			return entity.Transaction{}, err
		}

		return transaction, nil
	}

	return entity.Transaction{}, errors.New("transaction not found")
}

func (repository *transactionRepository) Create(ctx context.Context, tx *sql.Tx, transaction entity.Transaction) (entity.Transaction, error) {
	id := 0
	query := "INSERT INTO transactions(product_name,customer_name,no_transaction,created_at,updated_at) VALUES($1,$2,$3,$4,$5) RETURNING id_transaction"
	row := tx.QueryRowContext(
		ctx,
		query,
		transaction.ProductName,
		transaction.CustomerName,
		transaction.NoTransaction,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	)
	err := row.Scan(&id)
	if err != nil {
		return entity.Transaction{}, err
	}

	transaction.IdTransaction = uint64(id)

	return transaction, nil
}

func (repository *transactionRepository) CreateFromCsv(ctx context.Context, tx *sql.Tx, transactions []entity.Transaction) error {
	for _, transaction := range transactions {
		query := `INSERT INTO transactions(product_name,customer_name,no_transaction,created_at,updated_at) VALUES ($1,$2,$3,$4,$5)`
		productName := strings.ToLower(transaction.ProductName)
		_, err := tx.ExecContext(
			ctx,
			query,
			productName,
			transaction.CustomerName,
			transaction.NoTransaction,
			transaction.CreatedAt,
			transaction.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repository *transactionRepository) Update(ctx context.Context, tx *sql.Tx, transaction entity.Transaction) (entity.Transaction, error) {
	query := `UPDATE transactions SET product_name = $1, customer_name = $2, updated_at = $3 WHERE no_transaction = $4`

	_, err := tx.ExecContext(
		ctx,
		query,
		transaction.ProductName,
		transaction.CustomerName,
		transaction.UpdatedAt,
		transaction.NoTransaction,
	)
	if err != nil {
		return entity.Transaction{}, err
	}

	return transaction, nil
}

func (repository *transactionRepository) Delete(ctx context.Context, tx *sql.Tx, noTransaction string) error {
	query := "DELETE FROM transactions WHERE no_transaction = $1"

	_, err := tx.ExecContext(ctx, query, noTransaction)
	if err != nil {
		return err
	}

	return nil
}

func (repository *transactionRepository) Truncate(ctx context.Context, tx *sql.Tx) error {
	query := `DELETE FROM transactions`

	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
