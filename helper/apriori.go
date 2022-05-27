package helper

import (
	"apriori/entity"
	"apriori/model"
	"errors"
	"sort"
	"strconv"
	"strings"
)

func FindFirstItemSet(transactionsSet []entity.Transaction, minimumSupport int) ([]model.GetTransactionResponses, map[string]float64, []string) {
	// Generate all product
	var transactions []model.GetTransactionResponses
	for _, transaction := range transactionsSet {
		productName := strings.Split(transaction.ProductName, ", ")
		transactions = append(transactions, model.GetTransactionResponses{
			ProductName: productName,
		})
	}

	// Count all every product name
	var productName = make(map[string]float64)
	for _, value := range transactions {
		for _, product := range value.ProductName {
			product = strings.ToLower(product)
			productName[product] = productName[product] + 1
		}
	}

	// Finding one item set
	var propertyProduct []string
	for nameOfProduct, total := range productName {
		support := total / float64(len(transactionsSet)) * 100
		if support >= float64(minimumSupport) {
			supportValue := strconv.Itoa(int(support))
			totalValue := strconv.Itoa(int(total))

			propertyProduct = append(propertyProduct, nameOfProduct+":"+supportValue+","+totalValue)
		}
	}

	return transactions, productName, propertyProduct
}
func HandleMapsProblem(propertyProduct []string) ([]string, []float64, []int32) {
	var oneSet []string
	var support []float64
	var totalTransaction []int32

	sort.Strings(propertyProduct)

	for i := 0; i < len(propertyProduct); i++ {
		// Split property
		nameOfProduct := strings.Split(propertyProduct[i], ":")
		transaction := strings.Split(nameOfProduct[1], ",")

		// Insert product name
		oneSet = append(oneSet, nameOfProduct[0])

		// Convert and insert support
		number, _ := strconv.Atoi(transaction[0])
		support = append(support, float64(number))

		// Convert and insert total transaction
		transactionNumber, _ := strconv.Atoi(transaction[1])
		totalTransaction = append(totalTransaction, int32(transactionNumber))
	}

	return oneSet, support, totalTransaction
}

func FindConfidence(apriori []model.GetAprioriResponses, productName map[string]float64) []model.GetAprioriResponses {
	var confidence []model.GetAprioriResponses
	for _, value := range apriori {
		if value.Iterate == apriori[len(apriori)-1].Iterate {
			if val, ok := productName[value.ItemSet[0]]; ok {
				confidence = append(confidence, model.GetAprioriResponses{
					ItemSet:     value.ItemSet,
					Support:     value.Support,
					Iterate:     value.Iterate,
					Transaction: value.Transaction,
					Confidence:  float64(value.Transaction) / val * 100,
				})
			}
		}
	}

	return confidence
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

func FindDiscount(apriori []model.GetAprioriResponses, minDiscount float64, maxDiscount float64) ([]model.GetAprioriResponses, error) {
	if maxDiscount < minDiscount {
		return []model.GetAprioriResponses{}, errors.New("the maximum discount cannot be less than the minimum discount")
	}
	var discounts []model.GetAprioriResponses
	var calculateDiscount = (maxDiscount - minDiscount) / float64(len(apriori))

	// Sorting if the value is greater, then the discount given will be large
	sort.Slice(apriori, func(i, j int) bool {
		return apriori[i].Confidence < apriori[j].Confidence
	})

	for _, value := range apriori {
		minDiscount += calculateDiscount
		discounts = append(discounts, model.GetAprioriResponses{
			ItemSet:     value.ItemSet,
			Support:     value.Support,
			Iterate:     value.Iterate,
			Transaction: value.Transaction,
			Confidence:  value.Confidence,
			Discount:    minDiscount,
		})
	}

	return discounts, nil
}

func FilterCandidateInSlice(dataTemp [][]string) [][]string {
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

	return dataTemp
}
