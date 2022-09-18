package postgres

import (
	"context"
	"database/sql"
	"github.com/arvians-id/apriori/entity"
	"github.com/arvians-id/apriori/repository"
	"log"
)

type CategoryRepositoryImpl struct {
}

func NewCategoryRepository() repository.CategoryRepository {
	return &CategoryRepositoryImpl{}
}

func (repository *CategoryRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]*entity.Category, error) {
	query := "SELECT * FROM categories ORDER BY id_category DESC"
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

	var categories []*entity.Category
	for rows.Next() {
		var category entity.Category
		err := rows.Scan(&category.IdCategory, &category.Name, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}

	return categories, nil
}

func (repository *CategoryRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (*entity.Category, error) {
	query := "SELECT * FROM categories WHERE id_category = $1"
	row := tx.QueryRowContext(ctx, query, id)

	var category entity.Category
	err := row.Scan(&category.IdCategory, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (repository *CategoryRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, category *entity.Category) (*entity.Category, error) {
	id := 0
	query := "INSERT INTO categories (name,created_at,updated_at) VALUES($1,$2,$3) RETURNING id_category"
	row := tx.QueryRowContext(ctx, query, category.Name, category.CreatedAt, category.UpdatedAt)
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}

	category.IdCategory = id

	return category, nil
}

func (repository *CategoryRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, category *entity.Category) (*entity.Category, error) {
	query := "UPDATE categories SET name = $1, updated_at = $2 WHERE id_category = $3"
	_, err := tx.ExecContext(ctx, query, category.Name, category.UpdatedAt, category.IdCategory)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (repository *CategoryRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	query := "DELETE FROM categories WHERE id_category = $1"
	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
