package service

import (
	"context"
	"database/sql"
	"github.com/arvians-id/apriori/helper"
	"github.com/arvians-id/apriori/model"
	"github.com/arvians-id/apriori/repository"
)

type UserOrderServiceImpl struct {
	PaymentRepository   repository.PaymentRepository
	UserOrderRepository repository.UserOrderRepository
	UserRepository      repository.UserRepository
	DB                  *sql.DB
}

func NewUserOrderService(
	paymentRepository *repository.PaymentRepository,
	userOrderRepository *repository.UserOrderRepository,
	userRepository *repository.UserRepository,
	db *sql.DB,
) UserOrderService {
	return &UserOrderServiceImpl{
		PaymentRepository:   *paymentRepository,
		UserOrderRepository: *userOrderRepository,
		UserRepository:      *userRepository,
		DB:                  db,
	}
}

func (service *UserOrderServiceImpl) FindAllByPayloadId(ctx context.Context, payloadId int) ([]*model.UserOrder, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	userOrders, err := service.UserOrderRepository.FindAllByPayloadId(ctx, tx, helper.IntToStr(payloadId))
	if err != nil {
		return nil, err
	}

	return userOrders, nil
}

func (service *UserOrderServiceImpl) FindAllByUserId(ctx context.Context, userId int) ([]*model.UserOrder, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	_, err = service.UserRepository.FindById(ctx, tx, userId)
	if err != nil {
		return nil, err
	}

	userOrders, err := service.UserOrderRepository.FindAllByUserId(ctx, tx, userId)
	if err != nil {
		return nil, err
	}

	return userOrders, nil
}

func (service *UserOrderServiceImpl) FindById(ctx context.Context, id int) (*model.UserOrder, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	userOrder, err := service.UserOrderRepository.FindById(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return userOrder, nil
}
