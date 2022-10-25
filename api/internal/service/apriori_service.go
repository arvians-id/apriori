package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/arvians-id/apriori/internal/http/presenter/request"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/arvians-id/apriori/internal/repository"
	"github.com/arvians-id/apriori/util"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

type AprioriServiceImpl struct {
	TransactionRepository repository.TransactionRepository
	AprioriRepository     repository.AprioriRepository
	ProductRepository     repository.ProductRepository
	DB                    *sql.DB
}

func NewAprioriService(
	transactionRepository *repository.TransactionRepository,
	productRepository *repository.ProductRepository,
	aprioriRepository *repository.AprioriRepository,
	db *sql.DB,
) AprioriService {
	return &AprioriServiceImpl{
		TransactionRepository: *transactionRepository,
		AprioriRepository:     *aprioriRepository,
		ProductRepository:     *productRepository,
		DB:                    db,
	}
}

func (service *AprioriServiceImpl) FindAll(ctx context.Context) ([]*model.Apriori, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[AprioriService][FindAll] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	apriories, err := service.AprioriRepository.FindAll(ctx, tx)
	if err != nil {
		log.Println("[AprioriService][FindAll] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return apriories, nil
}

func (service *AprioriServiceImpl) FindAllByActive(ctx context.Context) ([]*model.Apriori, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[AprioriService][FindAllByActive] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	apriories, err := service.AprioriRepository.FindAllByActive(ctx, tx)
	if err != nil {
		log.Println("[AprioriService][FindAllByActive][FindAllByActive] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return apriories, nil
}

func (service *AprioriServiceImpl) FindAllByCode(ctx context.Context, code string) ([]*model.Apriori, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[AprioriService][FindAllByCode] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	apriories, err := service.AprioriRepository.FindAllByCode(ctx, tx, code)
	if err != nil {
		log.Println("[AprioriService][FindAllByCode][FindAllByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return apriories, nil
}

func (service *AprioriServiceImpl) FindByCodeAndId(ctx context.Context, code string, id int) (*model.ProductRecommendation, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[AprioriService][FindByCodeAndId] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	apriori, err := service.AprioriRepository.FindByCodeAndId(ctx, tx, code, id)
	if err != nil {
		log.Println("[AprioriService][FindByCodeAndId][FindByCodeAndId] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var totalPrice, mass int
	productNames := strings.Split(apriori.Item, ",")
	for _, productName := range productNames {
		product, _ := service.ProductRepository.FindByName(ctx, tx, util.UpperWords(productName))
		totalPrice += product.Price
		mass += product.Mass
	}

	return &model.ProductRecommendation{
		AprioriId:          apriori.IdApriori,
		AprioriCode:        apriori.Code,
		AprioriItem:        apriori.Item,
		AprioriDiscount:    apriori.Discount,
		ProductTotalPrice:  totalPrice,
		PriceDiscount:      totalPrice - (totalPrice * int(apriori.Discount) / 100),
		AprioriImage:       apriori.Image,
		Mass:               mass,
		AprioriDescription: apriori.Description,
	}, nil
}

func (service *AprioriServiceImpl) Create(ctx context.Context, requests []*request.CreateAprioriRequest) error {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[AprioriService][Create] problem in db transaction, err: ", err.Error())
		return err
	}
	defer util.CommitOrRollback(tx)

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[AprioriService][Create] problem in parsing to time, err: ", err.Error())
		return err
	}

	var aprioriRequests []*model.Apriori
	code := util.RandomString(10)
	for _, requestItem := range requests {
		image := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/assets/%s", os.Getenv("AWS_BUCKET"), os.Getenv("AWS_REGION"), "no-image.png")
		aprioriRequests = append(aprioriRequests, &model.Apriori{
			Code:       code,
			Item:       requestItem.Item,
			Discount:   requestItem.Discount,
			Support:    requestItem.Support,
			Confidence: requestItem.Confidence,
			RangeDate:  requestItem.RangeDate,
			IsActive:   false,
			Image:      &image,
			CreatedAt:  timeNow,
		})
	}

	err = service.AprioriRepository.Create(ctx, tx, aprioriRequests)
	if err != nil {
		log.Println("[AprioriService][Create][Create] problem in getting from repository, err: ", err.Error())
		return err
	}

	return nil
}

func (service *AprioriServiceImpl) Update(ctx context.Context, request *request.UpdateAprioriRequest) (*model.Apriori, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[AprioriService][Update] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	apriori, err := service.AprioriRepository.FindByCodeAndId(ctx, tx, request.Code, request.IdApriori)
	if err != nil {
		log.Println("[AprioriService][Update][FindByCodeAndId] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	image := apriori.Image
	if request.Image != "" {
		image = &request.Image
	}

	apriori.Image = image
	apriori.Description = &request.Description

	_, err = service.AprioriRepository.Update(ctx, tx, apriori)
	if err != nil {
		log.Println("[AprioriService][Update][Update] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return apriori, nil
}

func (service *AprioriServiceImpl) UpdateStatus(ctx context.Context, code string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[AprioriService][UpdateStatus] problem in db transaction, err: ", err.Error())
		return err
	}
	defer util.CommitOrRollback(tx)

	apriories, err := service.AprioriRepository.FindAllByCode(ctx, tx, code)
	if err != nil {
		log.Println("[AprioriService][UpdateStatus][FindAllByCode] problem in getting from repository, err: ", err.Error())
		return err
	}

	err = service.AprioriRepository.UpdateAllStatus(ctx, tx, false)
	if err != nil {
		log.Println("[AprioriService][UpdateStatus][UpdateAllStatus] problem in getting from repository, err: ", err.Error())
		return err
	}

	status := true
	if apriories[0].IsActive {
		status = false
	}

	err = service.AprioriRepository.UpdateStatusByCode(ctx, tx, apriories[0].Code, status)
	if err != nil {
		log.Println("[AprioriService][UpdateStatus][UpdateStatusByCode] problem in getting from repository, err: ", err.Error())
		return err
	}

	return nil
}

func (service *AprioriServiceImpl) Delete(ctx context.Context, code string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[AprioriService][Delete] problem in db transaction, err: ", err.Error())
		return err
	}
	defer util.CommitOrRollback(tx)

	apriories, err := service.AprioriRepository.FindAllByCode(ctx, tx, code)
	if err != nil {
		log.Println("[AprioriService][Delete][FindAllByCode] problem in getting from repository, err: ", err.Error())
		return err
	}

	err = service.AprioriRepository.Delete(ctx, tx, apriories[0].Code)
	if err != nil {
		log.Println("[AprioriService][Delete][Delete] problem in getting from repository, err: ", err.Error())
		return err
	}

	return nil
}

func (service *AprioriServiceImpl) Generate(ctx context.Context, request *request.GenerateAprioriRequest) ([]*model.GenerateApriori, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[AprioriService][Generate] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	var apriori []*model.GenerateApriori
	// Get all transaction from database
	transactionsSet, err := service.TransactionRepository.FindAllItemSet(ctx, tx, request.StartDate, request.EndDate)
	if err != nil {
		log.Println("[AprioriService][Generate][FindAllItemSet] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	// Find first item set
	transactions, productName, propertyProduct := util.FindFirstItemSet(transactionsSet, request.MinimumSupport)

	// Handle random maps problem
	oneSet, support, totalTransaction, isEligible, cleanSet := util.HandleMapsProblem(propertyProduct, request.MinimumSupport)

	// Get one item set
	for i := 0; i < len(oneSet); i++ {
		apriori = append(apriori, &model.GenerateApriori{
			ItemSet:     []string{oneSet[i]},
			Support:     support[i],
			Iterate:     1,
			Transaction: totalTransaction[i],
			Description: isEligible[i],
			RangeDate:   request.StartDate + " - " + request.EndDate,
		})
	}

	oneSet = cleanSet
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
			if isDuplicate := util.IsDuplicate(dataTemp[i]); !isDuplicate {
				cleanValues = append(cleanValues, dataTemp[i])
			}
		}
		dataTemp = cleanValues
		// Filter candidates by comparing slice to slice
		dataTemp = util.FilterCandidateInSlice(dataTemp)

		// Find item set by minimum support
		for i := 0; i < len(dataTemp); i++ {
			countCandidates := util.FindCandidate(dataTemp[i], transactions)
			result := float64(countCandidates) / float64(len(transactionsSet)) * 100
			if result >= request.MinimumSupport {
				apriori = append(apriori, &model.GenerateApriori{
					ItemSet:     dataTemp[i],
					Support:     math.Round(result*100) / 100,
					Iterate:     iterate + 2,
					Transaction: countCandidates,
					Description: "Eligible",
					RangeDate:   request.StartDate + " - " + request.EndDate,
				})
			} else {
				apriori = append(apriori, &model.GenerateApriori{
					ItemSet:     dataTemp[i],
					Support:     math.Round(result*100) / 100,
					Iterate:     iterate + 2,
					Transaction: countCandidates,
					Description: "Not Eligible",
					RangeDate:   request.StartDate + " - " + request.EndDate,
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

		var checkClean bool
		for _, value := range apriori {
			if value.Iterate == iterate+2 && value.Description == "Eligible" {
				checkClean = true
				break
			}
		}

		countIterate := 0
		for i := 0; i < len(apriori); i++ {
			if apriori[i].Iterate == iterate+2 {
				countIterate++
			}
		}

		if checkClean == false {
			for i := 0; i < len(apriori); i++ {
				if apriori[i].Iterate == iterate+2 {
					apriori = append(apriori[:i], apriori[i+countIterate:]...)
				}
			}
			break
		}

		// if nothing else is sent to the candidate, then break it
		if iterate+2 > apriori[len(apriori)-1].Iterate {
			break
		}

		// Add increment, if any candidate is submitted
		iterate++
	}

	// Find Association rules
	// Set confidence
	confidence := util.FindConfidence(apriori, productName, request.MinimumSupport, request.MinimumConfidence)

	// Set discount
	discount := util.FindDiscount(confidence, float64(request.MinimumDiscount), float64(request.MaximumDiscount))

	//// Remove last element in apriori as many association rules
	//for i := 0; i < len(discount); i++ {
	//	apriori = apriori[:len(apriori)-1]
	//}

	// Replace the last item set and add discount and confidence
	for i := 0; i < len(discount); i++ {
		if discount[i].Confidence >= request.MinimumConfidence {
			apriori = append(apriori, &model.GenerateApriori{
				ItemSet:     discount[i].ItemSet,
				Support:     math.Round(discount[i].Support*100) / 100,
				Iterate:     discount[i].Iterate + 1,
				Transaction: discount[i].Transaction,
				Confidence:  math.Round(discount[i].Confidence*100) / 100,
				Discount:    discount[i].Discount,
				Description: "Rules",
				RangeDate:   request.StartDate + " - " + request.EndDate,
			})
		}
	}

	return apriori, nil
}
