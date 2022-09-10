package service

import (
	"apriori/entity"
	"apriori/helper"
	"apriori/model"
	"apriori/repository"
	"context"
	"database/sql"
	"time"
)

type CategoryServiceImpl struct {
	CategoryRepository repository.CategoryRepository
	DB                 *sql.DB
}

func NewCategoryService(categoryRepository *repository.CategoryRepository, db *sql.DB) CategoryService {
	return &CategoryServiceImpl{
		CategoryRepository: *categoryRepository,
		DB:                 db,
	}
}

func (service *CategoryServiceImpl) FindAll(ctx context.Context) ([]*model.GetCategoryResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	categories, err := service.CategoryRepository.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	var categoryResponses []*model.GetCategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, category.ToCategoryResponse())
	}

	return categoryResponses, nil
}

func (service *CategoryServiceImpl) FindById(ctx context.Context, id int) (*model.GetCategoryResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return &model.GetCategoryResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	categoryResponse, err := service.CategoryRepository.FindById(ctx, tx, id)
	if err != nil {
		return &model.GetCategoryResponse{}, err
	}

	return categoryResponse.ToCategoryResponse(), nil
}

func (service *CategoryServiceImpl) Create(ctx context.Context, request *model.CreateCategoryRequest) (*model.GetCategoryResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return &model.GetCategoryResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	timeNow, err := time.Parse(helper.TimeFormat, time.Now().Format(helper.TimeFormat))
	if err != nil {
		return &model.GetCategoryResponse{}, err
	}

	categoryRequest := entity.Category{
		Name:      helper.UpperWords(request.Name),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	categoryResponse, err := service.CategoryRepository.Create(ctx, tx, &categoryRequest)
	if err != nil {
		return &model.GetCategoryResponse{}, err
	}

	return categoryResponse.ToCategoryResponse(), nil

}

func (service *CategoryServiceImpl) Update(ctx context.Context, request *model.UpdateCategoryRequest) (*model.GetCategoryResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return &model.GetCategoryResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	category, err := service.CategoryRepository.FindById(ctx, tx, request.IdCategory)
	if err != nil {
		return &model.GetCategoryResponse{}, err
	}

	timeNow, err := time.Parse(helper.TimeFormat, time.Now().Format(helper.TimeFormat))
	if err != nil {
		return &model.GetCategoryResponse{}, err
	}
	category.Name = helper.UpperWords(request.Name)
	category.UpdatedAt = timeNow

	categoryResponse, err := service.CategoryRepository.Update(ctx, tx, category)
	if err != nil {
		return &model.GetCategoryResponse{}, err
	}

	return categoryResponse.ToCategoryResponse(), nil
}

func (service *CategoryServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	category, err := service.CategoryRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}

	err = service.CategoryRepository.Delete(ctx, tx, category.IdCategory)
	if err != nil {
		return err
	}

	return nil
}
