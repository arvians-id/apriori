package redis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvians-id/apriori/cmd/config"
	"github.com/arvians-id/apriori/util"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type Redis struct {
	Addr     string
	Password string
	DB       int
	Expired  int
}

func NewCacheService(configuration config.Config) *Redis {
	host := fmt.Sprintf("%s:%s", configuration.Get("REDIS_HOST"), configuration.Get("REDIS_PORT"))
	db := util.StrToInt(configuration.Get("REDIS_DB"))

	return &Redis{
		Addr:     host,
		Password: configuration.Get("REDIS_PASSWORD"),
		DB:       db,
		Expired:  util.StrToInt(configuration.Get("REDIS_EXPIRED")),
	}
}

func (cache *Redis) GetClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cache.Addr,
		Password: cache.Password,
		DB:       cache.DB,
	})

	return rdb, nil
}

func (cache *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	rdb, err := cache.GetClient()
	if err != nil {
		log.Println("[Redis][Get][GetClient] problem in connecting to redis, err: ", err.Error())
		return nil, err
	}

	value, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		log.Println("[Redis][Get][Get] problem in getting cache from redis, err: ", err.Error())
		return nil, err
	}

	return value, nil
}

func (cache *Redis) Set(ctx context.Context, key string, value interface{}) error {
	rdb, err := cache.GetClient()
	if err != nil {
		log.Println("[Redis][Set][GetClient] problem in connecting to redis, err: ", err.Error())
		return err
	}

	b, err := json.Marshal(value)
	if err != nil {
		log.Println("[Redis][Set] unable to marshal json, err: ", err.Error())
		return err
	}

	err = rdb.Set(ctx, key, bytes.NewBuffer(b).Bytes(), time.Duration(cache.Expired)*time.Minute).Err()
	if err != nil {
		log.Println("[Redis][Set][Set] problem in setting cache from redis, err: ", err.Error())
		return err
	}

	return nil
}

func (cache *Redis) Del(ctx context.Context, key ...string) error {
	rdb, err := cache.GetClient()
	if err != nil {
		log.Println("[Redis][Del][GetClient] problem in connecting to redis, err: ", err.Error())
		return err
	}

	for _, k := range key {
		err = rdb.Del(ctx, k).Err()
		if err != nil {
			log.Println("[Redis][Del][Del] problem in deleting cache from redis, err: ", err.Error())
			return err
		}
	}

	return nil
}

func (cache *Redis) FlushDB(ctx context.Context) error {
	rdb, err := cache.GetClient()
	if err != nil {
		log.Println("[Redis][FlushDB][GetClient] problem in connecting to redis, err: ", err.Error())
		return err
	}

	err = rdb.FlushDB(ctx).Err()
	if err != nil {
		log.Println("[Redis][FlushDB][FlushDB] problem in flushing db cache from redis, err: ", err.Error())
		return err
	}

	return nil
}
