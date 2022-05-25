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
	Generate(ctx context.Context) ([]model.GetAprioriResponses, error)
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

//func (service aprioriService) Generate(ctx context.Context) ([]model.GetAprioriResponses, error) {
//	tx, err := service.DB.Begin()
//	if err != nil {
//		return nil, err
//	}
//	defer helper.CommitOrRollback(tx)
//
//	transactionsSet, err := service.TransactionRepository.FindItemSet(ctx, tx)
//	if err != nil {
//		return nil, err
//	}
//
//	// Generate all product
//	var transactions []model.GetTransactionResponses
//	for _, transaction := range transactionsSet {
//		texts := strings.Split(transaction.ProductName, ", ")
//		transactions = append(transactions, model.GetTransactionResponses{
//			ProductName: texts,
//		})
//	}
//
//	// Count all every product
//	productName := make(map[string]float32)
//	for _, value := range transactions {
//		for _, text := range value.ProductName {
//			text = strings.ToLower(text)
//			productName[text] = productName[text] + 1
//		}
//	}
//
//	// Count Minimum Support
//	// One item set
//	var apriori []model.GetAprioriResponses
//	for key, value := range productName {
//		result := value / float32(len(transactionsSet)) * 100
//		if result >= float32(70) {
//			apriori = append(apriori, model.GetAprioriResponses{
//				ItemSet: []string{key},
//				Support: result,
//				Number:  1,
//			})
//		}
//	}
//
//	// Two item set
//	var twoSet [][]string
//	for i := 0; i < len(apriori); i++ {
//		for j := 0; j < len(apriori); j++ {
//			if i != j && j > i {
//				twoSet = append(twoSet, []string{apriori[i].ItemSet[0], apriori[j].ItemSet[0]})
//			}
//		}
//	}
//
//	for i := 0; i < len(twoSet); i++ {
//		test := FindCandidate(twoSet[i], transactions)
//		result := float32(test) / float32(len(transactionsSet)) * 100
//		if result >= float32(70) {
//			apriori = append(apriori, model.GetAprioriResponses{
//				ItemSet: twoSet[i],
//				Support: result,
//				Number:  2,
//			})
//		}
//	}
//	return apriori, nil
//}

func (service aprioriService) Generate(ctx context.Context) ([]model.GetAprioriResponses, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	transactionsSet, err := service.TransactionRepository.FindItemSet(ctx, tx)
	if err != nil {
		return nil, err
	}

	// Generate all product
	var transactions []model.GetTransactionResponses
	for _, transaction := range transactionsSet {
		texts := strings.Split(transaction.ProductName, ", ")
		transactions = append(transactions, model.GetTransactionResponses{
			ProductName: texts,
		})
	}

	// Count all every product
	productName := make(map[string]float32)
	for _, value := range transactions {
		for _, text := range value.ProductName {
			text = strings.ToLower(text)
			productName[text] = productName[text] + 1
		}
	}

	// Count Minimum Support
	// One item set
	var apriori []model.GetAprioriResponses
	var oneSet []string
	for key, value := range productName {
		result := value / float32(len(transactionsSet)) * 100
		if result >= float32(70) {
			oneSet = append(oneSet, key)
			apriori = append(apriori, model.GetAprioriResponses{
				ItemSet: []string{key},
				Support: result,
				Number:  1,
			})
		}
	}

	var inc int
	var isEnded = true
	var dataTemp [][]string
	for isEnded {
		// Finding candidate
		for i := 0; i < len(oneSet)-inc; i++ {
			for j := 0; j < len(oneSet); j++ {
				if j > i {
					if inc == 0 {
						dataTemp = append(dataTemp, []string{oneSet[i], oneSet[j]})
					} else if inc == 1 {
						if oneSet[i] != oneSet[i+1] && oneSet[i] != oneSet[j] && oneSet[i+1] != oneSet[j] {
							dataTemp = append(dataTemp, []string{oneSet[i], oneSet[i+1], oneSet[j]})
						}
					}
				}
			}
		}

		// Item Set Operation
		for i := 0; i < len(dataTemp); i++ {
			tests := FindCandidate(dataTemp[i], transactions)
			result := float32(tests) / float32(len(transactionsSet)) * 100
			if result >= float32(1) {
				apriori = append(apriori, model.GetAprioriResponses{
					ItemSet: dataTemp[i],
					Support: result,
					Number:  int32(inc) + 2,
				})
			}
		}

		// Convert Item Set
		var test []string
		for i := 0; i < len(dataTemp); i++ {
			test = append(test, dataTemp[i]...)
		}
		oneSet = test

		// After finish operation clear array before
		dataTemp = [][]string{}

		if inc == 1 {
			isEnded = false
		}
		inc++
	}

	return apriori, nil
}

func FilterCandidate(first, second []string) bool {
	exists := make(map[string]bool)
	for _, value := range first {
		exists[value] = true
	}
	for _, value := range second {
		if !exists[value] {
			return false
		}
	}

	return true
}

func FindCandidate(data []string, transactions []model.GetTransactionResponses) int {
	var counter int
	for _, j := range transactions {
		results := make([]string, 0) // slice tostore the result

		for i := 0; i < len(data); i++ {
			for k := 0; k < len(j.ProductName); k++ {
				if data[i] != j.ProductName[k] {
					continue
				}
				// append a value in result only if
				// it exists both in names and board
				results = append(results, data[i])
			}
			if len(results) == len(data) {
				counter++
			}
		}
	}
	return counter
}
