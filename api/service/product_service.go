package service

import (
	"apriori/entity"
	"apriori/model"
	"apriori/repository"
	"apriori/utils"
	"context"
	"database/sql"
	"time"
)

type ProductService interface {
	FindAll(ctx context.Context) ([]model.GetProductResponse, error)
	FindByCode(ctx context.Context, code string) (model.GetProductResponse, error)
	Create(ctx context.Context, request model.CreateProductRequest) (model.GetProductResponse, error)
	Update(ctx context.Context, request model.UpdateProductRequest) (model.GetProductResponse, error)
	Delete(ctx context.Context, code string) error
}

type productService struct {
	ProductRepository repository.ProductRepository
	DB                *sql.DB
	date              string
}

func NewProductService(ProductRepository *repository.ProductRepository, db *sql.DB) ProductService {
	return &productService{
		ProductRepository: *ProductRepository,
		DB:                db,
		date:              "2006-01-02 15:04:05",
	}
}

func (service *productService) FindAll(ctx context.Context) ([]model.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []model.GetProductResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	products, err := service.ProductRepository.FindAll(ctx, tx)
	if err != nil {
		return []model.GetProductResponse{}, err
	}

	var productResponse []model.GetProductResponse
	for _, product := range products {
		productResponse = append(productResponse, utils.ToProductResponse(product))
	}

	return productResponse, nil
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

	createProduct := entity.Product{
		Code:        request.Code,
		Name:        request.Name,
		Description: request.Description,
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

	getProduct.Name = request.Name
	getProduct.Description = request.Description
	getProduct.UpdatedAt = updatedAt

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
