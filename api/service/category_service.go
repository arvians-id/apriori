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
	FindAll(ctx context.Context) ([]*model.GetCategoryResponse, error)
	FindById(ctx context.Context, id int) (*model.GetCategoryResponse, error)
	Create(ctx context.Context, request *model.CreateCategoryRequest) (*model.GetCategoryResponse, error)
	Update(ctx context.Context, request *model.UpdateCategoryRequest) (*model.GetCategoryResponse, error)
	Delete(ctx context.Context, id int) error
}

type categoryService struct {
	categoryRepository repository.CategoryRepository
	DB                 *sql.DB
}

func NewCategoryService(categoryRepository *repository.CategoryRepository, db *sql.DB) CategoryService {
	return &categoryService{
		categoryRepository: *categoryRepository,
		DB:                 db,
	}
}

func (service *categoryService) FindAll(ctx context.Context) ([]*model.GetCategoryResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	categories, err := service.categoryRepository.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	var categoryResponses []*model.GetCategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, utils.ToCategoryResponse(category))
	}

	return categoryResponses, nil
}

func (service *categoryService) FindById(ctx context.Context, id int) (*model.GetCategoryResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return &model.GetCategoryResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	categoryResponse, err := service.categoryRepository.FindById(ctx, tx, id)
	if err != nil {
		return &model.GetCategoryResponse{}, err
	}

	return utils.ToCategoryResponse(categoryResponse), nil
}

func (service *categoryService) Create(ctx context.Context, request *model.CreateCategoryRequest) (*model.GetCategoryResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return &model.GetCategoryResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	timeNow, err := time.Parse(utils.TimeFormat, time.Now().Format(utils.TimeFormat))
	if err != nil {
		return &model.GetCategoryResponse{}, err
	}

	categoryRequest := entity.Category{
		Name:      utils.UpperWords(request.Name),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	categoryResponse, err := service.categoryRepository.Create(ctx, tx, &categoryRequest)
	if err != nil {
		return &model.GetCategoryResponse{}, err
	}

	return utils.ToCategoryResponse(categoryResponse), nil

}

func (service *categoryService) Update(ctx context.Context, request *model.UpdateCategoryRequest) (*model.GetCategoryResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return &model.GetCategoryResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	category, err := service.categoryRepository.FindById(ctx, tx, request.IdCategory)
	if err != nil {
		return &model.GetCategoryResponse{}, err
	}

	timeNow, err := time.Parse(utils.TimeFormat, time.Now().Format(utils.TimeFormat))
	if err != nil {
		return &model.GetCategoryResponse{}, err
	}
	category.Name = utils.UpperWords(request.Name)
	category.UpdatedAt = timeNow

	categoryResponse, err := service.categoryRepository.Update(ctx, tx, category)
	if err != nil {
		return &model.GetCategoryResponse{}, err
	}

	return utils.ToCategoryResponse(categoryResponse), nil
}

func (service *categoryService) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	category, err := service.categoryRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}

	err = service.categoryRepository.Delete(ctx, tx, category.IdCategory)
	if err != nil {
		return err
	}

	return nil
}
