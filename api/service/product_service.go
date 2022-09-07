package service

import (
	"apriori/entity"
	"apriori/helper"
	"apriori/model"
	"apriori/repository"
	"context"
	"database/sql"
	"strings"
	"time"
)

type ProductServiceImpl struct {
	ProductRepository repository.ProductRepository
	AprioriRepository repository.AprioriRepository
	StorageService
	DB *sql.DB
}

func NewProductService(
	productRepository *repository.ProductRepository,
	storageService StorageService,
	aprioriRepository *repository.AprioriRepository,
	db *sql.DB,
) ProductService {
	return &ProductServiceImpl{
		ProductRepository: *productRepository,
		AprioriRepository: *aprioriRepository,
		StorageService:    storageService,
		DB:                db,
	}
}

func (service *ProductServiceImpl) FindAllByAdmin(ctx context.Context) ([]*model.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	products, err := service.ProductRepository.FindAllByAdmin(ctx, tx)
	if err != nil {
		return nil, err
	}

	var productResponses []*model.GetProductResponse
	for _, product := range products {
		productResponses = append(productResponses, product.ToProductResponse())
	}

	return productResponses, nil
}

func (service *ProductServiceImpl) FindAll(ctx context.Context, search string, category string) ([]*model.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	products, err := service.ProductRepository.FindAll(ctx, tx, search, category)
	if err != nil {
		return nil, err
	}

	var productResponses []*model.GetProductResponse
	for _, product := range products {
		productResponses = append(productResponses, product.ToProductResponse())
	}

	return productResponses, nil
}

func (service *ProductServiceImpl) FindAllBySimilarCategory(ctx context.Context, code string) ([]*model.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, code)
	if err != nil {
		return nil, err
	}

	categoryArray := strings.Split(product.Category, ", ")
	categoryString := strings.Join(categoryArray, "|")
	productCategories, err := service.ProductRepository.FindAllBySimilarCategory(ctx, tx, categoryString)
	if err != nil {
		return nil, err
	}

	var productResponses []*model.GetProductResponse
	for _, productCategory := range productCategories {
		if productCategory.Code != code {
			productResponses = append(productResponses, productCategory.ToProductResponse())
		}
	}

	return productResponses, nil
}

func (service *ProductServiceImpl) FindAllRecommendation(ctx context.Context, code string) ([]*model.GetProductRecommendationResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, code)
	if err != nil {
		return nil, err
	}

	apriories, err := service.AprioriRepository.FindAllByActive(ctx, tx)
	if err != nil {
		return nil, err
	}

	var productResponses []*model.GetProductRecommendationResponse
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
				productByName, _ := service.ProductRepository.FindByName(ctx, tx, helper.UpperWords(productName))
				totalPrice += productByName.Price
			}

			productResponses = append(productResponses, &model.GetProductRecommendationResponse{
				AprioriId:          apriori.IdApriori,
				AprioriCode:        apriori.Code,
				AprioriItem:        apriori.Item,
				AprioriDiscount:    apriori.Discount,
				ProductTotalPrice:  totalPrice,
				PriceAfterDiscount: totalPrice - (totalPrice * int(apriori.Discount) / 100),
				Image:              apriori.Image,
			})
		}
	}

	return productResponses, nil
}

func (service *ProductServiceImpl) FindByCode(ctx context.Context, code string) (*model.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return &model.GetProductResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	productResponse, err := service.ProductRepository.FindByCode(ctx, tx, code)
	if err != nil {
		return &model.GetProductResponse{}, err
	}

	return productResponse.ToProductResponse(), nil
}

func (service *ProductServiceImpl) Create(ctx context.Context, request *model.CreateProductRequest) (*model.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return &model.GetProductResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	timeNow, err := time.Parse(helper.TimeFormat, time.Now().Format(helper.TimeFormat))
	if err != nil {
		return &model.GetProductResponse{}, err
	}
	if request.Image == "" {
		request.Image = "no-image.png"
	}

	productRequest := entity.Product{
		Code:        request.Code,
		Name:        helper.UpperWords(request.Name),
		Description: request.Description,
		Price:       request.Price,
		Image:       request.Image,
		Category:    helper.UpperWords(request.Category),
		IsEmpty:     0,
		Mass:        request.Mass,
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	}

	productResponse, err := service.ProductRepository.Create(ctx, tx, &productRequest)
	if err != nil {
		return &model.GetProductResponse{}, err
	}

	return productResponse.ToProductResponse(), nil
}

func (service *ProductServiceImpl) Update(ctx context.Context, request *model.UpdateProductRequest) (*model.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return &model.GetProductResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, request.Code)
	if err != nil {
		return &model.GetProductResponse{}, err
	}

	timeNow, err := time.Parse(helper.TimeFormat, time.Now().Format(helper.TimeFormat))
	if err != nil {
		return &model.GetProductResponse{}, err
	}

	product.Name = helper.UpperWords(request.Name)
	product.Description = request.Description
	product.Price = request.Price
	product.Category = helper.UpperWords(request.Category)
	product.IsEmpty = request.IsEmpty
	product.Mass = request.Mass
	product.UpdatedAt = timeNow
	if request.Image != "" {
		product.Image = request.Image
	}

	productResponse, err := service.ProductRepository.Update(ctx, tx, product)
	if err != nil {
		return &model.GetProductResponse{}, err
	}

	return productResponse.ToProductResponse(), nil
}

func (service *ProductServiceImpl) Delete(ctx context.Context, code string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, code)
	if err != nil {
		return err
	}

	err = service.ProductRepository.Delete(ctx, tx, product.Code)
	if err != nil {
		return err
	}

	return nil
}
