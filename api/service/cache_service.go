package service

import (
	"apriori/config"
	"apriori/model"
	"apriori/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type CacheService interface {
	Get(ctx context.Context, key string) ([]model.GetTransactionResponse, error)
	Set(ctx context.Context, key string, value []model.GetTransactionResponse) error
}

type cacheService struct {
	Addr     string
	Password string
	DB       int
}

func NewCacheService(configuration config.Config) CacheService {
	host := fmt.Sprintf("%s:%s", configuration.Get("REDIS_HOST"), configuration.Get("REDIS_PORT"))
	db := utils.StrToInt(configuration.Get("REDIS_DB"))

	return &cacheService{
		Addr:     host,
		Password: configuration.Get("REDIS_PASSWORD"),
		DB:       db,
	}
}

func (service *cacheService) GetClient(ctx context.Context) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     service.Addr,
		Password: service.Password,
		DB:       service.DB,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return rdb, nil
}

func (service *cacheService) Get(ctx context.Context, key string) ([]model.GetTransactionResponse, error) {
	rdb, err := service.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	value, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var transaction []model.GetTransactionResponse
	err = json.Unmarshal(bytes.NewBufferString(value).Bytes(), &transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (service *cacheService) Set(ctx context.Context, key string, value []model.GetTransactionResponse) error {
	rdb, err := service.GetClient(ctx)
	if err != nil {
		return err
	}

	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = rdb.Set(ctx, key, bytes.NewBuffer(b).Bytes(), time.Duration(30)*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}
