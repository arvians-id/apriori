package cache

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

type UserOrderCache interface {
	SingleGet(ctx context.Context, key string) (model.GetUserOrderResponse, error)
	Get(ctx context.Context, key string) ([]model.GetUserOrderResponse, error)
	Set(ctx context.Context, key string, value []model.GetUserOrderResponse) error
	SingleSet(ctx context.Context, key string, value model.GetUserOrderResponse) error
}

type userOrderCache struct {
	Addr     string
	Password string
	DB       int
}

func NewUserOrderCache(configuration config.Config) UserOrderCache {
	host := fmt.Sprintf("%s:%s", configuration.Get("REDIS_HOST"), configuration.Get("REDIS_PORT"))
	db := utils.StrToInt(configuration.Get("REDIS_DB"))

	return &userOrderCache{
		Addr:     host,
		Password: configuration.Get("REDIS_PASSWORD"),
		DB:       db,
	}
}

func (cache *userOrderCache) GetClient(ctx context.Context) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cache.Addr,
		Password: cache.Password,
		DB:       cache.DB,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return rdb, nil
}

func (cache *userOrderCache) Get(ctx context.Context, key string) ([]model.GetUserOrderResponse, error) {
	rdb, err := cache.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	value, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var userOrder []model.GetUserOrderResponse
	err = json.Unmarshal(bytes.NewBufferString(value).Bytes(), &userOrder)
	if err != nil {
		return nil, err
	}

	return userOrder, nil
}

func (cache *userOrderCache) SingleGet(ctx context.Context, key string) (model.GetUserOrderResponse, error) {
	rdb, err := cache.GetClient(ctx)
	if err != nil {
		return model.GetUserOrderResponse{}, err
	}

	value, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return model.GetUserOrderResponse{}, err
	}

	var userOrder model.GetUserOrderResponse
	err = json.Unmarshal(bytes.NewBufferString(value).Bytes(), &userOrder)
	if err != nil {
		return model.GetUserOrderResponse{}, err
	}

	return userOrder, nil
}

func (cache *userOrderCache) Set(ctx context.Context, key string, value []model.GetUserOrderResponse) error {
	rdb, err := cache.GetClient(ctx)
	if err != nil {
		return err
	}

	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = rdb.Set(ctx, key, bytes.NewBuffer(b).Bytes(), time.Duration(60)*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

func (cache *userOrderCache) SingleSet(ctx context.Context, key string, value model.GetUserOrderResponse) error {
	rdb, err := cache.GetClient(ctx)
	if err != nil {
		return err
	}

	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = rdb.Set(ctx, key, bytes.NewBuffer(b).Bytes(), time.Duration(60)*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}
