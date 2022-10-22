package mysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/arvians-id/apriori/internal/http/presenter/response"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/arvians-id/apriori/internal/repository"
	"log"
)

type AprioriRepositoryImpl struct {
}

func NewAprioriRepository() repository.AprioriRepository {
	return &AprioriRepositoryImpl{}
}

func (repository *AprioriRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]*model.Apriori, error) {
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
			log.Println(err)
			return
		}
	}(rows)

	var apriories []*model.Apriori
	for rows.Next() {
		var apriori model.Apriori
		err := rows.Scan(&apriori.Code, &apriori.RangeDate, &apriori.CreatedAt, &apriori.IsActive)
		if err != nil {
			return nil, err
		}

		apriories = append(apriories, &apriori)
	}

	return apriories, nil
}

func (repository *AprioriRepositoryImpl) FindAllByActive(ctx context.Context, tx *sql.Tx) ([]*model.Apriori, error) {
	query := `SELECT * FROM apriories WHERE is_active = 1 ORDER BY discount DESC`
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

	var apriories []*model.Apriori
	for rows.Next() {
		var apriori model.Apriori
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

		apriories = append(apriories, &apriori)
	}

	if apriories == nil {
		return nil, errors.New(response.ErrorNotFound)
	}

	return apriories, nil
}

func (repository *AprioriRepositoryImpl) FindAllByCode(ctx context.Context, tx *sql.Tx, code string) ([]*model.Apriori, error) {
	query := `SELECT * FROM apriories WHERE code = ?`
	rows, err := tx.QueryContext(ctx, query, code)
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

	var apriories []*model.Apriori
	for rows.Next() {
		var apriori model.Apriori
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

		apriories = append(apriories, &apriori)
	}

	if apriories == nil {
		return nil, errors.New(response.ErrorNotFound)
	}

	return apriories, nil
}

func (repository *AprioriRepositoryImpl) FindByCodeAndId(ctx context.Context, tx *sql.Tx, code string, id int) (*model.Apriori, error) {
	query := `SELECT * FROM apriories WHERE code = ? AND id_apriori = ?`
	row := tx.QueryRowContext(ctx, query, code, id)

	var apriori model.Apriori
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
		return nil, err
	}

	return &apriori, nil
}

func (repository *AprioriRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, apriories []*model.Apriori) error {
	for _, apriori := range apriories {
		query := `INSERT INTO apriories(code,item,discount,support,confidence,range_date,is_active,image,created_at) 
				  VALUES (?,?,?,?,?,?,?,?,?)`
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

func (repository *AprioriRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, apriori *model.Apriori) (*model.Apriori, error) {
	query := `UPDATE apriories SET description = ?, image = ? WHERE code = ? AND id_apriori = ?`
	_, err := tx.ExecContext(ctx, query, apriori.Description, apriori.Image, apriori.Code, apriori.IdApriori)
	if err != nil {
		return nil, err
	}

	return apriori, nil
}

func (repository *AprioriRepositoryImpl) UpdateAllStatus(ctx context.Context, tx *sql.Tx, status bool) error {
	query := `UPDATE apriories SET is_active = ?`
	_, err := tx.ExecContext(ctx, query, status)
	if err != nil {
		return err
	}

	return nil
}

func (repository *AprioriRepositoryImpl) UpdateStatusByCode(ctx context.Context, tx *sql.Tx, code string, status bool) error {
	query := `UPDATE apriories SET is_active = ? WHERE code = ?`
	_, err := tx.ExecContext(ctx, query, status, code)
	if err != nil {
		return err
	}

	return nil
}

func (repository *AprioriRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, code string) error {
	query := `DELETE FROM apriories WHERE code = ?`
	_, err := tx.ExecContext(ctx, query, code)
	if err != nil {
		return err
	}

	return nil
}
