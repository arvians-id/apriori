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

type CommentService interface {
	FindAllRatingByProductCode(ctx context.Context, productCode string) ([]model.GetRatingResponse, error)
	FindAllByProductCode(ctx context.Context, productCode string, rating string, tags string) ([]model.GetCommentResponse, error)
	FindById(ctx context.Context, id int) (model.GetCommentResponse, error)
	FindByUserOrderId(ctx context.Context, userOrderId int) (model.GetCommentResponse, error)
	Create(ctx context.Context, request model.CreateCommentRequest) (model.GetCommentResponse, error)
}

type commentService struct {
	CommentRepository repository.CommentRepository
	ProductRepository repository.ProductRepository
	db                *sql.DB
}

func NewCommentService(commentRepository *repository.CommentRepository, productRepository *repository.ProductRepository, db *sql.DB) CommentService {
	return &commentService{
		CommentRepository: *commentRepository,
		ProductRepository: *productRepository,
		db:                db,
	}
}

func (service *commentService) FindAllRatingByProductCode(ctx context.Context, productCode string) ([]model.GetRatingResponse, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, productCode)
	if err != nil {
		return nil, err
	}

	ratings, err := service.CommentRepository.FindAllRatingByProductCode(ctx, tx, product.Code)
	if err != nil {
		return nil, err
	}

	var ratingResponses []model.GetRatingResponse
	for _, comment := range ratings {
		ratingResponses = append(ratingResponses, utils.ToRatingResponse(comment))
	}

	return ratingResponses, nil
}

func (service *commentService) FindAllByProductCode(ctx context.Context, productCode string, rating string, tags string) ([]model.GetCommentResponse, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

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

	var commentResponses []model.GetCommentResponse
	for _, comment := range comments {
		commentResponses = append(commentResponses, utils.ToCommentResponse(comment))
	}

	return commentResponses, nil
}

func (service *commentService) FindById(ctx context.Context, id int) (model.GetCommentResponse, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return model.GetCommentResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	commentResponse, err := service.CommentRepository.FindById(ctx, tx, id)
	if err != nil {
		return model.GetCommentResponse{}, err
	}

	return utils.ToCommentResponse(commentResponse), nil
}

func (service *commentService) FindByUserOrderId(ctx context.Context, userOrderId int) (model.GetCommentResponse, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return model.GetCommentResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	commentResponse, err := service.CommentRepository.FindByUserOrderId(ctx, tx, userOrderId)
	if err != nil {
		return model.GetCommentResponse{}, err
	}

	return utils.ToCommentResponse(commentResponse), nil
}

func (service *commentService) Create(ctx context.Context, request model.CreateCommentRequest) (model.GetCommentResponse, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return model.GetCommentResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	timeNow, err := time.Parse(utils.TimeFormat, time.Now().Format(utils.TimeFormat))
	if err != nil {
		return model.GetCommentResponse{}, err
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

	commentResponse, err := service.CommentRepository.Create(ctx, tx, commentRequest)
	if err != nil {
		return model.GetCommentResponse{}, err
	}

	return utils.ToCommentResponse(commentResponse), nil
}
