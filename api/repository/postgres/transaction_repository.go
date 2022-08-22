package postgres

import (
	"apriori/entity"
	"apriori/repository"
	"context"
	"database/sql"
	"log"
	"strings"
)

type TransactionRepositoryImpl struct {
}

func NewTransactionRepository() repository.TransactionRepository {
	return &TransactionRepositoryImpl{}
}

func (repository *TransactionRepositoryImpl) FindAllItemSet(ctx context.Context, tx *sql.Tx, startDate string, endDate string) ([]*entity.Transaction, error) {
	query := `SELECT * FROM transactions 
			  WHERE DATE(created_at) >= $1 AND DATE(created_at) <= $2`
	rows, err := tx.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(rows)

	var transactions []*entity.Transaction
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
			return nil, err
		}

		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

func (repository *TransactionRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]*entity.Transaction, error) {
	query := `SELECT * FROM transactions ORDER BY id_transaction DESC`
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(rows)

	var transactions []*entity.Transaction
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
			return nil, err
		}

		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

func (repository *TransactionRepositoryImpl) FindByNoTransaction(ctx context.Context, tx *sql.Tx, noTransaction string) (*entity.Transaction, error) {
	query := `SELECT * FROM transactions WHERE no_transaction = $1 LIMIT 1`
	row := tx.QueryRowContext(ctx, query, noTransaction)

	var transaction entity.Transaction
	err := row.Scan(
		&transaction.IdTransaction,
		&transaction.ProductName,
		&transaction.CustomerName,
		&transaction.NoTransaction,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)
	if err != nil {
		return &entity.Transaction{}, err
	}

	return &transaction, nil
}

func (repository *TransactionRepositoryImpl) CreateByCsv(ctx context.Context, tx *sql.Tx, transactions []*entity.Transaction) error {
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

func (repository *TransactionRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, transaction *entity.Transaction) (*entity.Transaction, error) {
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
		return &entity.Transaction{}, err
	}

	transaction.IdTransaction = id

	return transaction, nil
}

func (repository *TransactionRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, transaction *entity.Transaction) (*entity.Transaction, error) {
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
		return &entity.Transaction{}, err
	}

	return transaction, nil
}

func (repository *TransactionRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, noTransaction string) error {
	query := "DELETE FROM transactions WHERE no_transaction = $1"
	_, err := tx.ExecContext(ctx, query, noTransaction)
	if err != nil {
		return err
	}

	return nil
}

func (repository *TransactionRepositoryImpl) Truncate(ctx context.Context, tx *sql.Tx) error {
	query := `DELETE FROM transactions`
	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
