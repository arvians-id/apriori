package postgres

import (
	"apriori/entity"
	"apriori/repository"
	"context"
	"database/sql"
	"errors"
)

type aprioriRepository struct {
}

func NewAprioriRepository() repository.AprioriRepository {
	return &aprioriRepository{}
}

func (repository *aprioriRepository) FindAll(ctx context.Context, tx *sql.Tx) ([]entity.Apriori, error) {
	query := `SELECT code,range_date,created_at,is_active 
			  FROM apriories 
			  GROUP BY code,range_date,created_at,is_active 
			  ORDER BY created_at DESC`
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

func (repository *aprioriRepository) FindAllByActive(ctx context.Context, tx *sql.Tx) ([]entity.Apriori, error) {
	query := `SELECT * FROM apriories WHERE is_active = 1 ORDER BY discount DESC`
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
			&apriori.Description,
			&apriori.Image,
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

func (repository *aprioriRepository) FindAllByCode(ctx context.Context, tx *sql.Tx, code string) ([]entity.Apriori, error) {
	query := `SELECT * FROM apriories WHERE code = $1`
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
			&apriori.Description,
			&apriori.Image,
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

func (repository *aprioriRepository) FindByCodeAndId(ctx context.Context, tx *sql.Tx, code string, id int) (entity.Apriori, error) {
	query := `SELECT * FROM apriories WHERE code = $1 AND id_apriori = $2`
	row := tx.QueryRowContext(ctx, query, code, id)

	var apriori entity.Apriori
	err := row.Scan(
		&apriori.IdApriori,
		&apriori.Code,
		&apriori.Item,
		&apriori.Discount,
		&apriori.Support,
		&apriori.Confidence,
		&apriori.RangeDate,
		&apriori.IsActive,
		&apriori.Description,
		&apriori.Image,
		&apriori.CreatedAt,
	)
	if err != nil {
		return entity.Apriori{}, err
	}

	return apriori, nil
}

func (repository *aprioriRepository) Create(ctx context.Context, tx *sql.Tx, apriories []entity.Apriori) error {
	for _, apriori := range apriories {
		query := `INSERT INTO apriories(code,item,discount,support,confidence,range_date,is_active,image,created_at) 
				  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
		_, err := tx.ExecContext(
			ctx,
			query,
			apriori.Code,
			apriori.Item,
			apriori.Discount,
			apriori.Support,
			apriori.Confidence,
			apriori.RangeDate,
			apriori.IsActive,
			apriori.Image,
			apriori.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repository *aprioriRepository) Update(ctx context.Context, tx *sql.Tx, apriori entity.Apriori) (entity.Apriori, error) {
	query := `UPDATE apriories SET description = $1, image = $2 WHERE code = $3 AND id_apriori = $4`
	_, err := tx.ExecContext(ctx, query, apriori.Description, apriori.Image, apriori.Code, apriori.IdApriori)
	if err != nil {
		return entity.Apriori{}, err
	}

	return apriori, nil
}

func (repository *aprioriRepository) UpdateAllStatus(ctx context.Context, tx *sql.Tx, status int) error {
	query := `UPDATE apriories SET is_active = $1`
	_, err := tx.ExecContext(ctx, query, status)
	if err != nil {
		return err
	}

	return nil
}

func (repository *aprioriRepository) UpdateStatusByCode(ctx context.Context, tx *sql.Tx, code string, status int) error {
	query := `UPDATE apriories SET is_active = $1 WHERE code = $2`
	_, err := tx.ExecContext(ctx, query, status, code)
	if err != nil {
		return err
	}

	return nil
}

func (repository *aprioriRepository) Delete(ctx context.Context, tx *sql.Tx, code string) error {
	query := `DELETE FROM apriories WHERE code = $1`
	_, err := tx.ExecContext(ctx, query, code)
	if err != nil {
		return err
	}

	return nil
}
