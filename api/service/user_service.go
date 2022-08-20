package service

import (
	"apriori/entity"
	"apriori/model"
	"apriori/repository"
	"apriori/utils"
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService interface {
	FindAll(ctx context.Context) ([]model.GetUserResponse, error)
	FindById(ctx context.Context, id int) (model.GetUserResponse, error)
	FindByEmail(ctx context.Context, request model.GetUserCredentialRequest) (model.GetUserResponse, error)
	Create(ctx context.Context, request model.CreateUserRequest) (model.GetUserResponse, error)
	Update(ctx context.Context, request model.UpdateUserRequest) (model.GetUserResponse, error)
	Delete(ctx context.Context, id int) error
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

	var userResponses []model.GetUserResponse
	for _, user := range users {
		userResponses = append(userResponses, utils.ToUserResponse(user))
	}

	return userResponses, nil
}

func (service *userService) FindById(ctx context.Context, id int) (model.GetUserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetUserResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	userResponse, err := service.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		return model.GetUserResponse{}, err
	}

	return utils.ToUserResponse(userResponse), nil
}

func (service *userService) FindByEmail(ctx context.Context, request model.GetUserCredentialRequest) (model.GetUserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetUserResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	userResponse, err := service.UserRepository.FindByEmail(ctx, tx, request.Email)
	if err != nil {
		return model.GetUserResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userResponse.Password), []byte(request.Password))
	if err != nil {
		return model.GetUserResponse{}, errors.New("wrong password")
	}

	return utils.ToUserResponse(userResponse), nil
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

	userRequest := entity.User{
		Role:      2,
		Name:      request.Name,
		Email:     request.Email,
		Address:   request.Address,
		Phone:     request.Phone,
		Password:  string(password),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}
	userResponse, err := service.UserRepository.Create(ctx, tx, userRequest)
	if err != nil {
		return model.GetUserResponse{}, err
	}

	return utils.ToUserResponse(userResponse), nil
}

func (service *userService) Update(ctx context.Context, request model.UpdateUserRequest) (model.GetUserResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetUserResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, request.IdUser)
	if err != nil {
		return model.GetUserResponse{}, err
	}

	newPassword := user.Password
	if request.Password != "" {
		password, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return model.GetUserResponse{}, err
		}

		newPassword = string(password)
	}

	timeNow, err := time.Parse(service.date, time.Now().Format(service.date))
	if err != nil {
		return model.GetUserResponse{}, err
	}

	user.Name = request.Name
	user.Email = request.Email
	user.Address = request.Address
	user.Phone = request.Phone
	user.Password = newPassword
	user.UpdatedAt = timeNow

	userResponse, err := service.UserRepository.Update(ctx, tx, user)
	if err != nil {
		return model.GetUserResponse{}, err
	}

	return utils.ToUserResponse(userResponse), nil
}

func (service *userService) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}

	err = service.UserRepository.Delete(ctx, tx, user.IdUser)
	if err != nil {
		return err
	}

	return nil
}
