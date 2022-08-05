package mysql

import (
	"apriori/entity"
	"apriori/repository"
	"context"
	"database/sql"
	"errors"
)

type productRepository struct {
}

func NewProductRepository() repository.ProductRepository {
	return &productRepository{}
}

func (repository *productRepository) FindAll(ctx context.Context, tx *sql.Tx, search string) ([]entity.Product, error) {
	query := "SELECT * FROM products WHERE name LIKE ? ORDER BY id_product DESC"
	queryContext, err := tx.QueryContext(ctx, query, "%"+search+"%")
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
		err := queryContext.Scan(&product.IdProduct, &product.Code, &product.Name, &product.Description, &product.Price, &product.Image, &product.CreatedAt, &product.UpdatedAt)
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
		err := queryContext.Scan(&product.IdProduct, &product.Code, &product.Name, &product.Description, &product.Price, &product.Image, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return entity.Product{}, err
		}

		return product, nil
	}

	return product, errors.New("product not found")
}

func (repository *productRepository) FindByName(ctx context.Context, tx *sql.Tx, name string) (entity.Product, error) {
	query := "SELECT * FROM products WHERE name = ?"
	queryContext, err := tx.QueryContext(ctx, query, name)
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
		err := queryContext.Scan(&product.IdProduct, &product.Code, &product.Name, &product.Description, &product.Price, &product.Image, &product.CreatedAt, &product.UpdatedAt)
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
		err := queryContext.Scan(&product.IdProduct, &product.Code, &product.Name, &product.Description, &product.Price, &product.Image, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return entity.Product{}, err
		}

		return product, nil
	}

	return product, errors.New("product not found")
}

func (repository *productRepository) Create(ctx context.Context, tx *sql.Tx, product entity.Product) (entity.Product, error) {
	query := "INSERT INTO products (code,name,description,price,image,created_at,updated_at) VALUES(?,?,?,?,?,?,?)"
	row, err := tx.ExecContext(ctx, query, product.Code, product.Name, product.Description, product.Price, product.Image, product.CreatedAt, product.UpdatedAt)
	if err != nil {
		return entity.Product{}, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return entity.Product{}, err
	}

	product.IdProduct = uint64(id)

	return product, nil
}

func (repository *productRepository) Update(ctx context.Context, tx *sql.Tx, product entity.Product) (entity.Product, error) {
	query := "UPDATE products SET name = ?, description = ?, price = ?, image = ?, updated_at = ? WHERE code = ?"
	_, err := tx.ExecContext(ctx, query, product.Name, product.Description, product.Price, product.Image, product.UpdatedAt, product.Code)
	if err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func (repository *productRepository) Delete(ctx context.Context, tx *sql.Tx, code string) error {
	query := "DELETE FROM products WHERE code = ?"
	_, err := tx.ExecContext(ctx, query, code)
	if err != nil {
		return err
	}

	return nil
}
