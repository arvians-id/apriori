package service

import (
	"context"
	"database/sql"
	"github.com/arvians-id/apriori/internal/http/presenter/request"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/arvians-id/apriori/internal/repository"
	"github.com/arvians-id/apriori/util"
	"log"
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
		log.Println("[CategoryService][FindAll] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	categories, err := service.CategoryRepository.FindAll(ctx, tx)
	if err != nil {
		log.Println("[CategoryService][FindAll][FindAll] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return categories, nil
}

func (service *CategoryServiceImpl) FindById(ctx context.Context, id int) (*model.Category, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[CategoryService][FindById] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	category, err := service.CategoryRepository.FindById(ctx, tx, id)
	if err != nil {
		log.Println("[CategoryService][FindById][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return category, nil
}

func (service *CategoryServiceImpl) Create(ctx context.Context, request *request.CreateCategoryRequest) (*model.Category, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[CategoryService][Create] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[CategoryService][Create] problem in parsing to time, err: ", err.Error())
		return nil, err
	}

	categoryRequest := model.Category{
		Name:      util.UpperWords(request.Name),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	category, err := service.CategoryRepository.Create(ctx, tx, &categoryRequest)
	if err != nil {
		log.Println("[CategoryService][Create][Create] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return category, nil

}

func (service *CategoryServiceImpl) Update(ctx context.Context, request *request.UpdateCategoryRequest) (*model.Category, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[CategoryService][Update] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	category, err := service.CategoryRepository.FindById(ctx, tx, request.IdCategory)
	if err != nil {
		log.Println("[CategoryService][Update][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[CategoryService][Update] problem in parsing to time, err: ", err.Error())
		return nil, err
	}
	category.Name = util.UpperWords(request.Name)
	category.UpdatedAt = timeNow

	_, err = service.CategoryRepository.Update(ctx, tx, category)
	if err != nil {
		log.Println("[CategoryService][Update][Update] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return category, nil
}

func (service *CategoryServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[CategoryService][Delete] problem in db transaction, err: ", err.Error())
		return err
	}
	defer util.CommitOrRollback(tx)

	category, err := service.CategoryRepository.FindById(ctx, tx, id)
	if err != nil {
		log.Println("[CategoryService][Delete][FindById] problem in getting from repository, err: ", err.Error())
		return err
	}

	err = service.CategoryRepository.Delete(ctx, tx, category.IdCategory)
	if err != nil {
		log.Println("[CategoryService][Delete][Delete] problem in getting from repository, err: ", err.Error())
		return err
	}

	return nil
}
