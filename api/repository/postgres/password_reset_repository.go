package postgres

import (
	"context"
	"database/sql"
	"github.com/arvians-id/apriori/model"
	"github.com/arvians-id/apriori/repository"
)

type PasswordResetRepositoryImpl struct {
}

func NewPasswordResetRepository() repository.PasswordResetRepository {
	return &PasswordResetRepositoryImpl{}
}

func (repository *PasswordResetRepositoryImpl) FindByEmailAndToken(ctx context.Context, tx *sql.Tx, passwordReset *model.PasswordReset) (*model.PasswordReset, error) {
	query := "SELECT * FROM password_resets WHERE email = $1 AND token = $2"
	row := tx.QueryRowContext(ctx, query, passwordReset.Email, passwordReset.Token)

	err := row.Scan(&passwordReset.Email, &passwordReset.Token, &passwordReset.Expired)
	if err != nil {
		return nil, err
	}

	return passwordReset, nil

}

func (repository *PasswordResetRepositoryImpl) FindByEmail(ctx context.Context, tx *sql.Tx, email string) (*model.PasswordReset, error) {
	query := "SELECT * FROM password_resets WHERE email = $1"
	row := tx.QueryRowContext(ctx, query, email)

	var passwordReset model.PasswordReset
	err := row.Scan(&passwordReset.Email, &passwordReset.Token, &passwordReset.Expired)
	if err != nil {
		return nil, err
	}

	return &passwordReset, nil
}

func (repository *PasswordResetRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, passwordReset *model.PasswordReset) (*model.PasswordReset, error) {
	query := "INSERT INTO password_resets (email,token,expired) VALUES($1,$2,$3)"
	_, err := tx.ExecContext(ctx, query, passwordReset.Email, passwordReset.Token, passwordReset.Expired)
	if err != nil {
		return nil, err
	}

	return passwordReset, nil
}

func (repository *PasswordResetRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, passwordReset *model.PasswordReset) (*model.PasswordReset, error) {
	query := "UPDATE password_resets SET token = $1, expired = $2 WHERE email = $3"
	_, err := tx.ExecContext(ctx, query, passwordReset.Token, passwordReset.Expired, passwordReset.Email)
	if err != nil {
		return nil, err
	}

	return passwordReset, nil
}

func (repository *PasswordResetRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, email string) error {
	query := "DELETE FROM password_resets WHERE email = $1"
	_, err := tx.ExecContext(ctx, query, email)
	if err != nil {
		return err
	}

	return nil
}

func (repository *PasswordResetRepositoryImpl) Truncate(ctx context.Context, tx *sql.Tx) error {
	query := `DELETE FROM password_resets`
	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
