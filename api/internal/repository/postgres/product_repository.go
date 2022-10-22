package postgres

import (
	"context"
	"database/sql"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/arvians-id/apriori/internal/repository"
	"log"
)

type ProductRepositoryImpl struct {
}

func NewProductRepository() repository.ProductRepository {
	return &ProductRepositoryImpl{}
}

func (repository *ProductRepositoryImpl) FindAllByAdmin(ctx context.Context, tx *sql.Tx) ([]*model.Product, error) {
	query := "SELECT * FROM products ORDER BY id_product DESC"
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

	var products []*model.Product
	for rows.Next() {
		var product model.Product
		err := rows.Scan(
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
			return nil, err
		}

		products = append(products, &product)
	}

	return products, nil
}

func (repository *ProductRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, search string, category string) ([]*model.Product, error) {
	query := `SELECT * FROM products 
			  WHERE LOWER(name) LIKE $1 AND LOWER(category) LIKE $2 AND is_empty = false 
			  ORDER BY id_product DESC`
	rows, err := tx.QueryContext(ctx, query, "%"+search+"%", "%"+category+"%")
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

	var products []*model.Product
	for rows.Next() {
		var product model.Product
		err := rows.Scan(
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
			return nil, err
		}

		products = append(products, &product)
	}

	return products, nil
}

func (repository *ProductRepositoryImpl) FindAllBySimilarCategory(ctx context.Context, tx *sql.Tx, category string) ([]*model.Product, error) {
	query := `SELECT * FROM products 
			  WHERE category SIMILAR TO $1 AND is_empty = false 
			  ORDER BY random() DESC LIMIT 4`
	rows, err := tx.QueryContext(ctx, query, "%("+category+")%")
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

	var products []*model.Product
	for rows.Next() {
		var product model.Product
		err := rows.Scan(
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
			return nil, err
		}

		products = append(products, &product)
	}

	return products, nil
}

func (repository *ProductRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (*model.Product, error) {
	query := "SELECT * FROM products WHERE id_product = $1"
	row := tx.QueryRowContext(ctx, query, id)

	var product model.Product
	err := row.Scan(
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
		return nil, err
	}

	return &product, nil
}

func (repository *ProductRepositoryImpl) FindByName(ctx context.Context, tx *sql.Tx, name string) (*model.Product, error) {
	query := "SELECT * FROM products WHERE name = $1"
	row := tx.QueryRowContext(ctx, query, name)

	var product model.Product
	err := row.Scan(
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
		return nil, err
	}

	return &product, nil
}

func (repository *ProductRepositoryImpl) FindByCode(ctx context.Context, tx *sql.Tx, code string) (*model.Product, error) {
	query := "SELECT * FROM products WHERE code = $1"
	row := tx.QueryRowContext(ctx, query, code)

	var product model.Product
	err := row.Scan(
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
		return nil, err
	}

	return &product, nil
}

func (repository *ProductRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, product *model.Product) (*model.Product, error) {
	id := 0
	query := `INSERT INTO products (code,name,description,price,category,is_empty,mass,image,created_at,updated_at) 
			  VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id_product`
	row := tx.QueryRowContext(
		ctx,
		query,
		product.Code,
		product.Name,
		product.Description,
		product.Price,
		product.Category,
		product.IsEmpty,
		product.Mass,
		product.Image,
		product.CreatedAt,
		product.UpdatedAt,
	)
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}

	product.IdProduct = id

	return product, nil
}

func (repository *ProductRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, product *model.Product) (*model.Product, error) {
	query := `UPDATE products 
			  SET name = $1, 
			      description = $2, 
			      price = $3,
			      category = $4,
			      is_empty = $5, 
			      mass = $6, 
			      image = $7, 
			      updated_at = $8
			  WHERE code = $9`
	_, err := tx.ExecContext(
		ctx,
		query,
		product.Name,
		product.Description,
		product.Price,
		product.Category,
		product.IsEmpty,
		product.Mass,
		product.Image,
		product.UpdatedAt,
		product.Code,
	)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (repository *ProductRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, code string) error {
	query := "DELETE FROM products WHERE code = $1"
	_, err := tx.ExecContext(ctx, query, code)
	if err != nil {
		return err
	}

	return nil
}
