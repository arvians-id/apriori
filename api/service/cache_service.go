package service

import (
	"apriori/config"
	"apriori/helper"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

func (cache *CacheServiceImpl) Get(ctx context.Context, key string) (string, error) {
	rdb, err := cache.GetClient()
	if err != nil {
		return "", err
	}

	value, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
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

func (cache *CacheServiceImpl) Subscribe(ctx context.Context) (string, error) {
	rdb, err := cache.GetClient()
	if err != nil {
		return "", err
	}

	subscriber := rdb.Subscribe(ctx, "test")
	defer subscriber.Close()
	var str string
	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			break
		}

		str = msg.Payload
	}

	return str, nil
}

func (cache *CacheServiceImpl) Publish(ctx context.Context) error {
	rdb, err := cache.GetClient()
	if err != nil {
		return err
	}

	err = rdb.Publish(ctx, "test", "Redis test").Err()
	if err != nil {
		return err
	}

	return nil
}
