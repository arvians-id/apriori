package service

import (
	"apriori/config"
	"apriori/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type CacheService interface {
	GetClient(ctx context.Context) (*redis.Client, error)
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}) error
	Del(ctx context.Context, key ...string) error
	FlushDB(ctx context.Context) error
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

func (cache *cacheService) GetClient(ctx context.Context) (*redis.Client, error) {
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

func (cache *cacheService) Get(ctx context.Context, key string) (string, error) {
	rdb, err := cache.GetClient(ctx)
	if err != nil {
		return "", err
	}

	value, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (cache *cacheService) Set(ctx context.Context, key string, value interface{}) error {
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

func (cache *cacheService) Del(ctx context.Context, key ...string) error {
	rdb, err := cache.GetClient(ctx)
	if err != nil {
		return err
	}

	for _, k := range key {
		err = rdb.Del(ctx, k).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (cache *cacheService) FlushDB(ctx context.Context) error {
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
