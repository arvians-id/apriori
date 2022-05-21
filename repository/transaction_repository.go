package repository

import (
	"apriori/entity"
	"context"
	"database/sql"
	"errors"
)

type TransactionRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]entity.TransactionProduct, error)
	FindByTransaction(ctx context.Context, tx *sql.Tx, noTransaction string) (entity.TransactionProduct, error)
	Create(ctx context.Context, tx *sql.Tx, transaction entity.Transaction) (entity.Transaction, error)
	Update(ctx context.Context, tx *sql.Tx, transaction entity.Transaction) (entity.Transaction, error)
	Delete(ctx context.Context, tx *sql.Tx, noTransaction string) error
}

type transactionRepository struct {
}

func NewTransactionRepository() TransactionRepository {
	return &transactionRepository{}
}

func (repository transactionRepository) FindAll(ctx context.Context, tx *sql.Tx) ([]entity.TransactionProduct, error) {
	query := `SELECT 
						t.id_transaction,
						t.customer_name,
						t.no_transaction,
						t.quantity,
						t.created_at,
						p.id_product,
						p.code as product_code,
						p.name as product_name,
						p.description as product_description
					FROM transactions t
					LEFT JOIN products p ON p.id_product = t.product_id
					ORDER BY t.id_transaction`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return []entity.TransactionProduct{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var transactions []entity.TransactionProduct
	for rows.Next() {
		var transaction entity.TransactionProduct
		err := rows.Scan(
			&transaction.TransactionId,
			&transaction.TransactionCustomerName,
			&transaction.TransactionNo,
			&transaction.TransactionQuantity,
			&transaction.TransactionCreatedAt,
			&transaction.ProductId,
			&transaction.ProductCode,
			&transaction.ProductName,
			&transaction.ProductDescription,
		)
		if err != nil {
			return []entity.TransactionProduct{}, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (repository transactionRepository) FindByTransaction(ctx context.Context, tx *sql.Tx, noTransaction string) (entity.TransactionProduct, error) {
	query := `SELECT 
						t.id_transaction,
						t.customer_name,
						t.no_transaction,
						t.quantity,
						t.created_at,
						p.id_product,
						p.code as product_code,
						p.name as product_name,
						p.description as product_description
					FROM transactions t
					LEFT JOIN products p ON p.id_product = t.product_id
					WHERE t.no_transaction = ?
					LIMIT 1`

	rows, err := tx.QueryContext(ctx, query, noTransaction)
	if err != nil {
		return entity.TransactionProduct{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	if rows.Next() {
		var transaction entity.TransactionProduct
		err := rows.Scan(
			&transaction.TransactionId,
			&transaction.TransactionCustomerName,
			&transaction.TransactionNo,
			&transaction.TransactionQuantity,
			&transaction.TransactionCreatedAt,
			&transaction.ProductId,
			&transaction.ProductCode,
			&transaction.ProductName,
			&transaction.ProductDescription,
		)

		if err != nil {
			return entity.TransactionProduct{}, err
		}

		return transaction, nil
	}

	return entity.TransactionProduct{}, errors.New("transaction not found")
}

func (repository transactionRepository) Create(ctx context.Context, tx *sql.Tx, transaction entity.Transaction) (entity.Transaction, error) {
	query := `INSERT INTO transactions(product_id,customer_name,no_transaction,quantity,created_at) VALUES(?,?,?,?,?)`

	rows, err := tx.ExecContext(
		ctx,
		query,
		transaction.IdProduct,
		transaction.CustomerName,
		transaction.NoTransaction,
		transaction.Quantity,
		transaction.CreatedAt,
	)
	if err != nil {
		return entity.Transaction{}, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return entity.Transaction{}, err
	}

	transaction.IdTransaction = uint64(id)

	return transaction, nil
}

func (repository transactionRepository) Update(ctx context.Context, tx *sql.Tx, transaction entity.Transaction) (entity.Transaction, error) {
	query := `UPDATE transactions SET product_id = ?, customer_name = ?, quantity = ? WHERE no_transaction = ?`

	_, err := tx.ExecContext(
		ctx,
		query,
		transaction.IdProduct,
		transaction.CustomerName,
		transaction.Quantity,
		transaction.NoTransaction,
	)
	if err != nil {
		return entity.Transaction{}, err
	}

	return transaction, nil
}

func (repository transactionRepository) Delete(ctx context.Context, tx *sql.Tx, noTransaction string) error {
	query := "DELETE FROM transactions WHERE no_transaction = ?"

	_, err := tx.ExecContext(ctx, query, noTransaction)
	if err != nil {
		return err
	}

	return nil
}
