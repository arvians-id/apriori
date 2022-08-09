package mysql

import (
	"apriori/entity"
	"apriori/repository"
	"context"
	"database/sql"
	"errors"
)

type categoryRepository struct {
}

func NewCategoryRepository() repository.CategoryRepository {
	return &categoryRepository{}
}

func (repository *categoryRepository) FindAll(ctx context.Context, tx *sql.Tx) ([]entity.Category, error) {
	query := "SELECT * FROM categories ORDER BY id_category DESC"
	queryContext, err := tx.QueryContext(ctx, query)
	if err != nil {
		return []entity.Category{}, err
	}
	defer func(queryContext *sql.Rows) {
		err := queryContext.Close()
		if err != nil {
			return
		}
	}(queryContext)

	var categories []entity.Category
	for queryContext.Next() {
		var category entity.Category
		err := queryContext.Scan(&category.IdCategory, &category.Name, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return []entity.Category{}, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (repository *categoryRepository) FindById(ctx context.Context, tx *sql.Tx, categoryId int) (entity.Category, error) {
	query := "SELECT * FROM categories WHERE id_category = ?"
	queryContext, err := tx.QueryContext(ctx, query, categoryId)
	if err != nil {
		return entity.Category{}, err
	}
	defer func(queryContext *sql.Rows) {
		err := queryContext.Close()
		if err != nil {
			return
		}
	}(queryContext)

	var category entity.Category
	if queryContext.Next() {
		err := queryContext.Scan(&category.IdCategory, &category.Name, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return entity.Category{}, err
		}

		return category, nil
	}

	return category, errors.New("category not found")
}

func (repository *categoryRepository) Create(ctx context.Context, tx *sql.Tx, category entity.Category) (entity.Category, error) {
	query := "INSERT INTO categories (name,created_at,updated_at) VALUES(?,?,?)"
	row, err := tx.ExecContext(ctx, query, category.Name, category.CreatedAt, category.UpdatedAt)
	if err != nil {
		return entity.Category{}, err
	}
	id, err := row.LastInsertId()
	if err != nil {
		return entity.Category{}, err
	}

	category.IdCategory = int(id)

	return category, nil
}

func (repository *categoryRepository) Update(ctx context.Context, tx *sql.Tx, category entity.Category) (entity.Category, error) {
	query := "UPDATE categories SET name = ?, updated_at = ? WHERE id_category = ?"
	_, err := tx.ExecContext(ctx, query, category.Name, category.UpdatedAt, category.IdCategory)
	if err != nil {
		return entity.Category{}, err
	}

	return category, nil
}

func (repository *categoryRepository) Delete(ctx context.Context, tx *sql.Tx, categoryId int) error {
	query := "DELETE FROM categories WHERE id_category = ?"
	_, err := tx.ExecContext(ctx, query, categoryId)
	if err != nil {
		return err
	}

	return nil
}
