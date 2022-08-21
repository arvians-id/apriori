package postgres

import (
	"apriori/entity"
	"apriori/repository"
	"context"
	"database/sql"
)

type passwordResetRepository struct {
}

func NewPasswordResetRepository() repository.PasswordResetRepository {
	return &passwordResetRepository{}
}

func (repository *passwordResetRepository) FindByEmailAndToken(ctx context.Context, tx *sql.Tx, passwordReset *entity.PasswordReset) (*entity.PasswordReset, error) {
	query := "SELECT * FROM password_resets WHERE email = $1 AND token = $2"
	row := tx.QueryRowContext(ctx, query, passwordReset.Email, passwordReset.Token)

	err := row.Scan(&passwordReset.Email, &passwordReset.Token, &passwordReset.Expired)
	if err != nil {
		return &entity.PasswordReset{}, err
	}

	return passwordReset, nil

}

func (repository *passwordResetRepository) FindByEmail(ctx context.Context, tx *sql.Tx, email string) (*entity.PasswordReset, error) {
	query := "SELECT * FROM password_resets WHERE email = $1"
	row := tx.QueryRowContext(ctx, query, email)

	var passwordReset entity.PasswordReset
	err := row.Scan(&passwordReset.Email, &passwordReset.Token, &passwordReset.Expired)
	if err != nil {
		return &entity.PasswordReset{}, err
	}

	return &passwordReset, nil
}

func (repository *passwordResetRepository) Create(ctx context.Context, tx *sql.Tx, passwordReset *entity.PasswordReset) (*entity.PasswordReset, error) {
	query := "INSERT INTO password_resets (email,token,expired) VALUES($1,$2,$3)"
	_, err := tx.ExecContext(ctx, query, passwordReset.Email, passwordReset.Token, passwordReset.Expired)
	if err != nil {
		return &entity.PasswordReset{}, err
	}

	return passwordReset, nil
}

func (repository *passwordResetRepository) Update(ctx context.Context, tx *sql.Tx, passwordReset *entity.PasswordReset) (*entity.PasswordReset, error) {
	query := "UPDATE password_resets SET token = $1, expired = $2 WHERE email = $3"
	_, err := tx.ExecContext(ctx, query, passwordReset.Token, passwordReset.Expired, passwordReset.Email)
	if err != nil {
		return &entity.PasswordReset{}, err
	}

	return passwordReset, nil
}

func (repository *passwordResetRepository) Delete(ctx context.Context, tx *sql.Tx, email string) error {
	query := "DELETE FROM password_resets WHERE email = $1"
	_, err := tx.ExecContext(ctx, query, email)
	if err != nil {
		return err
	}

	return nil
}
