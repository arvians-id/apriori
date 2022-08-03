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

type ProductCache interface {
	SingleGet(ctx context.Context, key string) (model.GetProductResponse, error)
	Get(ctx context.Context, key string) ([]model.GetProductResponse, error)
	Set(ctx context.Context, key string, value []model.GetProductResponse) error
	SingleSet(ctx context.Context, key string, value model.GetProductResponse) error
}

type productCache struct {
	Addr     string
	Password string
	DB       int
}

func NewProductCache(configuration config.Config) ProductCache {
	host := fmt.Sprintf("%s:%s", configuration.Get("REDIS_HOST"), configuration.Get("REDIS_PORT"))
	db := utils.StrToInt(configuration.Get("REDIS_DB"))

	return &productCache{
		Addr:     host,
		Password: configuration.Get("REDIS_PASSWORD"),
		DB:       db,
	}
}

func (cache *productCache) GetClient(ctx context.Context) (*redis.Client, error) {
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

func (cache *productCache) Get(ctx context.Context, key string) ([]model.GetProductResponse, error) {
	rdb, err := cache.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	value, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var product []model.GetProductResponse
	err = json.Unmarshal(bytes.NewBufferString(value).Bytes(), &product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (cache *productCache) SingleGet(ctx context.Context, key string) (model.GetProductResponse, error) {
	rdb, err := cache.GetClient(ctx)
	if err != nil {
		return model.GetProductResponse{}, err
	}

	value, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return model.GetProductResponse{}, err
	}

	var product model.GetProductResponse
	err = json.Unmarshal(bytes.NewBufferString(value).Bytes(), &product)
	if err != nil {
		return model.GetProductResponse{}, err
	}

	return product, nil
}

func (cache *productCache) Set(ctx context.Context, key string, value []model.GetProductResponse) error {
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

func (cache *productCache) SingleSet(ctx context.Context, key string, value model.GetProductResponse) error {
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
