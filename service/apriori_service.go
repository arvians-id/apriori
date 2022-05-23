package service

import (
	"apriori/helper"
	"apriori/model"
	"apriori/repository"
	"context"
	"database/sql"
)

type AprioriService interface {
	Generate(ctx context.Context) ([]model.GetTransactionResponse, error)
}

type aprioriService struct {
	TransactionRepository repository.TransactionRepository
	DB                    *sql.DB
}

func NewAprioriService(transactionRepository *repository.TransactionRepository, db *sql.DB) AprioriService {
	return &aprioriService{
		TransactionRepository: *transactionRepository,
		DB:                    db,
	}
}

func (service aprioriService) Generate(ctx context.Context) ([]model.GetTransactionResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []model.GetTransactionResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	transactionsSet, err := service.TransactionRepository.FindItemSet(ctx, tx)
	if err != nil {
		return []model.GetTransactionResponse{}, err
	}

	var transactions []model.GetTransactionResponse
	for _, _ = range transactionsSet {
		//transactions = append(transactions, helper.ToTransactionResponse(transaction))
	}

	return transactions, nil
}
