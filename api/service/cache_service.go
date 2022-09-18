package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvians-id/apriori/config"
	"github.com/arvians-id/apriori/helper"
	"github.com/go-redis/redis/v8"
	"time"
)

type CacheServiceImpl struct {
	Addr     string
	Password string
	DB       int
	Expired  int
}

func NewCacheService(configuration config.Config) CacheService {
	host := fmt.Sprintf("%s:%s", configuration.Get("REDIS_HOST"), configuration.Get("REDIS_PORT"))
	db := helper.StrToInt(configuration.Get("REDIS_DB"))

	return &CacheServiceImpl{
		Addr:     host,
		Password: configuration.Get("REDIS_PASSWORD"),
		DB:       db,
		Expired:  helper.StrToInt(configuration.Get("REDIS_EXPIRED")),
	}
}

func (cache *CacheServiceImpl) GetClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cache.Addr,
		Password: cache.Password,
		DB:       cache.DB,
	})

	return rdb, nil
}

func (cache *CacheServiceImpl) Get(ctx context.Context, key string) ([]byte, error) {
	rdb, err := cache.GetClient()
	if err != nil {
		return nil, err
	}

	value, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (cache *CacheServiceImpl) Set(ctx context.Context, key string, value interface{}) error {
	rdb, err := cache.GetClient()
	if err != nil {
		return err
	}

	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = rdb.Set(ctx, key, bytes.NewBuffer(b).Bytes(), time.Duration(cache.Expired)*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (cache *CacheServiceImpl) Del(ctx context.Context, key ...string) error {
	rdb, err := cache.GetClient()
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

func (cache *CacheServiceImpl) FlushDB(ctx context.Context) error {
	rdb, err := cache.GetClient()
	if err != nil {
		return err
	}

	err = rdb.FlushDB(ctx).Err()
	if err != nil {
		return err
	}

	return nil
}
