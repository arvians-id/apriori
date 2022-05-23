package service

import (
	"apriori/helper"
	"apriori/model"
	"apriori/repository"
	"context"
	"database/sql"
	"strings"
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

	minimumSupport := 70
	var transactions []model.GetTransactionResponse
	productName := make(map[string]float32)
	for _, transaction := range transactionsSet {
		transactions = append(transactions, helper.ToTransactionResponse(transaction))
		texts := strings.Split(transaction.ProductName, ", ")
		for _, text := range texts {
			text = strings.ToLower(text)
			productName[text] = productName[text] + 1
		}
	}

	temp := make(map[string]float32)
	for key, d := range productName {
		result := d / float32(len(transactionsSet)) * 100
		if result >= float32(minimumSupport) {
			temp[key] = result
		}
	}

	return transactions, nil
}
