package cache

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/arvians-id/apriori/cmd/library/aws"
	redisLibrary "github.com/arvians-id/apriori/cmd/library/redis"
	"github.com/arvians-id/apriori/internal/http/presenter/request"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/arvians-id/apriori/internal/repository"
	"github.com/arvians-id/apriori/internal/service"
	"github.com/arvians-id/apriori/util"
	"github.com/go-redis/redis/v8"
	"log"
	"strings"
	"time"
)

type ProductCacheServiceImpl struct {
	ProductRepository repository.ProductRepository
	AprioriRepository repository.AprioriRepository
	Redis             redisLibrary.Redis
	StorageS3         aws.StorageS3
	DB                *sql.DB
}

func NewProductCacheService(
	productRepository *repository.ProductRepository,
	aprioriRepository *repository.AprioriRepository,
	redis *redisLibrary.Redis,
	storageS3 *aws.StorageS3,
	db *sql.DB,
) service.ProductService {
	return &ProductCacheServiceImpl{
		ProductRepository: *productRepository,
		AprioriRepository: *aprioriRepository,
		Redis:             *redis,
		StorageS3:         *storageS3,
		DB:                db,
	}
}

func (cache *ProductCacheServiceImpl) FindAllByAdmin(ctx context.Context) ([]*model.Product, error) {
	tx, err := cache.DB.Begin()
	if err != nil {
		log.Println("[ProductCacheService][FindAllByAdmin] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	productsCache, err := cache.Redis.Get(ctx, "products:admin")
	if err != redis.Nil {
		var products []*model.Product
		err = json.Unmarshal(productsCache, &products)
		if err != nil {
			log.Println("[ProductCacheService][FindAllByAdmin] unable to unmarshal json, err: ", err.Error())
			return nil, err
		}

		return products, nil
	}

	products, err := cache.ProductRepository.FindAllByAdmin(ctx, tx)
	if err != nil {
		log.Println("[ProductCacheService][FindAllByAdmin][FindAllByAdmin] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = cache.Redis.Set(ctx, "products:admin", products)
	if err != nil {
		log.Println("[ProductCacheService][FindAllByAdmin][Set] unable to set value to redis cache, err: ", err.Error())
		return nil, err
	}

	return products, nil
}

func (cache *ProductCacheServiceImpl) FindAll(ctx context.Context, search string, category string) ([]*model.Product, error) {
	tx, err := cache.DB.Begin()
	if err != nil {
		log.Println("[ProductCacheService][FindAll] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	productsCache, err := cache.Redis.Get(ctx, "products:all")
	if err != redis.Nil {
		var products []*model.Product
		err = json.Unmarshal(productsCache, &products)
		if err != nil {
			log.Println("[ProductCacheService][FindAll] unable to unmarshal json, err: ", err.Error())
			return nil, err
		}

		return products, nil
	}

	products, err := cache.ProductRepository.FindAll(ctx, tx, search, category)
	if err != nil {
		log.Println("[ProductCacheService][FindAll][FindAll] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = cache.Redis.Set(ctx, "products:all", products)
	if err != nil {
		log.Println("[ProductCacheService][FindAllByAdmin][Set] unable to set value to redis cache, err: ", err.Error())
		return nil, err
	}

	return products, nil
}

func (cache *ProductCacheServiceImpl) FindAllBySimilarCategory(ctx context.Context, code string) ([]*model.Product, error) {
	tx, err := cache.DB.Begin()
	if err != nil {
		log.Println("[ProductCacheService][FindAllBySimilarCategory] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	productsCache, err := cache.Redis.Get(ctx, "products:similar")
	if err != redis.Nil {
		var products []*model.Product
		err = json.Unmarshal(productsCache, &products)
		if err != nil {
			log.Println("[ProductCacheService][FindAllBySimilarCategory] unable to unmarshal json, err: ", err.Error())
			return nil, err
		}

		return products, nil
	}

	product, err := cache.ProductRepository.FindByCode(ctx, tx, code)
	if err != nil {
		log.Println("[ProductCacheService][FindAllBySimilarCategory][FindByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	categoryArray := strings.Split(product.Category, ", ")
	categoryString := strings.Join(categoryArray, "|")
	productCategories, err := cache.ProductRepository.FindAllBySimilarCategory(ctx, tx, categoryString)
	if err != nil {
		log.Println("[ProductCacheService][FindAllBySimilarCategory][FindAllBySimilarCategory] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var productResponses []*model.Product
	for _, productCategory := range productCategories {
		if productCategory.Code != code {
			productResponses = append(productResponses, productCategory)
		}
	}

	err = cache.Redis.Set(ctx, "products:similar", productResponses)
	if err != nil {
		log.Println("[ProductCacheService][FindAllBySimilarCategory][Set] unable to set value to redis cache, err: ", err.Error())
		return nil, err
	}

	return productResponses, nil
}

func (cache *ProductCacheServiceImpl) FindAllRecommendation(ctx context.Context, code string) ([]*model.ProductRecommendation, error) {
	tx, err := cache.DB.Begin()
	if err != nil {
		log.Println("[ProductCacheService][FindAllRecommendation] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	product, err := cache.ProductRepository.FindByCode(ctx, tx, code)
	if err != nil {
		log.Println("[ProductCacheService][FindAllRecommendation][FindByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	apriories, err := cache.AprioriRepository.FindAllByActive(ctx, tx)
	if err != nil {
		log.Println("[ProductCacheService][FindAllRecommendation][FindAllByActive] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var productResponses []*model.ProductRecommendation
	for _, apriori := range apriories {
		productNames := strings.Split(apriori.Item, ",")
		var exists bool
		for _, productName := range productNames {
			if strings.ToLower(product.Name) == strings.TrimSpace(productName) {
				exists = true
			}
		}

		var totalPrice int
		if exists {
			for _, productName := range productNames {
				productByName, _ := cache.ProductRepository.FindByName(ctx, tx, util.UpperWords(productName))
				totalPrice += productByName.Price
			}

			productResponses = append(productResponses, &model.ProductRecommendation{
				AprioriId:         apriori.IdApriori,
				AprioriCode:       apriori.Code,
				AprioriItem:       apriori.Item,
				AprioriDiscount:   apriori.Discount,
				ProductTotalPrice: totalPrice,
				PriceDiscount:     totalPrice - (totalPrice * int(apriori.Discount) / 100),
				AprioriImage:      apriori.Image,
			})
		}
	}

	return productResponses, nil
}

func (cache *ProductCacheServiceImpl) FindByCode(ctx context.Context, code string) (*model.Product, error) {
	tx, err := cache.DB.Begin()
	if err != nil {
		log.Println("[ProductCacheService][FindByCode] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	key := fmt.Sprintf("product:%s", code)
	productCache, err := cache.Redis.Get(ctx, key)
	if err != redis.Nil {
		var product model.Product
		err = json.Unmarshal(productCache, &product)
		if err != nil {
			log.Println("[ProductCacheService][FindByCode] unable to unmarshal json, err: ", err.Error())
			return nil, err
		}

		return &product, nil
	}

	productResponse, err := cache.ProductRepository.FindByCode(ctx, tx, code)
	if err != nil {
		log.Println("[ProductCacheService][FindByCode][FindByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = cache.Redis.Set(ctx, key, productResponse)
	if err != nil {
		log.Println("[ProductCacheService][FindByCode][Set] unable to set value to redis cache, err: ", err.Error())
		return nil, err
	}

	return productResponse, nil
}

func (cache *ProductCacheServiceImpl) Create(ctx context.Context, request *request.CreateProductRequest) (*model.Product, error) {
	tx, err := cache.DB.Begin()
	if err != nil {
		log.Println("[ProductCacheService][Create] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[ProductCacheService][Create] problem in parsing to time, err: ", err.Error())
		return nil, err
	}
	if request.Image == "" {
		request.Image = "no-image.png"
	}

	productRequest := model.Product{
		Code:        request.Code,
		Name:        util.UpperWords(request.Name),
		Description: &request.Description,
		Price:       request.Price,
		Image:       &request.Image,
		Category:    util.UpperWords(request.Category),
		IsEmpty:     false,
		Mass:        request.Mass,
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	}

	productResponse, err := cache.ProductRepository.Create(ctx, tx, &productRequest)
	if err != nil {
		log.Println("[ProductCacheService][Create][Create] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = cache.Redis.Del(ctx, "products:all", "products:admin", "products:similar")
	if err != nil {
		log.Println("[ProductCacheService][Create][Del] unable to deleting specific key cache, err: ", err.Error())
	}

	return productResponse, nil
}

func (cache *ProductCacheServiceImpl) Update(ctx context.Context, request *request.UpdateProductRequest) (*model.Product, error) {
	tx, err := cache.DB.Begin()
	if err != nil {
		log.Println("[ProductCacheService][Update] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	product, err := cache.ProductRepository.FindByCode(ctx, tx, request.Code)
	if err != nil {
		log.Println("[ProductCacheService][Update][FindByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[ProductCacheService][Update] problem in parsing to time, err: ", err.Error())
		return nil, err
	}

	product.Name = util.UpperWords(request.Name)
	product.Description = &request.Description
	product.Price = request.Price
	product.Category = util.UpperWords(request.Category)
	product.IsEmpty = request.IsEmpty
	product.Mass = request.Mass
	product.UpdatedAt = timeNow
	if request.Image != "" {
		_ = cache.StorageS3.DeleteFromAWS(*product.Image)
		product.Image = &request.Image
	}

	productResponse, err := cache.ProductRepository.Update(ctx, tx, product)
	if err != nil {
		log.Println("[ProductCacheService][Update][Update] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = cache.Redis.Del(ctx, "products:all", "products:admin", "products:similar", fmt.Sprintf("product:%s", request.Code))
	if err != nil {
		log.Println("[ProductCacheService][Update][Del] unable to deleting specific key cache, err: ", err.Error())
	}

	return productResponse, nil
}

func (cache *ProductCacheServiceImpl) Delete(ctx context.Context, code string) (err error) {
	tx, err := cache.DB.Begin()
	if err != nil {
		log.Println("[ProductCacheService][Delete] problem in db transaction, err: ", err.Error())
		return err
	}
	defer util.CommitOrRollback(tx)

	product, err := cache.ProductRepository.FindByCode(ctx, tx, code)
	if err != nil {
		log.Println("[ProductCacheService][Delete][FindByCode] problem in getting from repository, err: ", err.Error())
		return err
	}

	err = cache.ProductRepository.Delete(ctx, tx, product.Code)
	if err != nil {
		log.Println("[ProductCacheService][Delete][Delete] problem in getting from repository, err: ", err.Error())
		return err
	}

	_ = cache.StorageS3.DeleteFromAWS(*product.Image)

	err = cache.Redis.Del(ctx, "products:all", "products:admin", "products:similar", fmt.Sprintf("product:%s", code))
	if err != nil {
		log.Println("[ProductCacheService][Delete][Del] unable to deleting specific key cache, err: ", err.Error())
	}

	return nil
}
