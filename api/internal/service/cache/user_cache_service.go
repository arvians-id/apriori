package cache

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	redisLibrary "github.com/arvians-id/apriori/cmd/library/redis"
	"github.com/arvians-id/apriori/internal/http/presenter/request"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/arvians-id/apriori/internal/repository"
	"github.com/arvians-id/apriori/internal/service"
	"github.com/arvians-id/apriori/util"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type UserCacheServiceImpl struct {
	UserRepository repository.UserRepository
	Redis          redisLibrary.Redis
	DB             *sql.DB
}

func NewUserCacheService(userRepository *repository.UserRepository, redis *redisLibrary.Redis, db *sql.DB) service.UserService {
	return &UserCacheServiceImpl{
		UserRepository: *userRepository,
		Redis:          *redis,
		DB:             db,
	}
}

func (cache *UserCacheServiceImpl) FindAll(ctx context.Context) ([]*model.User, error) {
	var tx *sql.Tx
	if cache.DB != nil {
		transaction, err := cache.DB.Begin()
		if err != nil {
			log.Println("[UserCacheService][FindAll] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	usersCache, err := cache.Redis.Get(ctx, "users")
	if err != redis.Nil {
		var users []*model.User
		err = json.Unmarshal(usersCache, &users)
		if err != nil {
			log.Println("[UserCacheService][FindAll] unable to unmarshal json, err: ", err.Error())
			return nil, err
		}

		return users, nil
	}

	users, err := cache.UserRepository.FindAll(ctx, tx)
	if err != nil {
		log.Println("[UserCacheService][FindAll][FindAll] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = cache.Redis.Set(ctx, "users", users)
	if err != nil {
		log.Println("[UserCacheService][FindAll][Set] unable to set value to redis cache, err: ", err.Error())
		return nil, err
	}

	return users, nil
}

func (cache *UserCacheServiceImpl) FindById(ctx context.Context, id int) (*model.User, error) {
	var tx *sql.Tx
	if cache.DB != nil {
		transaction, err := cache.DB.Begin()
		if err != nil {
			log.Println("[UserCacheService][FindById] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	key := fmt.Sprintf("user:%d", id)
	userCache, err := cache.Redis.Get(ctx, key)
	if err != redis.Nil {
		var user model.User
		err = json.Unmarshal(userCache, &user)
		if err != nil {
			log.Println("[UserCacheService][FindById] unable to unmarshal json, err: ", err.Error())
			return nil, err
		}

		return &user, nil
	}

	user, err := cache.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		log.Println("[UserCacheService][FindById][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = cache.Redis.Set(ctx, key, user)
	if err != nil {
		log.Println("[UserCacheService][FindById][Set] unable to set value to redis cache, err: ", err.Error())
		return nil, err
	}

	return user, nil
}

func (cache *UserCacheServiceImpl) FindByEmail(ctx context.Context, request *request.GetUserCredentialRequest) (*model.User, error) {
	var tx *sql.Tx
	if cache.DB != nil {
		transaction, err := cache.DB.Begin()
		if err != nil {
			log.Println("[UserCacheService][FindByEmail] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	key := fmt.Sprintf("user:%s", request.Email)
	userCache, err := cache.Redis.Get(ctx, key)
	if err != redis.Nil {
		var user model.User
		err = json.Unmarshal(userCache, &user)
		if err != nil {
			log.Println("[UserCacheService][FindByEmail] unable to unmarshal json, err: ", err.Error())
			return nil, err
		}

		return &user, nil
	}

	user, err := cache.UserRepository.FindByEmail(ctx, tx, request.Email)
	if err != nil {
		log.Println("[UserCacheService][FindByEmail][FindByEmail] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		log.Println("[UserCacheService][FindByEmail][CompareHashAndPassword] problem in comparing password, err: ", err.Error())
		return nil, errors.New("wrong password")
	}

	err = cache.Redis.Set(ctx, key, user)
	if err != nil {
		log.Println("[UserCacheService][FindByEmail][Set] unable to set value to redis cache, err: ", err.Error())
		return nil, err
	}

	return user, nil
}

func (cache *UserCacheServiceImpl) Create(ctx context.Context, request *request.CreateUserRequest) (*model.User, error) {
	var tx *sql.Tx
	if cache.DB != nil {
		transaction, err := cache.DB.Begin()
		if err != nil {
			log.Println("[UserCacheService][Create] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[UserCacheService][Create][GenerateFromPassword] problem in generating password, err: ", err.Error())
		return nil, err
	}

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[UserCacheService][Create] problem in parsing time, err: ", err.Error())
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
	user, err := cache.UserRepository.Create(ctx, tx, &userRequest)
	if err != nil {
		log.Println("[UserCacheService][Create][Create] problem in creating user, err: ", err.Error())
		return nil, err
	}

	err = cache.Redis.Del(ctx, "users")
	if err != nil {
		log.Println("[UserCacheService][Create][Del] unable to deleting specific key cache, err: ", err.Error())
	}

	return user, nil
}

func (cache *UserCacheServiceImpl) Update(ctx context.Context, request *request.UpdateUserRequest) (*model.User, error) {
	var tx *sql.Tx
	if cache.DB != nil {
		transaction, err := cache.DB.Begin()
		if err != nil {
			log.Println("[UserCacheService][Update] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	user, err := cache.UserRepository.FindById(ctx, tx, request.IdUser)
	if err != nil {
		log.Println("[UserCacheService][Update][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	newPassword := user.Password
	if request.Password != "" {
		password, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("[UserCacheService][Update] problem in generating password, err: ", err.Error())
			return nil, err
		}

		newPassword = string(password)
	}

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[UserCacheService][Update] problem in parsing time, err: ", err.Error())
		return nil, err
	}

	user.Name = request.Name
	user.Email = request.Email
	user.Address = request.Address
	user.Phone = request.Phone
	user.Password = newPassword
	user.UpdatedAt = timeNow

	_, err = cache.UserRepository.Update(ctx, tx, user)
	if err != nil {
		log.Println("[UserCacheService][Update][Update] problem in updating user, err: ", err.Error())
		return nil, err
	}

	key1 := fmt.Sprintf("user:%s", request.Email)
	key2 := fmt.Sprintf("user:%d", request.IdUser)
	err = cache.Redis.Del(ctx, "users", key1, key2)
	if err != nil {
		log.Println("[UserCacheService][Update][Del] unable to deleting specific key cache, err: ", err.Error())
	}

	return user, nil
}

func (cache *UserCacheServiceImpl) Delete(ctx context.Context, id int) error {
	var tx *sql.Tx
	if cache.DB != nil {
		transaction, err := cache.DB.Begin()
		if err != nil {
			log.Println("[UserCacheService][Delete] problem in db transaction, err: ", err.Error())
			return err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	user, err := cache.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		log.Println("[UserCacheService][Delete][FindById] problem in getting from repository, err: ", err.Error())
		return err
	}

	err = cache.UserRepository.Delete(ctx, tx, user.IdUser)
	if err != nil {
		log.Println("[UserCacheService][Delete][Delete] problem in deleting user, err: ", err.Error())
		return err
	}

	err = cache.Redis.Del(ctx, "users")
	if err != nil {
		log.Println("[UserCacheService][Delete][Del] unable to deleting specific key cache, err: ", err.Error())
	}

	return nil
}
