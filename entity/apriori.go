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

type ProductTransaction struct {
	ProductName string
	Count       int32
	Support     int32
}
