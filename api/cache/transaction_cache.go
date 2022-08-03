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

type TransactionCache interface {
	Get(ctx context.Context, key string) ([]model.GetTransactionResponse, error)
	Set(ctx context.Context, key string, value []model.GetTransactionResponse) error
	FlushDB(ctx context.Context) error
}

type transactionCache struct {
	Addr     string
	Password string
	DB       int
}

func NewTransactionCache(configuration config.Config) TransactionCache {
	host := fmt.Sprintf("%s:%s", configuration.Get("REDIS_HOST"), configuration.Get("REDIS_PORT"))
	db := utils.StrToInt(configuration.Get("REDIS_DB"))

	return &transactionCache{
		Addr:     host,
		Password: configuration.Get("REDIS_PASSWORD"),
		DB:       db,
	}
}

func (cache *transactionCache) GetClient(ctx context.Context) (*redis.Client, error) {
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

func (cache *transactionCache) Get(ctx context.Context, key string) ([]model.GetTransactionResponse, error) {
	rdb, err := cache.GetClient(ctx)
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

func (cache *transactionCache) Set(ctx context.Context, key string, value []model.GetTransactionResponse) error {
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

func (cache *transactionCache) FlushDB(ctx context.Context) error {
	rdb, err := cache.GetClient(ctx)
	if err != nil {
		return err
	}

	err = rdb.FlushDB(ctx).Err()
	if err != nil {
		return err
	}

	return nil
}
