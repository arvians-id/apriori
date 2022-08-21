package service

import (
	"apriori/entity"
	"apriori/model"
	"apriori/repository"
	"apriori/utils"
	"context"
	"database/sql"
	"time"
)

type TransactionService interface {
	FindAll(ctx context.Context) ([]*model.GetTransactionResponse, error)
	FindByNoTransaction(ctx context.Context, noTransaction string) (*model.GetTransactionResponse, error)
	Create(ctx context.Context, request *model.CreateTransactionRequest) (*model.GetTransactionResponse, error)
	CreateByCsv(ctx context.Context, data [][]string) error
	Update(ctx context.Context, request *model.UpdateTransactionRequest) (*model.GetTransactionResponse, error)
	Delete(ctx context.Context, noTransaction string) error
	Truncate(ctx context.Context) error
}

type transactionService struct {
	TransactionRepository repository.TransactionRepository
	ProductRepository     repository.ProductRepository
	DB                    *sql.DB
}

func NewTransactionService(transactionRepository *repository.TransactionRepository, productRepository *repository.ProductRepository, db *sql.DB) TransactionService {
	return &transactionService{
		TransactionRepository: *transactionRepository,
		ProductRepository:     *productRepository,
		DB:                    db,
	}
}

func (service *transactionService) FindAll(ctx context.Context) ([]*model.GetTransactionResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	transactions, err := service.TransactionRepository.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	var transactionResponses []*model.GetTransactionResponse
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, utils.ToTransactionResponse(transaction))
	}

	return transactionResponses, nil
}

func (service *transactionService) FindByNoTransaction(ctx context.Context, noTransaction string) (*model.GetTransactionResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return &model.GetTransactionResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	transactionResponse, err := service.TransactionRepository.FindByNoTransaction(ctx, tx, noTransaction)
	if err != nil {
		return &model.GetTransactionResponse{}, err
	}

	return utils.ToTransactionResponse(transactionResponse), nil
}

func (service *transactionService) Create(ctx context.Context, request *model.CreateTransactionRequest) (*model.GetTransactionResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return &model.GetTransactionResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	timeNow, err := time.Parse(utils.TimeFormat, time.Now().Format(utils.TimeFormat))
	if err != nil {
		return &model.GetTransactionResponse{}, err
	}

	transactionRequest := entity.Transaction{
		ProductName:   request.ProductName,
		CustomerName:  request.CustomerName,
		NoTransaction: utils.CreateTransaction(),
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
	}

	transactionResponse, err := service.TransactionRepository.Create(ctx, tx, &transactionRequest)
	if err != nil {
		return &model.GetTransactionResponse{}, err
	}

	return utils.ToTransactionResponse(transactionResponse), nil
}

func (service *transactionService) CreateByCsv(ctx context.Context, data [][]string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	var transactions []*entity.Transaction
	for _, transaction := range data {
		createdAt, _ := time.Parse(utils.TimeFormat, transaction[3]+" 00:00:00")
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

func (service *transactionService) Update(ctx context.Context, request *model.UpdateTransactionRequest) (*model.GetTransactionResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return &model.GetTransactionResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	// Find Transaction by number transaction
	transaction, err := service.TransactionRepository.FindByNoTransaction(ctx, tx, request.NoTransaction)
	if err != nil {
		return &model.GetTransactionResponse{}, err
	}

	timeNow, err := time.Parse(utils.TimeFormat, time.Now().Format(utils.TimeFormat))
	if err != nil {
		return &model.GetTransactionResponse{}, err
	}

	transactionRequest := entity.Transaction{
		IdTransaction: transaction.IdTransaction,
		ProductName:   request.ProductName,
		CustomerName:  request.CustomerName,
		NoTransaction: request.NoTransaction,
		UpdatedAt:     timeNow,
	}

	transactionResponse, err := service.TransactionRepository.Update(ctx, tx, &transactionRequest)
	if err != nil {
		return &model.GetTransactionResponse{}, err
	}

	return utils.ToTransactionResponse(transactionResponse), nil
}

func (service *transactionService) Delete(ctx context.Context, noTransaction string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

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

func (service *transactionService) Truncate(ctx context.Context) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	err = service.TransactionRepository.Truncate(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}
