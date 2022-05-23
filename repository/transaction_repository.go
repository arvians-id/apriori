package repository

import (
	"apriori/entity"
	"context"
	"database/sql"
	"errors"
)

type TransactionRepository interface {
	FindItemSet(ctx context.Context, tx *sql.Tx) ([]entity.ProductTransaction, error)
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

func (repository *transactionRepository) FindItemSet(ctx context.Context, tx *sql.Tx) ([]entity.ProductTransaction, error) {
	query := `SELECT 
						p.code,
						p.name,
       					COUNT(t.product_id) as transaction,
						ROUND((COUNT(t.product_id) / (SELECT COUNT(*) 
							FROM transactions t 
							WHERE DATE(t.created_at) >= "2022-05-20" 
							AND DATE(t.created_at) <= "2022-05-21") * 100),2) as support
					FROM products p
					INNER JOIN transactions t ON t.product_id = p.id_product
					WHERE DATE(t.created_at) >= "2022-05-20" 
					AND DATE(t.created_at) <= "2022-05-21"
					GROUP BY t.product_id
					HAVING support >= 20`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return []entity.ProductTransaction{}, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var transactions []entity.ProductTransaction

	for rows.Next() {
		var transaction entity.ProductTransaction
		err := rows.Scan(&transaction.Code, &transaction.ProductName, &transaction.Transaction, &transaction.Support)

		if err != nil {
			return []entity.ProductTransaction{}, err
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
		query += "(?,?,?,?,?),"
		values = append(values, row.ProductName, row.CustomerName, row.NoTransaction, row.CreatedAt, row.UpdatedAt)
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
