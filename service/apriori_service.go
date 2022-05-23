package service

import (
	"apriori/helper"
	"apriori/model"
	"apriori/repository"
	"context"
	"database/sql"
)

type AprioriService interface {
	Generate(ctx context.Context) ([]model.GetProductTransactionResponse, error)
}

type aprioriService struct {
	ProductRepository repository.ProductRepository
	DB                *sql.DB
}

func NewAprioriService(productRepository *repository.ProductRepository, db *sql.DB) AprioriService {
	return &aprioriService{
		ProductRepository: *productRepository,
		DB:                db,
	}
}

func (service aprioriService) Generate(ctx context.Context) ([]model.GetProductTransactionResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []model.GetProductTransactionResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	transactionsSet, err := service.ProductRepository.FindItemSet(ctx, tx)
	if err != nil {
		return []model.GetProductTransactionResponse{}, err
	}

	var transactions []model.GetProductTransactionResponse
	for _, transaction := range transactionsSet {
		transactions = append(transactions, helper.ToProductTransactionResponse(transaction))
	}

	return transactions, nil
}
