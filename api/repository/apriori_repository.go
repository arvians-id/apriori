package repository

import (
	"apriori/entity"
	"context"
	"database/sql"
	"errors"
)

type AprioriRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]entity.Apriori, error)
	FindByActive(ctx context.Context, tx *sql.Tx) ([]entity.Apriori, error)
	FindByCode(ctx context.Context, tx *sql.Tx, code string) ([]entity.Apriori, error)
	ChangeAllStatus(ctx context.Context, tx *sql.Tx, status bool) error
	ChangeStatusByCode(ctx context.Context, tx *sql.Tx, code string, status bool) error
	Create(ctx context.Context, tx *sql.Tx, apriories []entity.Apriori) error
	Delete(ctx context.Context, tx *sql.Tx, code string) error
}

type aprioriRepository struct {
}

func NewAprioriRepository() AprioriRepository {
	return &aprioriRepository{}
}

func (repository *aprioriRepository) FindAll(ctx context.Context, tx *sql.Tx) ([]entity.Apriori, error) {
	query := `SELECT code,range_date,created_at,is_active FROM apriories GROUP BY code,range_date,created_at,is_active`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var apriories []entity.Apriori
	for rows.Next() {
		var apriori entity.Apriori
		err := rows.Scan(&apriori.Code, &apriori.RangeDate, &apriori.CreatedAt, &apriori.IsActive)
		if err != nil {
			return nil, err
		}

		apriories = append(apriories, apriori)
	}

	return apriories, nil
}

func (repository *aprioriRepository) FindByActive(ctx context.Context, tx *sql.Tx) ([]entity.Apriori, error) {
	query := `SELECT * FROM apriories WHERE is_active = 1`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var apriories []entity.Apriori
	for rows.Next() {
		var apriori entity.Apriori
		err := rows.Scan(
			&apriori.IdApriori,
			&apriori.Code,
			&apriori.Item,
			&apriori.Discount,
			&apriori.Support,
			&apriori.Confidence,
			&apriori.RangeDate,
			&apriori.IsActive,
			&apriori.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		apriories = append(apriories, apriori)
	}
	if apriories == nil {
		return nil, errors.New("data not found")
	}

	return apriories, nil
}

func (repository *aprioriRepository) FindByCode(ctx context.Context, tx *sql.Tx, code string) ([]entity.Apriori, error) {
	query := `SELECT * FROM apriories WHERE code = ?`

	rows, err := tx.QueryContext(ctx, query, code)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var apriories []entity.Apriori
	for rows.Next() {
		var apriori entity.Apriori
		err := rows.Scan(
			&apriori.IdApriori,
			&apriori.Code,
			&apriori.Item,
			&apriori.Discount,
			&apriori.Support,
			&apriori.Confidence,
			&apriori.RangeDate,
			&apriori.IsActive,
			&apriori.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		apriories = append(apriories, apriori)
	}
	if apriories == nil {
		return nil, errors.New("data not found")
	}

	return apriories, nil
}

func (repository *aprioriRepository) ChangeAllStatus(ctx context.Context, tx *sql.Tx, status bool) error {
	query := `UPDATE apriories SET is_active = ?`

	_, err := tx.ExecContext(ctx, query, status)
	if err != nil {
		return err
	}

	return nil
}

func (repository *aprioriRepository) ChangeStatusByCode(ctx context.Context, tx *sql.Tx, code string, status bool) error {
	query := `UPDATE apriories SET is_active = ? WHERE code = ?`

	_, err := tx.ExecContext(ctx, query, status, code)
	if err != nil {
		return err
	}

	return nil
}

func (repository *aprioriRepository) Create(ctx context.Context, tx *sql.Tx, apriories []entity.Apriori) error {
	query := `INSERT INTO apriories(code,item,discount,support,confidence,range_date,is_active,created_at) VALUES`
	var values []interface{}

	for _, row := range apriories {
		query += "(?,?,?,?,?,?,?,?),"
		values = append(
			values,
			row.Code,
			row.Item,
			row.Discount,
			row.Support,
			row.Confidence,
			row.RangeDate,
			row.IsActive,
			row.CreatedAt,
		)
	}

	query = query[0 : len(query)-1]
	txNext, _ := tx.PrepareContext(ctx, query)
	_, err := txNext.ExecContext(
		ctx,
		values...,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repository *aprioriRepository) Delete(ctx context.Context, tx *sql.Tx, code string) error {
	query := `DELETE FROM apriories WHERE code = ?`

	_, err := tx.ExecContext(ctx, query, code)
	if err != nil {
		return err
	}

	return nil
}
