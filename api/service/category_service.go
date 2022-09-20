package service

import (
	"context"
	"database/sql"
	"github.com/arvians-id/apriori/helper"
	"github.com/arvians-id/apriori/http/controller/rest/request"
	"github.com/arvians-id/apriori/model"
	"github.com/arvians-id/apriori/repository"
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

func (service *CategoryServiceImpl) FindAll(ctx context.Context) ([]*model.Category, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	categories, err := service.CategoryRepository.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (service *CategoryServiceImpl) FindById(ctx context.Context, id int) (*model.Category, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	category, err := service.CategoryRepository.FindById(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (service *CategoryServiceImpl) Create(ctx context.Context, request *request.CreateCategoryRequest) (*model.Category, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	timeNow, err := time.Parse(helper.TimeFormat, time.Now().Format(helper.TimeFormat))
	if err != nil {
		return nil, err
	}

	categoryRequest := model.Category{
		Name:      helper.UpperWords(request.Name),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	category, err := service.CategoryRepository.Create(ctx, tx, &categoryRequest)
	if err != nil {
		return nil, err
	}

	return category, nil

}

func (service *CategoryServiceImpl) Update(ctx context.Context, request *request.UpdateCategoryRequest) (*model.Category, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	category, err := service.CategoryRepository.FindById(ctx, tx, request.IdCategory)
	if err != nil {
		return nil, err
	}

	timeNow, err := time.Parse(helper.TimeFormat, time.Now().Format(helper.TimeFormat))
	if err != nil {
		return nil, err
	}
	category.Name = helper.UpperWords(request.Name)
	category.UpdatedAt = timeNow

	_, err = service.CategoryRepository.Update(ctx, tx, category)
	if err != nil {
		return nil, err
	}

	return category, nil
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
