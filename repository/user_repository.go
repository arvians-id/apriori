package repository

import (
	"apriori/entity"
	"context"
	"database/sql"
	"errors"
)

type UserRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]entity.User, error)
	FindById(ctx context.Context, tx *sql.Tx, userId uint64) (entity.User, error)
	FindByEmail(ctx context.Context, tx *sql.Tx, email string) (entity.User, error)
	Create(ctx context.Context, tx *sql.Tx, user entity.User) error
	Update(ctx context.Context, tx *sql.Tx, user entity.User) error
	UpdatePassword(ctx context.Context, tx *sql.Tx, user entity.User) error
	Delete(ctx context.Context, tx *sql.Tx, userId uint64) error
}

type userRepository struct {
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (repository *userRepository) FindAll(ctx context.Context, tx *sql.Tx) ([]entity.User, error) {
	query := "SELECT * FROM users"
	queryContext, err := tx.QueryContext(ctx, query)
	if err != nil {
		return []entity.User{}, err
	}
	defer func(queryContext *sql.Rows) {
		err := queryContext.Close()
		if err != nil {
			return
		}
	}(queryContext)

	var users []entity.User
	for queryContext.Next() {
		var user entity.User
		err := queryContext.Scan(&user.IdUser, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return []entity.User{}, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (repository *userRepository) FindById(ctx context.Context, tx *sql.Tx, userId uint64) (entity.User, error) {
	query := "SELECT * FROM users WHERE id_user = ?"
	queryContext, err := tx.QueryContext(ctx, query, userId)
	if err != nil {
		return entity.User{}, err
	}
	defer func(queryContext *sql.Rows) {
		err := queryContext.Close()
		if err != nil {
			return
		}
	}(queryContext)

	var user entity.User
	if queryContext.Next() {
		err := queryContext.Scan(&user.IdUser, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return entity.User{}, err
		}

		return user, nil
	}

	return user, errors.New("user not found")
}

func (repository *userRepository) FindByEmail(ctx context.Context, tx *sql.Tx, email string) (entity.User, error) {
	query := "SELECT * FROM users WHERE email = ?"
	queryContext, err := tx.QueryContext(ctx, query, email)
	if err != nil {
		return entity.User{}, err
	}
	defer func(queryContext *sql.Rows) {
		err := queryContext.Close()
		if err != nil {
			return
		}
	}(queryContext)

	var user entity.User
	if queryContext.Next() {
		err := queryContext.Scan(&user.IdUser, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return entity.User{}, err
		}

		return user, nil
	}

	return user, errors.New("email not found")
}

func (repository *userRepository) Create(ctx context.Context, tx *sql.Tx, user entity.User) error {
	query := "INSERT INTO users (name,email,password,created_at,updated_at) VALUES(?,?,?,?,?)"
	_, err := tx.ExecContext(ctx, query, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (repository *userRepository) Update(ctx context.Context, tx *sql.Tx, user entity.User) error {
	query := "UPDATE users SET name = ?, email = ?, password = ?, updated_at = ? WHERE id_user = ?"
	_, err := tx.ExecContext(ctx, query, user.Name, user.Email, user.Password, user.UpdatedAt, user.IdUser)
	if err != nil {
		return err
	}

	return nil
}

func (repository *userRepository) UpdatePassword(ctx context.Context, tx *sql.Tx, user entity.User) error {
	query := "UPDATE users SET password = ?, updated_at = ? WHERE email = ?"
	_, err := tx.ExecContext(ctx, query, user.Password, user.UpdatedAt, user.Email)
	if err != nil {
		return err
	}

	return nil
}

func (repository *userRepository) Delete(ctx context.Context, tx *sql.Tx, userId uint64) error {
	query := "DELETE FROM users WHERE id_user = ?"
	_, err := tx.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}

	return nil
}
