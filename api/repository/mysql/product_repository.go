package mysql

import (
	"context"
	"database/sql"
	"github.com/arvians-id/apriori/entity"
	"github.com/arvians-id/apriori/repository"
	"log"
)

type ProductRepositoryImpl struct {
}

func NewProductRepository() repository.ProductRepository {
	return &ProductRepositoryImpl{}
}

func (repository *ProductRepositoryImpl) FindAllByAdmin(ctx context.Context, tx *sql.Tx) ([]*entity.Product, error) {
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

	var products []*entity.Product
	for rows.Next() {
		var product entity.Product
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

func (repository *ProductRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, search string, category string) ([]*entity.Product, error) {
	query := `SELECT * FROM products 
			  WHERE LOWER(name) LIKE ? AND LOWER(category) LIKE ? AND is_empty = 0 
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

	var products []*entity.Product
	for rows.Next() {
		var product entity.Product
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

func (repository *ProductRepositoryImpl) FindAllBySimilarCategory(ctx context.Context, tx *sql.Tx, category string) ([]*entity.Product, error) {
	query := `SELECT * FROM products 
			  WHERE category SIMILAR TO ? AND is_empty = 0 
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

	var products []*entity.Product
	for rows.Next() {
		var product entity.Product
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

func (repository *ProductRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (*entity.Product, error) {
	query := "SELECT * FROM products WHERE id_product = ?"
	row := tx.QueryRowContext(ctx, query, id)

	var product entity.Product
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

func (repository *ProductRepositoryImpl) FindByName(ctx context.Context, tx *sql.Tx, name string) (*entity.Product, error) {
	query := "SELECT * FROM products WHERE name = ?"
	row := tx.QueryRowContext(ctx, query, name)

	var product entity.Product
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

func (repository *ProductRepositoryImpl) FindByCode(ctx context.Context, tx *sql.Tx, code string) (*entity.Product, error) {
	query := "SELECT * FROM products WHERE code = ?"
	row := tx.QueryRowContext(ctx, query, code)

	var product entity.Product
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

func (repository *ProductRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, product *entity.Product) (*entity.Product, error) {
	query := `INSERT INTO products (code,name,description,price,category,is_empty,mass,image,created_at,updated_at) 
			  VALUES(?,?,?,?,?,?,?,?,?,?)`
	row, err := tx.ExecContext(
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
	if err != nil {
		return nil, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return nil, err
	}

	product.IdProduct = int(id)

	return product, nil
}

func (repository *ProductRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, product *entity.Product) (*entity.Product, error) {
	query := `UPDATE products 
			  SET name = ?, 
			      description = ?, 
			      price = ?,
			      category = ?,
			      is_empty = ?, 
			      mass = ?, 
			      image = ?, 
			      updated_at = ?
			  WHERE code = ?`
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
	query := "DELETE FROM products WHERE code = ?"
	_, err := tx.ExecContext(ctx, query, code)
	if err != nil {
		return err
	}

	return nil
}
