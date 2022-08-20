package service

import (
	"apriori/entity"
	"apriori/model"
	"apriori/repository"
	"apriori/utils"
	"context"
	"database/sql"
	"strings"
	"time"
)

type ProductService interface {
	FindAllByAdmin(ctx context.Context) ([]model.GetProductResponse, error)
	FindAll(ctx context.Context, search string, category string) ([]model.GetProductResponse, error)
	FindAllBySimilarCategory(ctx context.Context, code string) ([]model.GetProductResponse, error)
	FindAllRecommendation(ctx context.Context, code string) ([]model.GetProductRecommendationResponse, error)
	FindByCode(ctx context.Context, code string) (model.GetProductResponse, error)
	Create(ctx context.Context, request model.CreateProductRequest) (model.GetProductResponse, error)
	Update(ctx context.Context, request model.UpdateProductRequest) (model.GetProductResponse, error)
	Delete(ctx context.Context, code string) error
}

type productService struct {
	ProductRepository repository.ProductRepository
	AprioriRepository repository.AprioriRepository
	StorageService
	DB *sql.DB
}

func NewProductService(productRepository *repository.ProductRepository, storageService StorageService, aprioriRepository *repository.AprioriRepository, db *sql.DB) ProductService {
	return &productService{
		ProductRepository: *productRepository,
		AprioriRepository: *aprioriRepository,
		StorageService:    storageService,
		DB:                db,
	}
}

func (service *productService) FindAllByAdmin(ctx context.Context) ([]model.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []model.GetProductResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	products, err := service.ProductRepository.FindAllByAdmin(ctx, tx)
	if err != nil {
		return []model.GetProductResponse{}, err
	}

	var productResponses []model.GetProductResponse
	for _, product := range products {
		productResponses = append(productResponses, utils.ToProductResponse(product))
	}

	return productResponses, nil
}

func (service *productService) FindAll(ctx context.Context, search string, category string) ([]model.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []model.GetProductResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	products, err := service.ProductRepository.FindAll(ctx, tx, search, category)
	if err != nil {
		return []model.GetProductResponse{}, err
	}

	var productResponses []model.GetProductResponse
	for _, product := range products {
		productResponses = append(productResponses, utils.ToProductResponse(product))
	}

	return productResponses, nil
}

func (service *productService) FindAllBySimilarCategory(ctx context.Context, code string) ([]model.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []model.GetProductResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, code)
	if err != nil {
		return []model.GetProductResponse{}, err
	}

	categoryArray := strings.Split(product.Category, ", ")
	categoryString := strings.Join(categoryArray, "|")
	productCategories, err := service.ProductRepository.FindAllBySimilarCategory(ctx, tx, categoryString)
	if err != nil {
		return []model.GetProductResponse{}, err
	}

	var productResponses []model.GetProductResponse
	for _, productCategory := range productCategories {
		if productCategory.Code != code {
			productResponses = append(productResponses, utils.ToProductResponse(productCategory))
		}
	}

	return productResponses, nil
}

func (service *productService) FindAllRecommendation(ctx context.Context, code string) ([]model.GetProductRecommendationResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, code)
	if err != nil {
		return nil, err
	}

	apriories, err := service.AprioriRepository.FindAllByActive(ctx, tx)
	if err != nil {
		return nil, err
	}

	var productResponses []model.GetProductRecommendationResponse
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
				product, _ := service.ProductRepository.FindByName(ctx, tx, utils.UpperWords(productName))
				totalPrice += product.Price
			}

			productResponses = append(productResponses, model.GetProductRecommendationResponse{
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

func (service *productService) FindByCode(ctx context.Context, code string) (model.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetProductResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	productResponse, err := service.ProductRepository.FindByCode(ctx, tx, code)
	if err != nil {
		return model.GetProductResponse{}, err
	}

	return utils.ToProductResponse(productResponse), nil
}

func (service *productService) Create(ctx context.Context, request model.CreateProductRequest) (model.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetProductResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	timeNow, err := time.Parse(utils.TimeFormat, time.Now().Format(utils.TimeFormat))
	if err != nil {
		return model.GetProductResponse{}, err
	}
	if request.Image == "" {
		request.Image = "no-image.png"
	}

	productRequest := entity.Product{
		Code:        request.Code,
		Name:        utils.UpperWords(request.Name),
		Description: request.Description,
		Price:       request.Price,
		Image:       request.Image,
		Category:    utils.UpperWords(request.Category),
		IsEmpty:     0,
		Mass:        request.Mass,
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	}

	productResponse, err := service.ProductRepository.Create(ctx, tx, productRequest)
	if err != nil {
		return model.GetProductResponse{}, err
	}

	return utils.ToProductResponse(productResponse), nil
}

func (service *productService) Update(ctx context.Context, request model.UpdateProductRequest) (model.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetProductResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, request.Code)
	if err != nil {
		return model.GetProductResponse{}, err
	}

	timeNow, err := time.Parse(utils.TimeFormat, time.Now().Format(utils.TimeFormat))
	if err != nil {
		return model.GetProductResponse{}, err
	}

	product.Name = utils.UpperWords(request.Name)
	product.Description = request.Description
	product.Price = request.Price
	product.Category = utils.UpperWords(request.Category)
	product.IsEmpty = request.IsEmpty
	product.Mass = request.Mass
	product.UpdatedAt = timeNow
	if request.Image != "" {
		product.Image = request.Image
	}

	productResponse, err := service.ProductRepository.Update(ctx, tx, product)
	if err != nil {
		return model.GetProductResponse{}, err
	}

	return utils.ToProductResponse(productResponse), nil
}

func (service *productService) Delete(ctx context.Context, code string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

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
