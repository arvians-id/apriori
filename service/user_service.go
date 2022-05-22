package service

import (
	"apriori/entity"
	"apriori/helper"
	"apriori/model"
	"apriori/repository"
	"context"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService interface {
	FindAll(ctx context.Context) ([]model.GetUserResponse, error)
	FindById(ctx context.Context, userId uint64) (model.GetUserResponse, error)
	Create(ctx context.Context, request model.CreateUserRequest) error
	Update(ctx context.Context, request model.UpdateUserRequest) error
	Delete(ctx context.Context, userId uint64) error
}

type userService struct {
	UserRepository repository.UserRepository
	DB             *sql.DB
	date           string
}

func NewUserService(userRepository *repository.UserRepository, db *sql.DB) UserService {
	return &userService{
		UserRepository: *userRepository,
		DB:             db,
		date:           "2006-01-02 15:04:05",
	}
}

func (service *userService) FindAll(ctx context.Context) ([]model.GetUserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []model.GetUserResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	users, err := service.UserRepository.FindAll(ctx, tx)
	if err != nil {
		return []model.GetUserResponse{}, err
	}

	var userResponse []model.GetUserResponse
	for _, user := range users {
		userResponse = append(userResponse, helper.ToUserResponse(user))
	}

	return userResponse, nil
}

func (service *userService) FindById(ctx context.Context, userId uint64) (model.GetUserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetUserResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, userId)
	if err != nil {
		return model.GetUserResponse{}, err
	}

	return helper.ToUserResponse(user), nil
}

func (service *userService) Create(ctx context.Context, request model.CreateUserRequest) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	createdAt, err := time.Parse(service.date, time.Now().Format(service.date))
	if err != nil {
		return err
	}
	updatedAt, err := time.Parse(service.date, time.Now().Format(service.date))
	if err != nil {
		return err
	}

	createUser := entity.User{
		Name:      request.Name,
		Email:     request.Email,
		Password:  string(password),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	err = service.UserRepository.Create(ctx, tx, createUser)
	if err != nil {
		return err
	}

	return nil
}

func (service *userService) Update(ctx context.Context, request model.UpdateUserRequest) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	getUser, err := service.UserRepository.FindById(ctx, tx, request.IdUser)
	if err != nil {
		return err
	}

	newPassword := getUser.Password
	if request.Password != "" {
		password, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		newPassword = string(password)
	}

	updatedAt, err := time.Parse(service.date, time.Now().Format(service.date))
	if err != nil {
		return err
	}

	getUser.Name = request.Name
	getUser.Email = request.Email
	getUser.Password = newPassword
	getUser.UpdatedAt = updatedAt

	err = service.UserRepository.Update(ctx, tx, getUser)
	if err != nil {
		return err
	}

	return nil
}

func (service *userService) Delete(ctx context.Context, userId uint64) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	getUser, err := service.UserRepository.FindById(ctx, tx, userId)
	if err != nil {
		return err
	}

	err = service.UserRepository.Delete(ctx, tx, getUser.IdUser)
	if err != nil {
		return err
	}

	return nil
}
