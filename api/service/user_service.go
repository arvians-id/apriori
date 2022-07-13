package service

import (
	"apriori/entity"
	"apriori/model"
	"apriori/repository"
	"apriori/utils"
	"context"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"time"
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
	defer utils.CommitOrRollback(tx)

	users, err := service.UserRepository.FindAll(ctx, tx)
	if err != nil {
		return []model.GetUserResponse{}, err
	}

	var userResponse []model.GetUserResponse
	for _, user := range users {
		userResponse = append(userResponse, utils.ToUserResponse(user))
	}

	return userResponse, nil
}

func (service *userService) FindById(ctx context.Context, userId uint64) (model.GetUserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetUserResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, userId)
	if err != nil {
		return model.GetUserResponse{}, err
	}

	return utils.ToUserResponse(user), nil
}

func (service *userService) Create(ctx context.Context, request model.CreateUserRequest) (model.GetUserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetUserResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.GetUserResponse{}, err
	}

	timeNow, err := time.Parse(service.date, time.Now().Format(service.date))
	if err != nil {
		return model.GetUserResponse{}, err
	}

	createUser := entity.User{
		Role:      2,
		Name:      request.Name,
		Email:     request.Email,
		Address:   request.Address,
		Phone:     request.Phone,
		Password:  string(password),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	user, err := service.UserRepository.Create(ctx, tx, createUser)
	if err != nil {
		return model.GetUserResponse{}, err
	}

	return utils.ToUserResponse(user), nil
}

func (service *userService) Update(ctx context.Context, request model.UpdateUserRequest) (model.GetUserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetUserResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	getUser, err := service.UserRepository.FindById(ctx, tx, request.IdUser)
	if err != nil {
		return model.GetUserResponse{}, err
	}

	newPassword := getUser.Password
	if request.Password != "" {
		password, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return model.GetUserResponse{}, err
		}

		newPassword = string(password)
	}

	updatedAt, err := time.Parse(service.date, time.Now().Format(service.date))
	if err != nil {
		return model.GetUserResponse{}, err
	}

	getUser.Name = request.Name
	getUser.Email = request.Email
	getUser.Address = request.Address
	getUser.Phone = request.Phone
	getUser.Password = newPassword
	getUser.UpdatedAt = updatedAt

	user, err := service.UserRepository.Update(ctx, tx, getUser)
	if err != nil {
		return model.GetUserResponse{}, err
	}

	return utils.ToUserResponse(user), nil
}

func (service *userService) Delete(ctx context.Context, userId uint64) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

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
