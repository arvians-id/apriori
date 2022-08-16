package service

import (
	"apriori/model"
	"apriori/repository"
	"apriori/utils"
	"context"
	"database/sql"
)

type UserOrderService interface {
	FindAllByPayload(ctx context.Context, payloadId int) ([]model.GetUserOrderResponse, error)
	FindAllByUserId(ctx context.Context, userId int) ([]model.GetUserOrderRelationByUserIdResponse, error)
	FindById(ctx context.Context, orderId int) (model.GetUserOrderResponse, error)
}

type userOrderService struct {
	DB                  *sql.DB
	PaymentRepository   repository.PaymentRepository
	UserOrderRepository repository.UserOrderRepository
	UserRepository      repository.UserRepository
}

func NewUserOrderService(paymentRepository *repository.PaymentRepository, userOrderRepository *repository.UserOrderRepository, userRepository *repository.UserRepository, db *sql.DB) UserOrderService {
	return &userOrderService{
		PaymentRepository:   *paymentRepository,
		UserOrderRepository: *userOrderRepository,
		DB:                  db,
		UserRepository:      *userRepository,
	}
}

func (service *userOrderService) FindAllByPayload(ctx context.Context, payloadId int) ([]model.GetUserOrderResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	userOrders, err := service.UserOrderRepository.FindAll(ctx, tx, utils.IntToStr(payloadId))
	if err != nil {
		return nil, err
	}

	var userOrderResponse []model.GetUserOrderResponse
	for _, userOrder := range userOrders {
		userOrderResponse = append(userOrderResponse, utils.ToUserOrderResponse(userOrder))
	}

	return userOrderResponse, nil
}

func (service *userOrderService) FindAllByUserId(ctx context.Context, userId int) ([]model.GetUserOrderRelationByUserIdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	_, err = service.UserRepository.FindById(ctx, tx, uint64(userId))
	if err != nil {
		return nil, err
	}

	userOrders, err := service.UserOrderRepository.FindAllByUserId(ctx, tx, userId)
	if err != nil {
		return nil, err
	}

	var userOrderResponse []model.GetUserOrderRelationByUserIdResponse
	for _, userOrder := range userOrders {
		userOrderResponse = append(userOrderResponse, utils.ToUserOrderRelationByUserIdResponse(userOrder))
	}

	return userOrderResponse, nil
}

func (service *userOrderService) FindById(ctx context.Context, orderId int) (model.GetUserOrderResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetUserOrderResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	userOrder, err := service.UserOrderRepository.FindById(ctx, tx, orderId)
	if err != nil {
		return model.GetUserOrderResponse{}, err
	}

	return utils.ToUserOrderResponse(userOrder), nil
}
