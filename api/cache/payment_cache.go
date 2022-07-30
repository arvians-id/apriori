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

type PaymentCache interface {
	SingleGet(ctx context.Context, key string) (model.GetPaymentNullableResponse, error)
	Get(ctx context.Context, key string) ([]model.GetPaymentNullableResponse, error)
	Set(ctx context.Context, key string, value []model.GetPaymentNullableResponse) error
	SingleSet(ctx context.Context, key string, value model.GetPaymentNullableResponse) error
	RecoverCache(ctx context.Context, key string, userId int) error
}

type paymentCache struct {
	Addr     string
	Password string
	DB       int
}

func NewPaymentCache(configuration config.Config) PaymentCache {
	host := fmt.Sprintf("%s:%s", configuration.Get("REDIS_HOST"), configuration.Get("REDIS_PORT"))
	db := utils.StrToInt(configuration.Get("REDIS_DB"))

	return &paymentCache{
		Addr:     host,
		Password: configuration.Get("REDIS_PASSWORD"),
		DB:       db,
	}
}

func (cache *paymentCache) GetClient(ctx context.Context) (*redis.Client, error) {
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

func (cache *paymentCache) Get(ctx context.Context, key string) ([]model.GetPaymentNullableResponse, error) {
	rdb, err := cache.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	value, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var userOrder []model.GetPaymentNullableResponse
	err = json.Unmarshal(bytes.NewBufferString(value).Bytes(), &userOrder)
	if err != nil {
		return nil, err
	}

	return userOrder, nil
}

func (cache *paymentCache) SingleGet(ctx context.Context, key string) (model.GetPaymentNullableResponse, error) {
	rdb, err := cache.GetClient(ctx)
	if err != nil {
		return model.GetPaymentNullableResponse{}, err
	}

	value, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return model.GetPaymentNullableResponse{}, err
	}

	var userOrder model.GetPaymentNullableResponse
	err = json.Unmarshal(bytes.NewBufferString(value).Bytes(), &userOrder)
	if err != nil {
		return model.GetPaymentNullableResponse{}, err
	}

	return userOrder, nil
}

func (cache *paymentCache) Set(ctx context.Context, key string, value []model.GetPaymentNullableResponse) error {
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

func (cache *paymentCache) SingleSet(ctx context.Context, key string, value model.GetPaymentNullableResponse) error {
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

func (cache *paymentCache) RecoverCache(ctx context.Context, key string, userId int) error {
	return nil
}
