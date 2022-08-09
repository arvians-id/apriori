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

type CategoryService interface {
	FindAll(ctx context.Context) ([]model.GetCategoryResponse, error)
	FindById(ctx context.Context, categoryId int) (model.GetCategoryResponse, error)
	Create(ctx context.Context, request model.CreateCategoryRequest) (model.GetCategoryResponse, error)
	Update(ctx context.Context, request model.UpdateCategoryRequest) (model.GetCategoryResponse, error)
	Delete(ctx context.Context, categoryId int) error
}

type categoryService struct {
	categoryRepository repository.CategoryRepository
	DB                 *sql.DB
	date               string
}

func NewCategoryService(categoryRepository *repository.CategoryRepository, db *sql.DB) CategoryService {
	return &categoryService{
		categoryRepository: *categoryRepository,
		DB:                 db,
		date:               "2006-01-02 15:04:05",
	}
}

func (service *categoryService) FindAll(ctx context.Context) ([]model.GetCategoryResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	categories, err := service.categoryRepository.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	var categoryResponse []model.GetCategoryResponse
	for _, category := range categories {
		categoryResponse = append(categoryResponse, utils.ToCategoryResponse(category))
	}

	return categoryResponse, nil
}

func (service *categoryService) FindById(ctx context.Context, categoryId int) (model.GetCategoryResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetCategoryResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	category, err := service.categoryRepository.FindById(ctx, tx, categoryId)
	if err != nil {
		return model.GetCategoryResponse{}, err
	}

	return utils.ToCategoryResponse(category), nil
}

func (service *categoryService) Create(ctx context.Context, request model.CreateCategoryRequest) (model.GetCategoryResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetCategoryResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	createdAt, err := time.Parse(service.date, time.Now().Format(service.date))
	if err != nil {
		return model.GetCategoryResponse{}, err
	}
	updatedAt, err := time.Parse(service.date, time.Now().Format(service.date))
	if err != nil {
		return model.GetCategoryResponse{}, err
	}

	createCategory := entity.Category{
		Name:      utils.UpperWords(request.Name),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	category, err := service.categoryRepository.Create(ctx, tx, createCategory)
	if err != nil {
		return model.GetCategoryResponse{}, err
	}

	return utils.ToCategoryResponse(category), nil

}

func (service *categoryService) Update(ctx context.Context, request model.UpdateCategoryRequest) (model.GetCategoryResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetCategoryResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	category, err := service.categoryRepository.FindById(ctx, tx, request.IdCategory)
	if err != nil {
		return model.GetCategoryResponse{}, err
	}

	updatedAt, err := time.Parse(service.date, time.Now().Format(service.date))
	if err != nil {
		return model.GetCategoryResponse{}, err
	}
	category.Name = utils.UpperWords(request.Name)
	category.UpdatedAt = updatedAt

	categoryUpdate, err := service.categoryRepository.Update(ctx, tx, category)
	if err != nil {
		return model.GetCategoryResponse{}, err
	}

	return utils.ToCategoryResponse(categoryUpdate), nil
}

func (service *categoryService) Delete(ctx context.Context, categoryId int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	_, err = service.categoryRepository.FindById(ctx, tx, categoryId)
	if err != nil {
		return err
	}

	err = service.categoryRepository.Delete(ctx, tx, categoryId)
	if err != nil {
		return err
	}

	return nil
}
