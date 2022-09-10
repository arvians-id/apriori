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

func (service *CommentServiceImpl) FindAllRatingByProductCode(ctx context.Context, productCode string) ([]*model.GetRatingResponse, error) {
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

	var ratingResponses []*model.GetRatingResponse
	for _, comment := range ratings {
		ratingResponses = append(ratingResponses, comment.ToRatingResponse())
	}

	return ratingResponses, nil
}

func (service *CommentServiceImpl) FindAllByProductCode(ctx context.Context, productCode string, rating string, tags string) ([]*model.GetCommentResponse, error) {
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

	var commentResponses []*model.GetCommentResponse
	for _, comment := range comments {
		commentResponses = append(commentResponses, comment.ToCommentResponse())
	}

	return commentResponses, nil
}

func (service *CommentServiceImpl) FindById(ctx context.Context, id int) (*model.GetCommentResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return &model.GetCommentResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	commentResponse, err := service.CommentRepository.FindById(ctx, tx, id)
	if err != nil {
		return &model.GetCommentResponse{}, err
	}

	return commentResponse.ToCommentResponse(), nil
}

func (service *CommentServiceImpl) FindByUserOrderId(ctx context.Context, userOrderId int) (*model.GetCommentResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return &model.GetCommentResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	commentResponse, err := service.CommentRepository.FindByUserOrderId(ctx, tx, userOrderId)
	if err != nil {
		return &model.GetCommentResponse{}, err
	}

	return commentResponse.ToCommentResponse(), nil
}

func (service *CommentServiceImpl) Create(ctx context.Context, request *model.CreateCommentRequest) (*model.GetCommentResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return &model.GetCommentResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	timeNow, err := time.Parse(helper.TimeFormat, time.Now().Format(helper.TimeFormat))
	if err != nil {
		return &model.GetCommentResponse{}, err
	}

	commentRequest := entity.Comment{
		UserOrderId: request.UserOrderId,
		ProductCode: request.ProductCode,
		Description: sql.NullString{
			String: request.Description,
			Valid:  true,
		},
		Rating: request.Rating,
		Tag: sql.NullString{
			String: request.Tag,
			Valid:  true,
		},
		CreatedAt: timeNow,
	}

	commentResponse, err := service.CommentRepository.Create(ctx, tx, &commentRequest)
	if err != nil {
		return &model.GetCommentResponse{}, err
	}

	return commentResponse.ToCommentResponse(), nil
}
