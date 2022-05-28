package service

import (
	"apriori/lib"
	"apriori/model"
	"apriori/repository"
	"apriori/utils"
	"context"
	"database/sql"
)

type AprioriService interface {
	Generate(ctx context.Context, request model.GenerateAprioriRequest) ([]model.GetAprioriResponses, error)
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

func (service aprioriService) Generate(ctx context.Context, request model.GenerateAprioriRequest) ([]model.GetAprioriResponses, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	var apriori []model.GetAprioriResponses

	// Get all transaction from database
	transactionsSet, err := service.TransactionRepository.FindItemSet(ctx, tx, request.StartDate, request.EndDate)
	if err != nil {
		return nil, err
	}

	// Find first item set
	transactions, productName, propertyProduct := lib.FindFirstItemSet(transactionsSet, int(request.MinimumSupport))

	// Handle random maps problem
	oneSet, support, totalTransaction := lib.HandleMapsProblem(propertyProduct)

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
			if isDuplicate := lib.IsDuplicate(dataTemp[i]); !isDuplicate {
				cleanValues = append(cleanValues, dataTemp[i])
			}
		}
		dataTemp = cleanValues
		// Filter candidates by comparing slice to slice
		dataTemp = lib.FilterCandidateInSlice(dataTemp)

		// Find item set by minimum support
		for i := 0; i < len(dataTemp); i++ {
			countCandidates := lib.FindCandidate(dataTemp[i], transactions)
			result := float64(countCandidates) / float64(len(transactionsSet)) * 100
			if result >= float64(request.MinimumSupport) {
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
	confidence := lib.FindConfidence(apriori, productName)

	// Set discount
	discount := lib.FindDiscount(confidence, float64(request.MinimumDiscount), float64(request.MaximumDiscount))

	// Remove last element in apriori as many association rules
	for i := 0; i < len(discount); i++ {
		apriori = apriori[:len(apriori)-1]
	}

	// Replace the last item set and add discount and confidence
	for i := 0; i < len(discount); i++ {
		if discount[i].Confidence >= float64(request.MinimumConfidence) {
			apriori = append(apriori, model.GetAprioriResponses{
				ItemSet:     discount[i].ItemSet,
				Support:     discount[i].Support,
				Iterate:     discount[i].Iterate,
				Transaction: discount[i].Transaction,
				Confidence:  discount[i].Confidence,
				Discount:    discount[i].Discount,
				Description: "Eligible",
			})
		} else {
			apriori = append(apriori, model.GetAprioriResponses{
				ItemSet:     discount[i].ItemSet,
				Support:     discount[i].Support,
				Iterate:     discount[i].Iterate,
				Transaction: discount[i].Transaction,
				Confidence:  discount[i].Confidence,
				Description: "Not Eligible",
			})
		}
	}

	return apriori, nil
}
