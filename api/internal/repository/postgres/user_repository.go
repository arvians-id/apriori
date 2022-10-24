package postgres

import (
	"context"
	"database/sql"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/arvians-id/apriori/internal/repository"
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
		log.Println("[UserRepository][FindAll] problem querying to db, err: ", err.Error())
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("[UserRepository][FindAll] problem closing query from db, err: ", err.Error())
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
			log.Println("[UserRepository][FindAll] problem with scanning db row, err: ", err.Error())
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (repository *UserRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (*model.User, error) {
	query := "SELECT * FROM users WHERE id_user = $1"
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
		log.Println("[UserRepository][FindById] problem with scanning db row, err: ", err.Error())
		return nil, err
	}

	return &user, nil
}

func (repository *UserRepositoryImpl) FindByEmail(ctx context.Context, tx *sql.Tx, email string) (*model.User, error) {
	query := "SELECT * FROM users WHERE email = $1"
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
		log.Println("[UserRepository][FindByEmail] problem with scanning db row, err: ", err.Error())
		return nil, err
	}

	return &user, nil
}

func (repository *UserRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, user *model.User) (*model.User, error) {
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
		log.Println("[UserRepository][Create] problem with scanning db row, err: ", err.Error())
		return nil, err
	}

	user.IdUser = id

	return user, nil
}

func (repository *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user *model.User) (*model.User, error) {
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
		log.Println("[UserRepository][Update] problem querying to db, err: ", err.Error())
		return nil, err
	}

	return user, nil
}

func (repository *UserRepositoryImpl) UpdatePassword(ctx context.Context, tx *sql.Tx, user *model.User) error {
	query := "UPDATE users SET password = $1, updated_at = $2 WHERE email = $3"
	_, err := tx.ExecContext(ctx, query, user.Password, user.UpdatedAt, user.Email)
	if err != nil {
		log.Println("[UserRepository][UpdatePassword] problem querying to db, err: ", err.Error())
		return err
	}

	return nil
}

func (repository *UserRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	query := "DELETE FROM users WHERE id_user = $1"
	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		log.Println("[UserRepository][Delete] problem querying to db, err: ", err.Error())
		return err
	}

	return nil
}
