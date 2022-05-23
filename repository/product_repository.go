package repository

import (
	"apriori/entity"
	"context"
	"database/sql"
	"errors"
)

type ProductRepository interface {
	FindItemSet(ctx context.Context, tx *sql.Tx) ([]entity.ProductTransaction, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]entity.Product, error)
	FindById(ctx context.Context, tx *sql.Tx, productId uint64) (entity.Product, error)
	FindByCode(ctx context.Context, tx *sql.Tx, code string) (entity.Product, error)
	Create(ctx context.Context, tx *sql.Tx, product entity.Product) error
	Update(ctx context.Context, tx *sql.Tx, product entity.Product) error
	Delete(ctx context.Context, tx *sql.Tx, code string) error
}

type productRepository struct {
}

func NewProductRepository() ProductRepository {
	return &productRepository{}
}

func (repository *productRepository) FindItemSet(ctx context.Context, tx *sql.Tx) ([]entity.ProductTransaction, error) {
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

func (repository *productRepository) FindAll(ctx context.Context, tx *sql.Tx) ([]entity.Product, error) {
	query := "SELECT * FROM products"
	queryContext, err := tx.QueryContext(ctx, query)
	if err != nil {
		return []entity.Product{}, err
	}
	defer func(queryContext *sql.Rows) {
		err := queryContext.Close()
		if err != nil {
			return
		}
	}(queryContext)

	var products []entity.Product
	for queryContext.Next() {
		var product entity.Product
		err := queryContext.Scan(&product.IdProduct, &product.Code, &product.Name, &product.Description, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return []entity.Product{}, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (repository *productRepository) FindById(ctx context.Context, tx *sql.Tx, productId uint64) (entity.Product, error) {
	query := "SELECT * FROM products WHERE id_product = ?"
	queryContext, err := tx.QueryContext(ctx, query, productId)
	if err != nil {
		return entity.Product{}, err
	}
	defer func(queryContext *sql.Rows) {
		err := queryContext.Close()
		if err != nil {
			return
		}
	}(queryContext)

	var product entity.Product
	if queryContext.Next() {
		err := queryContext.Scan(&product.IdProduct, &product.Code, &product.Name, &product.Description, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return entity.Product{}, err
		}

		return product, nil
	}

	return product, errors.New("product not found")
}

func (repository *productRepository) FindByCode(ctx context.Context, tx *sql.Tx, code string) (entity.Product, error) {
	query := "SELECT * FROM products WHERE code = ?"
	queryContext, err := tx.QueryContext(ctx, query, code)
	if err != nil {
		return entity.Product{}, err
	}
	defer func(queryContext *sql.Rows) {
		err := queryContext.Close()
		if err != nil {
			return
		}
	}(queryContext)

	var product entity.Product
	if queryContext.Next() {
		err := queryContext.Scan(&product.IdProduct, &product.Code, &product.Name, &product.Description, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return entity.Product{}, err
		}

		return product, nil
	}

	return product, errors.New("product not found")
}

func (repository *productRepository) Create(ctx context.Context, tx *sql.Tx, product entity.Product) error {
	query := "INSERT INTO products (code,name,description,created_at,updated_at) VALUES(?,?,?,?,?)"
	_, err := tx.ExecContext(ctx, query, product.Code, product.Name, product.Description, product.CreatedAt, product.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (repository *productRepository) Update(ctx context.Context, tx *sql.Tx, product entity.Product) error {
	query := "UPDATE products SET code = ?, name = ?, description = ?, updated_at = ? WHERE code = ?"
	_, err := tx.ExecContext(ctx, query, product.Code, product.Name, product.Description, product.UpdatedAt, product.Code)
	if err != nil {
		return err
	}

	return nil
}

func (repository *productRepository) Delete(ctx context.Context, tx *sql.Tx, code string) error {
	query := "DELETE FROM products WHERE code = ?"
	_, err := tx.ExecContext(ctx, query, code)
	if err != nil {
		return err
	}

	return nil
}
