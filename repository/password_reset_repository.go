package repository

import (
	"apriori/entity"
	"context"
	"database/sql"
)

type PasswordResetRepository interface {
	Create(ctx context.Context, tx *sql.Tx, reset entity.PasswordReset) (entity.PasswordReset, error)
	Delete(ctx context.Context, tx *sql.Tx, reset entity.PasswordReset) error
}

type passwordResetRepository struct {
}

func NewPasswordResetRepository() PasswordResetRepository {
	return &passwordResetRepository{}
}

func (repository *passwordResetRepository) Create(ctx context.Context, tx *sql.Tx, reset entity.PasswordReset) (entity.PasswordReset, error) {
	query := "INSERT INTO password_resets (email,token,expired) VALUES(?,?,?)"
	_, err := tx.ExecContext(ctx, query, reset.Email, reset.Token, reset.Expired)
	if err != nil {
		return entity.PasswordReset{}, err
	}

	return reset, nil
}

func (repository *passwordResetRepository) Delete(ctx context.Context, tx *sql.Tx, reset entity.PasswordReset) error {
	query := "DELETE FROM password_resets WHERE email = ?"
	_, err := tx.ExecContext(ctx, query, reset.Email)
	if err != nil {
		return err
	}

	return nil
}
