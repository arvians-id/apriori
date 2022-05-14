package service

import (
	"apriori/model"
	"apriori/repository"
	"context"
	"database/sql"
)

type UserService interface {
	FindAll(ctx context.Context) ([]model.GetUserResponse, error)
	FindById(ctx context.Context, userId uint64) (model.GetUserResponse, error)
	Create(ctx context.Context, request model.CreateUserRequest) (model.GetUserResponse, error)
	Update(ctx context.Context, request model.UpdateUserRequest) (model.GetUserResponse, error)
	Delete(ctx context.Context, userId uint64) error
}

type userService struct {
	UserRepository repository.UserRepository
	DB             *sql.DB
}

func NewUserService(userRepository repository.UserRepository, db *sql.DB) UserService {
	return &userService{
		UserRepository: userRepository,
		DB:             db,
	}
}

func (service *userService) FindAll(ctx context.Context) ([]model.GetUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (service *userService) FindById(ctx context.Context, userId uint64) (model.GetUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (service *userService) Create(ctx context.Context, request model.CreateUserRequest) (model.GetUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (service *userService) Update(ctx context.Context, request model.UpdateUserRequest) (model.GetUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (service *userService) Delete(ctx context.Context, userId uint64) error {
	//TODO implement me
	panic("implement me")
}
