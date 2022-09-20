package mysql

import (
	"context"
	"database/sql"
	"github.com/arvians-id/apriori/model"
	"github.com/arvians-id/apriori/repository"
	"log"
)

type CategoryRepositoryImpl struct {
}

func NewCategoryRepository() repository.CategoryRepository {
	return &CategoryRepositoryImpl{}
}

func (repository *CategoryRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]*model.Category, error) {
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

	var categories []*model.Category
	for rows.Next() {
		var category model.Category
		err := rows.Scan(&category.IdCategory, &category.Name, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}

	return categories, nil
}

func (repository *CategoryRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (*model.Category, error) {
	query := "SELECT * FROM categories WHERE id_category = ?"
	row := tx.QueryRowContext(ctx, query, id)

	var category model.Category
	err := row.Scan(&category.IdCategory, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (repository *CategoryRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, category *model.Category) (*model.Category, error) {
	query := "INSERT INTO categories (name,created_at,updated_at) VALUES(?,?,?)"
	row, err := tx.ExecContext(ctx, query, category.Name, category.CreatedAt, category.UpdatedAt)
	if err != nil {
		return nil, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return nil, err
	}

	category.IdCategory = int(id)

	return category, nil
}

func (repository *CategoryRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, category *model.Category) (*model.Category, error) {
	query := "UPDATE categories SET name = ?, updated_at = ? WHERE id_category = ?"
	_, err := tx.ExecContext(ctx, query, category.Name, category.UpdatedAt, category.IdCategory)
	if err != nil {
		return &model.Category{}, err
	}

	return category, nil
}

func (repository *CategoryRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	query := "DELETE FROM categories WHERE id_category = ?"
	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
