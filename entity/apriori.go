package entity

import "time"

type Apriori struct {
	IdApriori  uint64
	Item       string
	Discount   int32
	Support    int32
	Confidence int32
	RangeDate  string
	Counter    uint64
	CreatedAt  time.Time
}

type RuleAssociation struct {
	ItemSets   []ProductTransaction
	Confidence int32
}

type ProductTransaction struct {
	Code        string
	ProductName string
	Transaction int32
	Support     float32
}
