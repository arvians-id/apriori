package service

import (
	"apriori/entity"
	"apriori/lib"
	"apriori/model"
	"apriori/repository"
	"apriori/utils"
	"context"
	"database/sql"
	"time"
)

type AprioriService interface {
	FindAll(ctx context.Context) ([]model.GetAprioriResponse, error)
	FindByCode(ctx context.Context, code string) ([]model.GetAprioriResponse, error)
	Create(ctx context.Context, requests []model.CreateAprioriRequest) error
	Delete(ctx context.Context, code string) error
	Generate(ctx context.Context, request model.GenerateAprioriRequest) ([]model.GetGenerateAprioriResponse, error)
}

type aprioriService struct {
	TransactionRepository repository.TransactionRepository
	AprioriRepository     repository.AprioriRepository
	DB                    *sql.DB
	date                  string
}

func NewAprioriService(transactionRepository *repository.TransactionRepository, aprioriRepository *repository.AprioriRepository, db *sql.DB) AprioriService {
	return &aprioriService{
		TransactionRepository: *transactionRepository,
		AprioriRepository:     *aprioriRepository,
		DB:                    db,
		date:                  "2006-01-02 15:04:05",
	}
}

func (service aprioriService) FindAll(ctx context.Context) ([]model.GetAprioriResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []model.GetAprioriResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	rows, err := service.AprioriRepository.FindAll(ctx, tx)
	if err != nil {
		return []model.GetAprioriResponse{}, err
	}

	var apriories []model.GetAprioriResponse
	for _, apriori := range rows {
		apriories = append(apriories, utils.ToAprioriResponse(apriori))
	}

	return apriories, nil
}

func (service *aprioriService) FindByCode(ctx context.Context, code string) ([]model.GetAprioriResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []model.GetAprioriResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	rows, err := service.AprioriRepository.FindByCode(ctx, tx, code)
	if err != nil {
		return []model.GetAprioriResponse{}, err
	}

	var apriories []model.GetAprioriResponse
	for _, apriori := range rows {
		apriories = append(apriories, utils.ToAprioriResponse(apriori))
	}

	return apriories, nil
}

func (service *aprioriService) Create(ctx context.Context, requests []model.CreateAprioriRequest) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	createdAt, _ := time.Parse(service.date, time.Now().Format(service.date))
	if err != nil {
		return err
	}

	var apriories []entity.Apriori

	code := utils.RandomString(10)
	for _, request := range requests {
		apriories = append(apriories, entity.Apriori{
			Code:       code,
			Item:       request.Item,
			Discount:   request.Discount,
			Support:    request.Support,
			Confidence: request.Confidence,
			RangeDate:  request.RangeDate,
			IsActive:   false,
			CreatedAt:  createdAt,
		})
	}

	err = service.AprioriRepository.Create(ctx, tx, apriories)
	if err != nil {
		return err
	}

	return nil
}

func (service *aprioriService) Delete(ctx context.Context, code string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	row, err := service.AprioriRepository.FindByCode(ctx, tx, code)
	if err != nil {
		return err
	}

	err = service.AprioriRepository.Delete(ctx, tx, row[0].Code)
	if err != nil {
		return err
	}

	return nil
}

func (service *aprioriService) Generate(ctx context.Context, request model.GenerateAprioriRequest) ([]model.GetGenerateAprioriResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	var apriori []model.GetGenerateAprioriResponse

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
		apriori = append(apriori, model.GetGenerateAprioriResponse{
			ItemSet:     []string{oneSet[i]},
			Support:     support[i],
			Iterate:     1,
			Transaction: totalTransaction[i],
			RangeDate:   request.StartDate + " - " + request.EndDate,
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
				apriori = append(apriori, model.GetGenerateAprioriResponse{
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
			apriori = append(apriori, model.GetGenerateAprioriResponse{
				ItemSet:     discount[i].ItemSet,
				Support:     discount[i].Support,
				Iterate:     discount[i].Iterate,
				Transaction: discount[i].Transaction,
				Confidence:  discount[i].Confidence,
				Discount:    discount[i].Discount,
				Description: "Eligible",
				RangeDate:   request.StartDate + " - " + request.EndDate,
			})
		} else {
			apriori = append(apriori, model.GetGenerateAprioriResponse{
				ItemSet:     discount[i].ItemSet,
				Support:     discount[i].Support,
				Iterate:     discount[i].Iterate,
				Transaction: discount[i].Transaction,
				Confidence:  discount[i].Confidence,
				Description: "Not Eligible",
				RangeDate:   request.StartDate + " - " + request.EndDate,
			})
		}
	}

	return apriori, nil
}
