package cache

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	redisLibrary "github.com/arvians-id/apriori/cmd/library/redis"
	"github.com/arvians-id/apriori/internal/http/presenter/request"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/arvians-id/apriori/internal/repository"
	"github.com/arvians-id/apriori/internal/service"
	"github.com/arvians-id/apriori/util"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type CategoryCacheServiceImpl struct {
	CategoryRepository repository.CategoryRepository
	Redis              redisLibrary.Redis
	DB                 *sql.DB
}

func NewCategoryCacheService(categoryRepository *repository.CategoryRepository, redis *redisLibrary.Redis, db *sql.DB) service.CategoryService {
	return &CategoryCacheServiceImpl{
		CategoryRepository: *categoryRepository,
		Redis:              *redis,
		DB:                 db,
	}
}

func (cache *CategoryCacheServiceImpl) FindAll(ctx context.Context) ([]*model.Category, error) {
	tx, err := cache.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	categoriesCache, err := cache.Redis.Get(ctx, "categories")
	if err != redis.Nil {
		var categories []*model.Category
		err = json.Unmarshal(categoriesCache, &categories)
		if err != nil {
			return nil, err
		}

		return categories, nil
	}

	categories, err := cache.CategoryRepository.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	err = cache.Redis.Set(ctx, "categories", categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (cache *CategoryCacheServiceImpl) FindById(ctx context.Context, id int) (*model.Category, error) {
	tx, err := cache.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	key := fmt.Sprintf("category:%d", id)
	categoryCache, err := cache.Redis.Get(ctx, key)
	if err != redis.Nil {
		var category model.Category
		err = json.Unmarshal(categoryCache, &category)
		if err != nil {
			return nil, err
		}

		return &category, nil
	}

	category, err := cache.CategoryRepository.FindById(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	err = cache.Redis.Set(ctx, key, category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (cache *CategoryCacheServiceImpl) Create(ctx context.Context, request *request.CreateCategoryRequest) (*model.Category, error) {
	tx, err := cache.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		return nil, err
	}

	categoryRequest := model.Category{
		Name:      util.UpperWords(request.Name),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	category, err := cache.CategoryRepository.Create(ctx, tx, &categoryRequest)
	if err != nil {
		return nil, err
	}

	err = cache.Redis.Del(ctx, "categories")
	if err != nil {
		log.Println(err)
	}

	return category, nil

}

func (cache *CategoryCacheServiceImpl) Update(ctx context.Context, request *request.UpdateCategoryRequest) (*model.Category, error) {
	tx, err := cache.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	category, err := cache.CategoryRepository.FindById(ctx, tx, request.IdCategory)
	if err != nil {
		return nil, err
	}

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		return nil, err
	}
	category.Name = util.UpperWords(request.Name)
	category.UpdatedAt = timeNow

	_, err = cache.CategoryRepository.Update(ctx, tx, category)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("category:%d", category.IdCategory)
	err = cache.Redis.Del(ctx, "categories", key)
	if err != nil {
		log.Println(err)
	}

	return category, nil
}

func (cache *CategoryCacheServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := cache.DB.Begin()
	if err != nil {
		return err
	}
	defer util.CommitOrRollback(tx)

	category, err := cache.CategoryRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}

	err = cache.CategoryRepository.Delete(ctx, tx, category.IdCategory)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("category:%d", category.IdCategory)
	err = cache.Redis.Del(ctx, "categories", key)
	if err != nil {
		log.Println(err)
	}

	return nil
}
