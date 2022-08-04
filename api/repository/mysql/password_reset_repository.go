package mysql

import (
	"apriori/entity"
	"apriori/repository"
	"context"
	"database/sql"
	"errors"
)

type passwordResetRepository struct {
}

func NewPasswordResetRepository() repository.PasswordResetRepository {
	return &passwordResetRepository{}
}

func (repository *passwordResetRepository) FindByEmailAndToken(ctx context.Context, tx *sql.Tx, reset entity.PasswordReset) (entity.PasswordReset, error) {
	query := "SELECT * FROM password_resets WHERE email = ? AND token = ?"
	rows, err := tx.QueryContext(ctx, query, reset.Email, reset.Token)
	if err != nil {
		return entity.PasswordReset{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var userRequestReset entity.PasswordReset
	if rows.Next() {
		err := rows.Scan(&userRequestReset.Email, &userRequestReset.Token, &userRequestReset.Expired)
		if err != nil {
			return entity.PasswordReset{}, err
		}

		return userRequestReset, nil
	}

	return userRequestReset, errors.New("invalid credentials")
}

func (repository *passwordResetRepository) FindByEmail(ctx context.Context, tx *sql.Tx, email string) (entity.PasswordReset, error) {
	query := "SELECT * FROM password_resets WHERE email = ?"
	rows, err := tx.QueryContext(ctx, query, email)
	if err != nil {
		return entity.PasswordReset{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var userRequestReset entity.PasswordReset
	if rows.Next() {
		err := rows.Scan(&userRequestReset.Email, &userRequestReset.Token, &userRequestReset.Expired)
		if err != nil {
			return entity.PasswordReset{}, err
		}

		return userRequestReset, nil
	}

	return userRequestReset, nil
}

func (repository *passwordResetRepository) Create(ctx context.Context, tx *sql.Tx, reset entity.PasswordReset) (entity.PasswordReset, error) {
	query := "INSERT INTO password_resets (email,token,expired) VALUES(?,?,?)"
	_, err := tx.ExecContext(ctx, query, reset.Email, reset.Token, reset.Expired)
	if err != nil {
		return entity.PasswordReset{}, err
	}

	return reset, nil
}

func (repository *passwordResetRepository) Update(ctx context.Context, tx *sql.Tx, reset entity.PasswordReset) (entity.PasswordReset, error) {
	query := "UPDATE password_resets SET token = ?, expired = ? WHERE email = ?"
	_, err := tx.ExecContext(ctx, query, reset.Token, reset.Expired, reset.Email)
	if err != nil {
		return entity.PasswordReset{}, err
	}

	return reset, nil
}

func (repository *passwordResetRepository) Delete(ctx context.Context, tx *sql.Tx, email string) error {
	query := "DELETE FROM password_resets WHERE email = ?"
	_, err := tx.ExecContext(ctx, query, email)
	if err != nil {
		return err
	}

	return nil
}
