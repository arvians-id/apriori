package postgres

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

func (repository *productRepository) FindAllOnAdmin(ctx context.Context, tx *sql.Tx) ([]entity.Product, error) {
	query := "SELECT * FROM products ORDER BY id_product DESC"
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
		err := queryContext.Scan(
			&product.IdProduct,
			&product.Code,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Category,
			&product.IsEmpty,
			&product.Mass,
			&product.Image,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return []entity.Product{}, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (repository *productRepository) FindAll(ctx context.Context, tx *sql.Tx, search string, category string) ([]entity.Product, error) {
	query := "SELECT * FROM products WHERE LOWER(name) LIKE $1 AND LOWER(category) LIKE $2 AND is_empty = 0 ORDER BY id_product DESC"
	queryContext, err := tx.QueryContext(ctx, query, "%"+search+"%", "%"+category+"%")
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
		err := queryContext.Scan(
			&product.IdProduct,
			&product.Code,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Category,
			&product.IsEmpty,
			&product.Mass,
			&product.Image,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return []entity.Product{}, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (repository *productRepository) FindAllSimilarCategory(ctx context.Context, tx *sql.Tx, category string) ([]entity.Product, error) {
	query := "SELECT * FROM products WHERE category SIMILAR TO $1 AND is_empty = 0 ORDER BY random() DESC LIMIT 4"
	queryContext, err := tx.QueryContext(ctx, query, "%("+category+")%")
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
		err := queryContext.Scan(
			&product.IdProduct,
			&product.Code,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Category,
			&product.IsEmpty,
			&product.Mass,
			&product.Image,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return []entity.Product{}, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (repository *productRepository) FindById(ctx context.Context, tx *sql.Tx, productId uint64) (entity.Product, error) {
	query := "SELECT * FROM products WHERE id_product = $1"
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
		err := queryContext.Scan(
			&product.IdProduct,
			&product.Code,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Category,
			&product.IsEmpty,
			&product.Mass,
			&product.Image,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return entity.Product{}, err
		}

		return product, nil
	}

	return product, errors.New("product not found")
}

func (repository *productRepository) FindByName(ctx context.Context, tx *sql.Tx, name string) (entity.Product, error) {
	query := "SELECT * FROM products WHERE name = $1"
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
		err := queryContext.Scan(
			&product.IdProduct,
			&product.Code,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Category,
			&product.IsEmpty,
			&product.Mass,
			&product.Image,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return entity.Product{}, err
		}

		return product, nil
	}

	return product, errors.New("product not found")
}

func (repository *productRepository) FindByCode(ctx context.Context, tx *sql.Tx, code string) (entity.Product, error) {
	query := "SELECT * FROM products WHERE code = $1"
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
		err := queryContext.Scan(
			&product.IdProduct,
			&product.Code,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Category,
			&product.IsEmpty,
			&product.Mass,
			&product.Image,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return entity.Product{}, err
		}

		return product, nil
	}

	return product, errors.New("product not found")
}

func (repository *productRepository) Create(ctx context.Context, tx *sql.Tx, product entity.Product) (entity.Product, error) {
	id := 0
	query := "INSERT INTO products (code,name,description,price,category,is_empty,mass,image,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id_product"
	row := tx.QueryRowContext(ctx, query, product.Code, product.Name, product.Description, product.Price, product.Category, product.IsEmpty, product.Mass, product.Image, product.CreatedAt, product.UpdatedAt)
	err := row.Scan(&id)
	if err != nil {
		return entity.Product{}, err
	}

	product.IdProduct = uint64(id)

	return product, nil
}

func (repository *productRepository) Update(ctx context.Context, tx *sql.Tx, product entity.Product) (entity.Product, error) {
	query := "UPDATE products SET name = $1, description = $2, price = $3, category = $4, is_empty = $5, mass = $6, image = $7, updated_at = $8 WHERE code = $9"
	_, err := tx.ExecContext(ctx, query, product.Name, product.Description, product.Price, product.Category, product.IsEmpty, product.Mass, product.Image, product.UpdatedAt, product.Code)
	if err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func (repository *productRepository) Delete(ctx context.Context, tx *sql.Tx, code string) error {
	query := "DELETE FROM products WHERE code = $1"
	_, err := tx.ExecContext(ctx, query, code)
	if err != nil {
		return err
	}

	return nil
}
