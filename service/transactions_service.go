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
	FindAll(ctx context.Context) ([]model.GetTransactionResponse, error)
	FindByTransaction(ctx context.Context, noTransaction string) (model.GetTransactionResponse, error)
	Create(ctx context.Context, request model.CreateTransactionRequest) error
	CreateFromCsv(ctx context.Context, data [][]string) error
	Update(ctx context.Context, request model.UpdateTransactionRequest) error
	Delete(ctx context.Context, noTransaction string) error
}

type transactionService struct {
	TransactionRepository repository.TransactionRepository
	ProductRepository     repository.ProductRepository
	DB                    *sql.DB
	date                  string
}

func NewTransactionService(transactionRepository *repository.TransactionRepository, productRepository *repository.ProductRepository, db *sql.DB) TransactionService {
	return &transactionService{
		TransactionRepository: *transactionRepository,
		ProductRepository:     *productRepository,
		DB:                    db,
		date:                  "2006-01-02 15:04:05",
	}
}

func (service *transactionService) FindAll(ctx context.Context) ([]model.GetTransactionResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []model.GetTransactionResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	transaction, err := service.TransactionRepository.FindAll(ctx, tx)
	if err != nil {
		return []model.GetTransactionResponse{}, err
	}

	var transactions []model.GetTransactionResponse
	for _, rows := range transaction {
		transactions = append(transactions, utils.ToTransactionResponse(rows))
	}

	return transactions, nil
}

func (service *transactionService) FindByTransaction(ctx context.Context, noTransaction string) (model.GetTransactionResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetTransactionResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	rows, err := service.TransactionRepository.FindByTransaction(ctx, tx, noTransaction)
	if err != nil {
		return model.GetTransactionResponse{}, err
	}

	return utils.ToTransactionResponse(rows), nil
}

func (service *transactionService) Create(ctx context.Context, request model.CreateTransactionRequest) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	timeNow, err := time.Parse(service.date, time.Now().Format(service.date))
	if err != nil {
		return err
	}

	transaction := entity.Transaction{
		ProductName:   request.ProductName,
		CustomerName:  request.CustomerName,
		NoTransaction: request.NoTransaction,
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
	}

	err = service.TransactionRepository.Create(ctx, tx, transaction)
	if err != nil {
		return err
	}

	return nil
}

func (service *transactionService) CreateFromCsv(ctx context.Context, data [][]string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	var transactions []entity.Transaction

	for _, transaction := range data {
		createdAt, _ := time.Parse(service.date, transaction[3]+" 00:00:00")

		transactions = append(transactions, entity.Transaction{
			ProductName:   transaction[0],
			CustomerName:  transaction[1],
			NoTransaction: transaction[2],
			CreatedAt:     createdAt,
			UpdatedAt:     createdAt,
		})
	}

	err = service.TransactionRepository.CreateFromCsv(ctx, tx, transactions)
	if err != nil {
		return err
	}

	return nil
}

func (service *transactionService) Update(ctx context.Context, request model.UpdateTransactionRequest) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	// Find Transaction by number transaction
	transaction, err := service.TransactionRepository.FindByTransaction(ctx, tx, request.NoTransaction)
	if err != nil {
		return err
	}

	updatedAt, err := time.Parse(service.date, time.Now().Format(service.date))
	if err != nil {
		return err
	}

	transaction.ProductName = request.ProductName
	transaction.CustomerName = request.CustomerName
	transaction.NoTransaction = request.NoTransaction
	transaction.UpdatedAt = updatedAt

	entityTransaction := entity.Transaction{
		ProductName:   transaction.ProductName,
		CustomerName:  transaction.CustomerName,
		NoTransaction: transaction.NoTransaction,
		UpdatedAt:     transaction.UpdatedAt,
	}

	err = service.TransactionRepository.Update(ctx, tx, entityTransaction)
	if err != nil {
		return err
	}

	return nil
}

func (service *transactionService) Delete(ctx context.Context, noTransaction string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	rows, err := service.TransactionRepository.FindByTransaction(ctx, tx, noTransaction)
	if err != nil {
		return err
	}

	err = service.TransactionRepository.Delete(ctx, tx, rows.NoTransaction)
	if err != nil {
		return err
	}

	return nil
}
