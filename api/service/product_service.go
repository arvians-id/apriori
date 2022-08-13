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
	FindAllOnAdmin(ctx context.Context) ([]model.GetProductResponse, error)
	FindAll(ctx context.Context, search string, category string) ([]model.GetProductResponse, error)
	FindAllSimilarCategory(ctx context.Context, code string) ([]model.GetProductResponse, error)
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
	DB   *sql.DB
	date string
}

func NewProductService(productRepository *repository.ProductRepository, storageService StorageService, aprioriRepository *repository.AprioriRepository, db *sql.DB) ProductService {
	return &productService{
		ProductRepository: *productRepository,
		AprioriRepository: *aprioriRepository,
		StorageService:    storageService,
		DB:                db,
		date:              "2006-01-02 15:04:05",
	}
}

func (service *productService) FindAllOnAdmin(ctx context.Context) ([]model.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []model.GetProductResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	products, err := service.ProductRepository.FindAllOnAdmin(ctx, tx)
	if err != nil {
		return []model.GetProductResponse{}, err
	}

	var productResponse []model.GetProductResponse
	for _, product := range products {
		productResponse = append(productResponse, utils.ToProductResponse(product))
	}

	return productResponse, nil
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

	var productResponse []model.GetProductResponse
	for _, product := range products {
		productResponse = append(productResponse, utils.ToProductResponse(product))
	}

	return productResponse, nil
}

func (service *productService) FindAllSimilarCategory(ctx context.Context, code string) ([]model.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []model.GetProductResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	getProduct, err := service.ProductRepository.FindByCode(ctx, tx, code)
	if err != nil {
		return []model.GetProductResponse{}, err
	}

	categoryArray := strings.Split(getProduct.Category, ", ")
	category := strings.Join(categoryArray, "|")
	products, err := service.ProductRepository.FindAllSimilarCategory(ctx, tx, category)
	if err != nil {
		return []model.GetProductResponse{}, err
	}

	var productResponse []model.GetProductResponse
	for _, product := range products {
		if product.Code != code {
			productResponse = append(productResponse, utils.ToProductResponse(product))
		}
	}

	return productResponse, nil
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

	apriories, err := service.AprioriRepository.FindByActive(ctx, tx)
	if err != nil {
		return nil, err
	}

	var aprioriResponse []model.GetProductRecommendationResponse
	for _, apriori := range apriories {
		items := strings.Split(apriori.Item, ",")
		var exists bool
		for _, nameProduct := range items {
			if strings.ToLower(product.Name) == strings.TrimSpace(nameProduct) {
				exists = true
			}
		}

		var totalPrice int
		if exists {
			for _, nameProduct := range items {
				filterProduct, _ := service.ProductRepository.FindByName(ctx, tx, utils.UpperWords(nameProduct))
				totalPrice += filterProduct.Price
			}

			aprioriResponse = append(aprioriResponse, model.GetProductRecommendationResponse{
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

	return aprioriResponse, nil
}

func (service *productService) FindByCode(ctx context.Context, code string) (model.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetProductResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, code)
	if err != nil {
		return model.GetProductResponse{}, err
	}

	return utils.ToProductResponse(product), nil
}

func (service *productService) Create(ctx context.Context, request model.CreateProductRequest) (model.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetProductResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	createdAt, err := time.Parse(service.date, time.Now().Format(service.date))
	if err != nil {
		return model.GetProductResponse{}, err
	}
	updatedAt, err := time.Parse(service.date, time.Now().Format(service.date))
	if err != nil {
		return model.GetProductResponse{}, err
	}
	if request.Image == "" {
		request.Image = "no-image.png"
	}

	createProduct := entity.Product{
		Code:        request.Code,
		Name:        utils.UpperWords(request.Name),
		Description: request.Description,
		Price:       request.Price,
		Image:       request.Image,
		Category:    utils.UpperWords(request.Category),
		IsEmpty:     0,
		Mass:        request.Mass,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	product, err := service.ProductRepository.Create(ctx, tx, createProduct)
	if err != nil {
		return model.GetProductResponse{}, err
	}

	return utils.ToProductResponse(product), nil
}

func (service *productService) Update(ctx context.Context, request model.UpdateProductRequest) (model.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetProductResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	getProduct, err := service.ProductRepository.FindByCode(ctx, tx, request.Code)
	if err != nil {
		return model.GetProductResponse{}, err
	}

	updatedAt, err := time.Parse(service.date, time.Now().Format(service.date))
	if err != nil {
		return model.GetProductResponse{}, err
	}

	getProduct.Name = utils.UpperWords(request.Name)
	getProduct.Description = request.Description
	getProduct.Price = request.Price
	getProduct.Category = utils.UpperWords(request.Category)
	getProduct.IsEmpty = request.IsEmpty
	getProduct.Mass = request.Mass
	getProduct.UpdatedAt = updatedAt
	if request.Image != "" {
		getProduct.Image = request.Image
	}

	product, err := service.ProductRepository.Update(ctx, tx, getProduct)
	if err != nil {
		return model.GetProductResponse{}, err
	}

	return utils.ToProductResponse(product), nil
}

func (service *productService) Delete(ctx context.Context, code string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	getProduct, err := service.ProductRepository.FindByCode(ctx, tx, code)
	if err != nil {
		return err
	}

	err = service.ProductRepository.Delete(ctx, tx, getProduct.Code)
	if err != nil {
		return err
	}

	return nil
}
