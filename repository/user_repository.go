package repository

import (
	"apriori/entity"
	"context"
	"database/sql"
)

type UserRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]entity.User, error)
	FindById(ctx context.Context, tx *sql.Tx, userId uint64) (entity.User, error)
	Create(ctx context.Context, tx *sql.Tx, user entity.User) (entity.User, error)
	Update(ctx context.Context, tx *sql.Tx, user entity.User) (entity.User, error)
	Delete(ctx context.Context, tx *sql.Tx, user entity.User) error
}

type userRepository struct {
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (repository *userRepository) FindAll(ctx context.Context, tx *sql.Tx) ([]entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (repository *userRepository) FindById(ctx context.Context, tx *sql.Tx, userId uint64) (entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (repository *userRepository) Create(ctx context.Context, tx *sql.Tx, user entity.User) (entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (repository *userRepository) Update(ctx context.Context, tx *sql.Tx, user entity.User) (entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (repository *userRepository) Delete(ctx context.Context, tx *sql.Tx, user entity.User) error {
	//TODO implement me
	panic("implement me")
}
