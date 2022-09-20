package service

import (
	"context"
	"database/sql"
	"github.com/arvians-id/apriori/helper"
	"github.com/arvians-id/apriori/http/request"
	"github.com/arvians-id/apriori/model"
	"github.com/arvians-id/apriori/repository"
	"strings"
	"time"
)

type CommentServiceImpl struct {
	CommentRepository repository.CommentRepository
	ProductRepository repository.ProductRepository
	DB                *sql.DB
}

func NewCommentService(
	commentRepository *repository.CommentRepository,
	productRepository *repository.ProductRepository,
	db *sql.DB,
) CommentService {
	return &CommentServiceImpl{
		CommentRepository: *commentRepository,
		ProductRepository: *productRepository,
		DB:                db,
	}
}

func (service *CommentServiceImpl) FindAllRatingByProductCode(ctx context.Context, productCode string) ([]*model.RatingFromComment, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, productCode)
	if err != nil {
		return nil, err
	}

	ratings, err := service.CommentRepository.FindAllRatingByProductCode(ctx, tx, product.Code)
	if err != nil {
		return nil, err
	}

	return ratings, nil
}

func (service *CommentServiceImpl) FindAllByProductCode(ctx context.Context, productCode string, rating string, tags string) ([]*model.Comment, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, productCode)
	if err != nil {
		return nil, err
	}

	tagArray := strings.Split(tags, ",")
	tag := strings.Join(tagArray, "|")
	comments, err := service.CommentRepository.FindAllByProductCode(ctx, tx, product.Code, rating, tag)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (service *CommentServiceImpl) FindById(ctx context.Context, id int) (*model.Comment, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	comment, err := service.CommentRepository.FindById(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (service *CommentServiceImpl) FindByUserOrderId(ctx context.Context, userOrderId int) (*model.Comment, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	comment, err := service.CommentRepository.FindByUserOrderId(ctx, tx, userOrderId)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (service *CommentServiceImpl) Create(ctx context.Context, request *request.CreateCommentRequest) (*model.Comment, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	timeNow, err := time.Parse(helper.TimeFormat, time.Now().Format(helper.TimeFormat))
	if err != nil {
		return nil, err
	}

	commentRequest := model.Comment{
		UserOrderId: request.UserOrderId,
		ProductCode: request.ProductCode,
		Description: &request.Description,
		Rating:      request.Rating,
		Tag:         &request.Tag,
		CreatedAt:   timeNow,
	}

	comment, err := service.CommentRepository.Create(ctx, tx, &commentRequest)
	if err != nil {
		return nil, err
	}

	return comment, nil
}
