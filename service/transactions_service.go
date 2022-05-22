package service

import (
	"apriori/entity"
	"apriori/helper"
	"apriori/model"
	"apriori/repository"
	"context"
	"database/sql"
	"strconv"
	"time"
)

type TransactionService interface {
	FindAll(ctx context.Context) ([]model.GetTransactionProductResponse, error)
	FindByTransaction(ctx context.Context, noTransaction string) (model.GetTransactionProductResponse, error)
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

func (service *transactionService) FindAll(ctx context.Context) ([]model.GetTransactionProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []model.GetTransactionProductResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	transaction, err := service.TransactionRepository.FindAll(ctx, tx)
	if err != nil {
		return []model.GetTransactionProductResponse{}, err
	}

	var transactions []model.GetTransactionProductResponse
	for _, rows := range transaction {
		transactions = append(transactions, helper.ToTransactionProductResponse(rows))
	}

	return transactions, nil
}

func (service *transactionService) FindByTransaction(ctx context.Context, noTransaction string) (model.GetTransactionProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetTransactionProductResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	rows, err := service.TransactionRepository.FindByTransaction(ctx, tx, noTransaction)
	if err != nil {
		return model.GetTransactionProductResponse{}, err
	}

	return helper.ToTransactionProductResponse(rows), nil
}

func (service *transactionService) Create(ctx context.Context, request model.CreateTransactionRequest) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	createdAt, err := time.Parse(service.date, time.Now().Format(service.date))
	if err != nil {
		return err
	}

	transaction := entity.Transaction{
		IdProduct:     request.IdProduct,
		CustomerName:  request.CustomerName,
		NoTransaction: request.NoTransaction,
		Quantity:      request.Quantity,
		CreatedAt:     createdAt,
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
	defer helper.CommitOrRollback(tx)

	var transactions []entity.Transaction

	for _, transaction := range data {
		idProduct, _ := strconv.Atoi(transaction[0])
		quantity, _ := strconv.Atoi(transaction[3])
		createdAt, _ := time.Parse(service.date, transaction[4]+" 00:00:00")
		createdAt.Add(time.Hour*time.Duration(1) +
			time.Minute*time.Duration(1) +
			time.Second*time.Duration(1))
		transactions = append(transactions, entity.Transaction{
			IdProduct:     uint64(idProduct),
			CustomerName:  transaction[1],
			NoTransaction: transaction[2],
			Quantity:      int32(quantity),
			CreatedAt:     createdAt,
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
	defer helper.CommitOrRollback(tx)

	// Find relation is exists
	_, err = service.ProductRepository.FindById(ctx, tx, request.IdProduct)
	if err != nil {
		return err
	}

	// Find Transaction by number transaction
	transaction, err := service.TransactionRepository.FindByTransaction(ctx, tx, request.NoTransaction)
	if err != nil {
		return err
	}

	transaction.ProductId = request.IdProduct
	transaction.TransactionCustomerName = request.CustomerName
	transaction.TransactionQuantity = request.Quantity

	entityTransaction := entity.Transaction{
		IdTransaction: transaction.TransactionId,
		IdProduct:     transaction.ProductId,
		CustomerName:  transaction.TransactionCustomerName,
		Quantity:      transaction.TransactionQuantity,
		NoTransaction: transaction.TransactionNo,
		CreatedAt:     transaction.TransactionCreatedAt,
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
	defer helper.CommitOrRollback(tx)

	rows, err := service.TransactionRepository.FindByTransaction(ctx, tx, noTransaction)
	if err != nil {
		return err
	}

	err = service.TransactionRepository.Delete(ctx, tx, rows.TransactionNo)
	if err != nil {
		return err
	}

	return nil
}
