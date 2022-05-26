package repository

import (
	"apriori/entity"
	"context"
	"database/sql"
	"errors"
	"strings"
)

type TransactionRepository interface {
	FindCandidate(ctx context.Context, tx *sql.Tx, data []string) ([]Test, error)
	FindItemSet(ctx context.Context, tx *sql.Tx) ([]entity.Transaction, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]entity.Transaction, error)
	FindByTransaction(ctx context.Context, tx *sql.Tx, noTransaction string) (entity.Transaction, error)
	Create(ctx context.Context, tx *sql.Tx, transaction entity.Transaction) error
	CreateFromCsv(ctx context.Context, tx *sql.Tx, transaction []entity.Transaction) error
	Update(ctx context.Context, tx *sql.Tx, transaction entity.Transaction) error
	Delete(ctx context.Context, tx *sql.Tx, noTransaction string) error
}

type transactionRepository struct {
}

func NewTransactionRepository() TransactionRepository {
	return &transactionRepository{}
}

type Test struct {
	productName string
}

func (repository *transactionRepository) FindCandidate(ctx context.Context, tx *sql.Tx, data []string) ([]Test, error) {
	query := `SELECT product_name FROM transactions WHERE CONTAINS(product_name, '"bantal cinta"')`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return []Test{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var test []Test
	for rows.Next() {
		var tes Test
		err := rows.Scan(
			&tes.productName,
		)

		if err != nil {
			return []Test{}, err
		}

		test = append(test, tes)
	}

	return test, nil
}
func (repository *transactionRepository) FindItemSet(ctx context.Context, tx *sql.Tx) ([]entity.Transaction, error) {
	query := `SELECT * FROM transactions`

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

func (repository *transactionRepository) FindAll(ctx context.Context, tx *sql.Tx) ([]entity.Transaction, error) {
	query := `SELECT * FROM transactions t ORDER BY t.id_transaction`

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
	query := `SELECT * FROM transactions WHERE no_transaction = ? LIMIT 1`

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

func (repository *transactionRepository) Create(ctx context.Context, tx *sql.Tx, transaction entity.Transaction) error {
	query := "INSERT INTO transactions(product_name,customer_name,no_transaction,created_at,updated_at) VALUES(?,?,?,?,?)"

	_, err := tx.ExecContext(
		ctx,
		query,
		transaction.ProductName,
		transaction.CustomerName,
		transaction.NoTransaction,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repository *transactionRepository) CreateFromCsv(ctx context.Context, tx *sql.Tx, transactions []entity.Transaction) error {
	query := `INSERT INTO transactions(product_name,customer_name,no_transaction,created_at,updated_at) VALUES `
	var values []interface{}

	for _, row := range transactions {
		productName := strings.ToLower(row.ProductName)
		query += "(?,?,?,?,?),"
		values = append(values, productName, row.CustomerName, row.NoTransaction, row.CreatedAt, row.UpdatedAt)
	}

	query = query[0 : len(query)-1]
	txNext, _ := tx.PrepareContext(ctx, query)
	_, err := txNext.ExecContext(
		ctx,
		values...,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repository *transactionRepository) Update(ctx context.Context, tx *sql.Tx, transaction entity.Transaction) error {
	query := `UPDATE transactions SET product_name = ?, customer_name = ?, updated_at = ? WHERE no_transaction = ?`

	_, err := tx.ExecContext(
		ctx,
		query,
		transaction.ProductName,
		transaction.CustomerName,
		transaction.UpdatedAt,
		transaction.NoTransaction,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repository *transactionRepository) Delete(ctx context.Context, tx *sql.Tx, noTransaction string) error {
	query := "DELETE FROM transactions WHERE no_transaction = ?"

	_, err := tx.ExecContext(ctx, query, noTransaction)
	if err != nil {
		return err
	}

	return nil
}
