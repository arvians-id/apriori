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
}

type userOrderService struct {
	DB                  *sql.DB
	PaymentRepository   repository.PaymentRepository
	UserOrderRepository repository.UserOrderRepository
}

func NewUserOrderService(paymentRepository *repository.PaymentRepository, userOrderRepository *repository.UserOrderRepository, db *sql.DB) UserOrderService {
	return &userOrderService{
		PaymentRepository:   *paymentRepository,
		UserOrderRepository: *userOrderRepository,
		DB:                  db,
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
