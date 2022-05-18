package repository

import (
	"apriori/entity"
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	VerifyCredential(ctx context.Context, tx *sql.Tx, email string, password string) (entity.User, error)
}

type authRepository struct {
}

func NewAuthRepository() AuthRepository {
	return &authRepository{}
}

func (repository *authRepository) VerifyCredential(ctx context.Context, tx *sql.Tx, email string, password string) (entity.User, error) {
	query := "SELECT * FROM users WHERE email = ?"
	rows, err := tx.QueryContext(ctx, query, email)
	if err != nil {
		return entity.User{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var user entity.User
	if rows.Next() {
		err := rows.Scan(&user.IdUser, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return entity.User{}, err
		}
	} else {
		return entity.User{}, errors.New("email not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return entity.User{}, errors.New("wrong password")
	}

	return user, nil
}
