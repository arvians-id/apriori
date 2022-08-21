package service

import (
	"apriori/model"
	"apriori/repository"
	"apriori/utils"
	"context"
	"database/sql"
)

type UserOrderService interface {
	FindAllByPayloadId(ctx context.Context, payloadId int) ([]*model.GetUserOrderResponse, error)
	FindAllByUserId(ctx context.Context, userId int) ([]*model.GetUserOrderRelationByUserIdResponse, error)
	FindById(ctx context.Context, id int) (*model.GetUserOrderResponse, error)
}

type userOrderService struct {
	PaymentRepository   repository.PaymentRepository
	UserOrderRepository repository.UserOrderRepository
	UserRepository      repository.UserRepository
	DB                  *sql.DB
}

func NewUserOrderService(paymentRepository *repository.PaymentRepository, userOrderRepository *repository.UserOrderRepository, userRepository *repository.UserRepository, db *sql.DB) UserOrderService {
	return &userOrderService{
		PaymentRepository:   *paymentRepository,
		UserOrderRepository: *userOrderRepository,
		UserRepository:      *userRepository,
		DB:                  db,
	}
}

func (service *userOrderService) FindAllByPayloadId(ctx context.Context, payloadId int) ([]*model.GetUserOrderResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	userOrders, err := service.UserOrderRepository.FindAllByPayloadId(ctx, tx, utils.IntToStr(payloadId))
	if err != nil {
		return nil, err
	}

	var userOrderResponses []*model.GetUserOrderResponse
	for _, userOrder := range userOrders {
		userOrderResponses = append(userOrderResponses, utils.ToUserOrderResponse(userOrder))
	}

	return userOrderResponses, nil
}

func (service *userOrderService) FindAllByUserId(ctx context.Context, userId int) ([]*model.GetUserOrderRelationByUserIdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	_, err = service.UserRepository.FindById(ctx, tx, userId)
	if err != nil {
		return nil, err
	}

	userOrders, err := service.UserOrderRepository.FindAllByUserId(ctx, tx, userId)
	if err != nil {
		return nil, err
	}

	var userOrderResponses []*model.GetUserOrderRelationByUserIdResponse
	for _, userOrder := range userOrders {
		userOrderResponses = append(userOrderResponses, utils.ToUserOrderRelationByUserIdResponse(userOrder))
	}

	return userOrderResponses, nil
}

func (service *userOrderService) FindById(ctx context.Context, id int) (*model.GetUserOrderResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return &model.GetUserOrderResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	userOrderResponse, err := service.UserOrderRepository.FindById(ctx, tx, id)
	if err != nil {
		return &model.GetUserOrderResponse{}, err
	}

	return utils.ToUserOrderResponse(userOrderResponse), nil
}
