package service

import (
	"apriori/helper"
	"apriori/model"
	"apriori/repository"
	"context"
	"database/sql"
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

func (service aprioriService) Generate(ctx context.Context, minimumSupport int) ([]model.GetAprioriResponses, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	var apriori []model.GetAprioriResponses

	// Get all transaction from database
	transactionsSet, err := service.TransactionRepository.FindItemSet(ctx, tx)
	if err != nil {
		return nil, err
	}

	// Find first item set
	transactions, productName, propertyProduct := helper.FindFirstItemSet(transactionsSet, minimumSupport)

	// Handle random maps problem
	oneSet, support, totalTransaction := helper.HandleMapsProblem(propertyProduct)

	// Get one item set
	for i := 0; i < len(oneSet); i++ {
		apriori = append(apriori, model.GetAprioriResponses{
			ItemSet:     []string{oneSet[i]},
			Support:     support[i],
			Iterate:     1,
			Transaction: totalTransaction[i],
		})
	}

	// Looking for more than one item set
	var iterate int
	var dataTemp [][]string
	for {
		for i := 0; i < len(oneSet)-iterate; i++ {
			for j := i + 1; j < len(oneSet); j++ {
				var iterateCandidate []string

				iterateCandidate = append(iterateCandidate, oneSet[i])
				for z := 1; z <= iterate; z++ {
					iterateCandidate = append(iterateCandidate, oneSet[i+z])
				}
				iterateCandidate = append(iterateCandidate, oneSet[j])

				dataTemp = append(dataTemp, iterateCandidate)
			}
		}
		// Filter when the slice has duplicate values
		var cleanValues [][]string
		for i := 0; i < len(dataTemp); i++ {
			if isDuplicate := helper.IsDuplicate(dataTemp[i]); !isDuplicate {
				cleanValues = append(cleanValues, dataTemp[i])
			}
		}
		dataTemp = cleanValues

		// Filter candidates by comparing slice to slice
		dataTemp = helper.FilterCandidateInSlice(dataTemp)

		// Find item set by minimum support
		for i := 0; i < len(dataTemp); i++ {
			countCandidates := helper.FindCandidate(dataTemp[i], transactions)
			result := float64(countCandidates) / float64(len(transactionsSet)) * 100
			if result >= float64(minimumSupport) {
				apriori = append(apriori, model.GetAprioriResponses{
					ItemSet:     dataTemp[i],
					Support:     result,
					Iterate:     int32(iterate) + 2,
					Transaction: int32(countCandidates),
				})
			}
		}

		// Convert dataTemp slice of slice to one slice
		var test []string
		for i := 0; i < len(dataTemp); i++ {
			test = append(test, dataTemp[i]...)
		}
		oneSet = test

		// After finish operating, then clean the array
		dataTemp = [][]string{}

		// if nothing else is sent to the candidate, then break it
		if int32(iterate)+2 > apriori[len(apriori)-1].Iterate {
			break
		}

		// Add increment, if any candidate is submitted
		iterate++
	}

	// Find Association rules
	// Set confidence
	confidence := helper.FindConfidence(apriori, productName)

	// Set discount
	discount, err := helper.FindDiscount(confidence, 10, 15)
	if err != nil {
		return nil, err
	}

	// Remove last element in apriori as many association rules
	for i := 0; i < len(discount); i++ {
		apriori = apriori[:len(apriori)-1]
	}

	// Insert new apriori with discount and confidence
	for i := 0; i < len(discount); i++ {
		apriori = append(apriori, model.GetAprioriResponses{
			ItemSet:     discount[i].ItemSet,
			Support:     discount[i].Support,
			Iterate:     discount[i].Iterate,
			Transaction: discount[i].Transaction,
			Confidence:  discount[i].Confidence,
			Discount:    discount[i].Discount,
		})
	}

	return apriori, nil
}
