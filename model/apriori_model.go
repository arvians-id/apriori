package model

import "time"

type CreateAprioriRequest struct {
	Item       string
	Discount   int32
	Support    int32
	Confidence int32
	RangeDate  string
	Counter    uint64
	CreatedAt  time.Time
}

type GetAprioriResponse struct {
	IdApriori  uint64
	Item       string
	Discount   int32
	Support    int32
	Confidence int32
	RangeDate  string
	Counter    uint64
	CreatedAt  time.Time
}
