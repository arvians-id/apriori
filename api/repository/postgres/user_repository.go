package postgres

import (
	"apriori/entity"
	"apriori/repository"
	"context"
	"database/sql"
	"log"
)

type UserRepositoryImpl struct {
}

func NewUserRepository() repository.UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]*entity.User, error) {
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

	var users []*entity.User
	for rows.Next() {
		var user entity.User
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

func (repository *UserRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (*entity.User, error) {
	query := "SELECT * FROM users WHERE id_user = $1"
	row := tx.QueryRowContext(ctx, query, id)

	var user entity.User
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

func (repository *UserRepositoryImpl) FindByEmail(ctx context.Context, tx *sql.Tx, email string) (*entity.User, error) {
	query := "SELECT * FROM users WHERE email = $1"
	row := tx.QueryRowContext(ctx, query, email)

	var user entity.User
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

func (repository *UserRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, user *entity.User) (*entity.User, error) {
	id := 0
	query := "INSERT INTO users (role,name,email,address,phone,password,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id_user"
	row := tx.QueryRowContext(
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
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}

	user.IdUser = id

	return user, nil
}

func (repository *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user *entity.User) (*entity.User, error) {
	query := "UPDATE users SET name = $1, email = $2, address = $3, phone = $4, password = $5, updated_at = $6 WHERE id_user = $7"
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
		return nil, err
	}

	return user, nil
}

func (repository *UserRepositoryImpl) UpdatePassword(ctx context.Context, tx *sql.Tx, user *entity.User) error {
	query := "UPDATE users SET password = $1, updated_at = $2 WHERE email = $3"
	_, err := tx.ExecContext(ctx, query, user.Password, user.UpdatedAt, user.Email)
	if err != nil {
		return err
	}

	return nil
}

func (repository *UserRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	query := "DELETE FROM users WHERE id_user = $1"
	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
