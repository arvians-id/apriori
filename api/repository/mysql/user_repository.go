package mysql

import (
	"context"
	"database/sql"
	"github.com/arvians-id/apriori/model"
	"github.com/arvians-id/apriori/repository"
	"log"
)

type UserRepositoryImpl struct {
}

func NewUserRepository() repository.UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]*model.User, error) {
	query := "SELECT * FROM users ORDER BY id_user DESC"
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

	var users []*model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.IdUser,
			&user.Role,
			&user.Name,
			&user.Email,
			&user.Address,
			&user.Phone,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (repository *UserRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (*model.User, error) {
	query := "SELECT * FROM users WHERE id_user = ?"
	row := tx.QueryRowContext(ctx, query, id)

	var user model.User
	err := row.Scan(
		&user.IdUser,
		&user.Role,
		&user.Name,
		&user.Email,
		&user.Address,
		&user.Phone,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repository *UserRepositoryImpl) FindByEmail(ctx context.Context, tx *sql.Tx, email string) (*model.User, error) {
	query := "SELECT * FROM users WHERE email = ?"
	row := tx.QueryRowContext(ctx, query, email)

	var user model.User
	err := row.Scan(
		&user.IdUser,
		&user.Role,
		&user.Name,
		&user.Email,
		&user.Address,
		&user.Phone,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repository *UserRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, user *model.User) (*model.User, error) {
	query := "INSERT INTO users (role,name,email,address,phone,password,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?)"
	row, err := tx.ExecContext(
		ctx,
		query,
		user.Role,
		user.Name,
		user.Email,
		user.Address,
		user.Phone,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return &model.User{}, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return &model.User{}, err
	}

	user.IdUser = int(id)

	return user, nil
}

func (repository *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user *model.User) (*model.User, error) {
	query := "UPDATE users SET name = ?, email = ?, address = ?, phone = ?, password = ?, updated_at = ? WHERE id_user = ?"
	_, err := tx.ExecContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Address,
		user.Phone,
		user.Password,
		user.UpdatedAt,
		user.IdUser,
	)
	if err != nil {
		return &model.User{}, err
	}

	return user, nil
}

func (repository *UserRepositoryImpl) UpdatePassword(ctx context.Context, tx *sql.Tx, user *model.User) error {
	query := "UPDATE users SET password = ?, updated_at = ? WHERE email = ?"
	_, err := tx.ExecContext(ctx, query, user.Password, user.UpdatedAt, user.Email)
	if err != nil {
		return err
	}

	return nil
}

func (repository *UserRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	query := "DELETE FROM users WHERE id_user = ?"
	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
