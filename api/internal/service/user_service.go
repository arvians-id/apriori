package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/arvians-id/apriori/internal/http/presenter/request"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/arvians-id/apriori/internal/repository"
	"github.com/arvians-id/apriori/util"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	DB             *sql.DB
}

func NewUserService(userRepository *repository.UserRepository, db *sql.DB) UserService {
	return &UserServiceImpl{
		UserRepository: *userRepository,
		DB:             db,
	}
}

func (service *UserServiceImpl) FindAll(ctx context.Context) ([]*model.User, error) {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserService][FindAll] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	users, err := service.UserRepository.FindAll(ctx, tx)
	if err != nil {
		log.Println("[UserService][FindAll][FindAll] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return users, nil
}

func (service *UserServiceImpl) FindById(ctx context.Context, id int) (*model.User, error) {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserService][FindById] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		log.Println("[UserService][FindById][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return user, nil
}

func (service *UserServiceImpl) FindByEmail(ctx context.Context, request *request.GetUserCredentialRequest) (*model.User, error) {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserService][FindByEmail] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	user, err := service.UserRepository.FindByEmail(ctx, tx, request.Email)
	if err != nil {
		log.Println("[UserService][FindByEmail][FindByEmail] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		log.Println("[UserService][FindByEmail] problem in comparing password, err: ", err.Error())
		return nil, errors.New("wrong password")
	}

	return user, nil
}

func (service *UserServiceImpl) Create(ctx context.Context, request *request.CreateUserRequest) (*model.User, error) {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserService][Create] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[UserService][Create] problem in generating password hashed, err: ", err.Error())
		return nil, err
	}

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[UserService][Create] problem in parsing to time, err: ", err.Error())
		return nil, err
	}

	userRequest := model.User{
		Role:      2,
		Name:      request.Name,
		Email:     request.Email,
		Address:   request.Address,
		Phone:     request.Phone,
		Password:  string(password),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}
	user, err := service.UserRepository.Create(ctx, tx, &userRequest)
	if err != nil {
		log.Println("[UserService][Create][Create] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return user, nil
}

func (service *UserServiceImpl) Update(ctx context.Context, request *request.UpdateUserRequest) (*model.User, error) {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserService][Update] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, request.IdUser)
	if err != nil {
		log.Println("[UserService][Update][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	newPassword := user.Password
	if request.Password != "" {
		password, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("[UserService][Update] problem in generating password hashed, err: ", err.Error())
			return nil, err
		}

		newPassword = string(password)
	}

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[UserService][Update] problem in parsing to time, err: ", err.Error())
		return nil, err
	}

	user.Name = request.Name
	user.Email = request.Email
	user.Address = request.Address
	user.Phone = request.Phone
	user.Password = newPassword
	user.UpdatedAt = timeNow

	_, err = service.UserRepository.Update(ctx, tx, user)
	if err != nil {
		log.Println("[UserService][Update][Update] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return user, nil
}

func (service *UserServiceImpl) Delete(ctx context.Context, id int) error {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserService][Delete] problem in db transaction, err: ", err.Error())
			return err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		log.Println("[UserService][Delete][FindById] problem in getting from repository, err: ", err.Error())
		return err
	}

	err = service.UserRepository.Delete(ctx, tx, user.IdUser)
	if err != nil {
		log.Println("[UserService][Delete][Delete] problem in getting from repository, err: ", err.Error())
		return err
	}

	return nil
}
