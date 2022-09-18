package service

import (
	"apriori/entity"
	"apriori/helper"
	"apriori/model"
	"apriori/repository"
	"context"
	"database/sql"
	"time"
)

type TransactionServiceImpl struct {
	TransactionRepository repository.TransactionRepository
	ProductRepository     repository.ProductRepository
	DB                    *sql.DB
}

func NewTransactionService(
	transactionRepository *repository.TransactionRepository,
	productRepository *repository.ProductRepository,
	db *sql.DB,
) TransactionService {
	return &TransactionServiceImpl{
		TransactionRepository: *transactionRepository,
		ProductRepository:     *productRepository,
		DB:                    db,
	}
}

func (service *TransactionServiceImpl) FindAll(ctx context.Context) ([]*entity.Transaction, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	transactions, err := service.TransactionRepository.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (service *TransactionServiceImpl) FindByNoTransaction(ctx context.Context, noTransaction string) (*entity.Transaction, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	transaction, err := service.TransactionRepository.FindByNoTransaction(ctx, tx, noTransaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (service *TransactionServiceImpl) Create(ctx context.Context, request *model.CreateTransactionRequest) (*entity.Transaction, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	timeNow, err := time.Parse(helper.TimeFormat, time.Now().Format(helper.TimeFormat))
	if err != nil {
		return nil, err
	}

	transactionRequest := entity.Transaction{
		ProductName:   request.ProductName,
		CustomerName:  request.CustomerName,
		NoTransaction: helper.CreateTransaction(),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
	}

	transaction, err := service.TransactionRepository.Create(ctx, tx, &transactionRequest)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (service *TransactionServiceImpl) CreateByCsv(ctx context.Context, data [][]string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	var transactions []*entity.Transaction
	for _, transaction := range data {
		createdAt, _ := time.Parse(helper.TimeFormat, transaction[3]+" 00:00:00")
		transactions = append(transactions, &entity.Transaction{
			ProductName:   transaction[0],
			CustomerName:  transaction[1],
			NoTransaction: transaction[2],
			CreatedAt:     createdAt,
			UpdatedAt:     createdAt,
		})
	}

	err = service.TransactionRepository.CreateByCsv(ctx, tx, transactions)
	if err != nil {
		return err
	}

	return nil
}

func (service *TransactionServiceImpl) Update(ctx context.Context, request *model.UpdateTransactionRequest) (*entity.Transaction, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Find Transaction by number transaction
	transaction, err := service.TransactionRepository.FindByNoTransaction(ctx, tx, request.NoTransaction)
	if err != nil {
		return nil, err
	}

	timeNow, err := time.Parse(helper.TimeFormat, time.Now().Format(helper.TimeFormat))
	if err != nil {
		return nil, err
	}

	transaction.ProductName = request.ProductName
	transaction.CustomerName = request.CustomerName
	transaction.NoTransaction = request.NoTransaction
	transaction.UpdatedAt = timeNow

	_, err = service.TransactionRepository.Update(ctx, tx, transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (service *TransactionServiceImpl) Delete(ctx context.Context, noTransaction string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	transaction, err := service.TransactionRepository.FindByNoTransaction(ctx, tx, noTransaction)
	if err != nil {
		return err
	}

	err = service.TransactionRepository.Delete(ctx, tx, transaction.NoTransaction)
	if err != nil {
		return err
	}

	return nil
}

func (service *TransactionServiceImpl) Truncate(ctx context.Context) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	err = service.TransactionRepository.Truncate(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}
