package service

import (
	"context"
	"database/sql"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/arvians-id/apriori/internal/repository"
	"github.com/arvians-id/apriori/util"
	"log"
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
		log.Println("[UserOrderService][FindAllByPayloadId] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	userOrders, err := service.UserOrderRepository.FindAllByPayloadId(ctx, tx, util.IntToStr(payloadId))
	if err != nil {
		log.Println("[UserOrderService][FindAllByPayloadId][FindAllByPayloadId] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return userOrders, nil
}

func (service *UserOrderServiceImpl) FindAllByUserId(ctx context.Context, userId int) ([]*model.UserOrder, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[UserOrderService][FindAllByUserId] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	_, err = service.UserRepository.FindById(ctx, tx, userId)
	if err != nil {
		log.Println("[UserOrderService][FindAllByUserId][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	userOrders, err := service.UserOrderRepository.FindAllByUserId(ctx, tx, userId)
	if err != nil {
		log.Println("[UserOrderService][FindAllByUserId][FindAllByUserId] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return userOrders, nil
}

func (service *UserOrderServiceImpl) FindById(ctx context.Context, id int) (*model.UserOrder, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[UserOrderService][FindById] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	userOrder, err := service.UserOrderRepository.FindById(ctx, tx, id)
	if err != nil {
		log.Println("[UserOrderService][FindById][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return userOrder, nil
}
