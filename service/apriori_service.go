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
	var productName = make(map[string]float64)
	for _, value := range transactions {
		for _, text := range value.ProductName {
			text = strings.ToLower(text)
			productName[text] = productName[text] + 1
		}
	}

	// Finding one item set
	minimumSupport := 40
	var apriori []model.GetAprioriResponses
	var tes []string
	for key, value := range productName {
		result := value / float64(len(transactionsSet)) * 100
		if result >= float64(minimumSupport) {
			str := strconv.Itoa(int(result))
			tes = append(tes, key+":"+str)
		}
	}

	// Handle random map problem
	sort.Strings(tes)

	var oneSet []string
	var val []float64
	for i := 0; i < len(tes); i++ {
		str := strings.Split(tes[i], ":")
		number, _ := strconv.Atoi(str[1])
		oneSet = append(oneSet, str[0])
		val = append(val, float64(number))
	}

	for i := 0; i < len(oneSet); i++ {
		apriori = append(apriori, model.GetAprioriResponses{
			ItemSet: []string{oneSet[i]},
			Support: val[i],
			Number:  1,
		})
	}

	// Finding a more than one item set
	var inc int
	var isEnded = true
	var dataTemp [][]string
	for isEnded {
		for i := 0; i < len(oneSet)-inc; i++ {
			for j := 0; j < len(oneSet); j++ {
				if j > i {
					if inc == 0 {
						dataTemp = append(dataTemp, []string{oneSet[i], oneSet[j]})
					} else if inc == 1 {
						if oneSet[i] != oneSet[i+1] && oneSet[i] != oneSet[j] && oneSet[i+1] != oneSet[j] {
							dataTemp = append(dataTemp, []string{oneSet[i], oneSet[i+1], oneSet[j]})
						}
					} else if inc == 2 {
						if oneSet[i] != oneSet[i+1] &&
							oneSet[i] != oneSet[j] &&
							oneSet[i+1] != oneSet[j] &&
							oneSet[i] != oneSet[i+2] &&
							oneSet[i+1] != oneSet[i+2] &&
							oneSet[j] != oneSet[i+2] {
							dataTemp = append(dataTemp, []string{oneSet[i], oneSet[i+1], oneSet[i+2], oneSet[j]})
						}
					}
				}
			}
		}

		// Filter Candidate
		var sets [][]string
		for i := 0; i < len(dataTemp)-1; i++ {
			var bol = false
			if i == 0 {
				sets = append(sets, dataTemp[0])
			}
			for j := 0; j < len(dataTemp); j++ {
				if j > i {
					filter := FilterCandidate(dataTemp[i], dataTemp[j])
					if !filter {
						bol = true
					} else {
						dataTemp = append(dataTemp[:i], dataTemp[j+1:]...)
					}
				}
			}
			if bol {
				sets = append(sets, dataTemp[i+1])
			}
			bol = false
		}
		dataTemp = sets

		// Item Set Operation
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
		if int32(inc+2) > apriori[len(apriori)-1].Number {
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

	exists = make(map[string]bool)
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
