package service

import (
	"apriori/helper"
	"apriori/model"
	"apriori/repository"
	"context"
	"database/sql"
	"sort"
	"strconv"
	"strings"
)

type AprioriService interface {
	Generate(ctx context.Context, support int) ([]model.GetAprioriResponses, error)
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

func (service aprioriService) Generate(ctx context.Context, support int) ([]model.GetAprioriResponses, error) {
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

	// Count all every product name
	var productName = make(map[string]float64)
	for _, value := range transactions {
		for _, text := range value.ProductName {
			text = strings.ToLower(text)
			productName[text] = productName[text] + 1
		}
	}

	// Finding one item set
	minimumSupport := support
	var apriori []model.GetAprioriResponses
	var tes []string
	for key, value := range productName {
		result := value / float64(len(transactionsSet)) * 100
		if result >= float64(minimumSupport) {
			str := strconv.Itoa(int(result))
			tes = append(tes, key+":"+str)
		}
	}

	// Handle random maps problem
	sort.Strings(tes)
	var oneSet []string
	var val []float64
	for i := 0; i < len(tes); i++ {
		str := strings.Split(tes[i], ":")
		number, _ := strconv.Atoi(str[1])
		oneSet = append(oneSet, str[0])
		val = append(val, float64(number))
	}

	// Get one item set
	for i := 0; i < len(oneSet); i++ {
		apriori = append(apriori, model.GetAprioriResponses{
			ItemSet: []string{oneSet[i]},
			Support: val[i],
			Number:  1,
		})
	}

	// Looking for more than one item set
	var inc int
	var dataTemp [][]string
	for {
		for i := 0; i < len(oneSet)-inc; i++ {
			for j := i + 1; j < len(oneSet); j++ {
				var v []string

				v = append(v, oneSet[i])
				for z := 1; z <= inc; z++ {
					v = append(v, oneSet[i+z])
				}
				v = append(v, oneSet[j])

				dataTemp = append(dataTemp, v)
			}
		}
		// Filter when the slice has duplicate values
		var temp [][]string
		for i := 0; i < len(dataTemp); i++ {
			isDuplicate := IsDuplicate(dataTemp[i])
			if !isDuplicate {
				temp = append(temp, dataTemp[i])
			}
		}
		dataTemp = temp

		// Filter candidates by comparing slice to slice
		for k := 0; k < 2; k++ {
			for i := 0; i < len(dataTemp)-1; i++ {
				for j := i + 1; j < len(dataTemp); j++ {
					filter := FilterCandidate(dataTemp[i], dataTemp[j])
					if filter {
						dataTemp = append(dataTemp[:i], dataTemp[j+1:]...)
					}
				}
			}
		}

		// Find item set by minimum support
		for i := 0; i < len(dataTemp); i++ {
			tests := FindCandidate(dataTemp[i], transactions)
			result := float64(tests) / float64(len(transactionsSet)) * 100
			if result >= float64(minimumSupport) {
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

		if int32(inc)+2 > apriori[len(apriori)-1].Number {
			break
		}
		inc++
	}

	return apriori, nil
}

func IsDuplicate(array []string) bool {
	visited := make(map[string]bool, 0)
	for i := 0; i < len(array); i++ {
		if visited[array[i]] == true {
			return true
		} else {
			visited[array[i]] = true
		}
	}
	return false
}

func FilterCandidate(first, second []string) bool {
	sort.Strings(first)
	sort.Strings(second)

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
		results := make([]string, 0) // slice to store the result

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
