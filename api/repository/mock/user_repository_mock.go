package mock

import (
	"apriori/entity"
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	Mock mock.Mock
}

func (repository *UserRepositoryMock) FindAll(ctx context.Context, tx *sql.Tx) ([]*entity.User, error) {
	arguments := repository.Mock.Called(ctx)
	if arguments.Get(0) == nil {
		return nil, errors.New("data not found")
	}

	return arguments.Get(0).([]*entity.User), nil
}
func (repository *UserRepositoryMock) FindById(ctx context.Context, tx *sql.Tx, id int) (*entity.User, error) {
	arguments := repository.Mock.Called(ctx, id)
	if arguments.Get(0) == nil {
		return nil, errors.New("data not found")
	}

	return arguments.Get(0).(*entity.User), nil

}
func (repository *UserRepositoryMock) FindByEmail(ctx context.Context, tx *sql.Tx, email string) (*entity.User, error) {
	arguments := repository.Mock.Called(ctx, email)
	if arguments.Get(0) == nil {
		return nil, errors.New("data not found")
	}

	return arguments.Get(0).(*entity.User), nil

}
func (repository *UserRepositoryMock) Create(ctx context.Context, tx *sql.Tx, user *entity.User) (*entity.User, error) {
	arguments := repository.Mock.Called(ctx)
	if arguments.Get(0) == nil {
		return nil, errors.New("data not found")
	}

	return arguments.Get(0).(*entity.User), nil
}
func (repository *UserRepositoryMock) Update(ctx context.Context, tx *sql.Tx, user *entity.User) (*entity.User, error) {
	arguments := repository.Mock.Called(ctx)
	if arguments.Get(0) == nil {
		return nil, errors.New("data not found")
	}

	return arguments.Get(0).(*entity.User), nil

}
func (repository *UserRepositoryMock) UpdatePassword(ctx context.Context, tx *sql.Tx, user *entity.User) error {
	arguments := repository.Mock.Called(ctx)
	if arguments.Get(0) == nil {
		return errors.New("data not found")
	}

	return nil

}
func (repository *UserRepositoryMock) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	arguments := repository.Mock.Called(ctx, id)
	if arguments.Get(0) == nil {
		return nil
	}

	return errors.New("something went wrong")
}
