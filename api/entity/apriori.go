package entity

import "time"

type Apriori struct {
	IdApriori  uint64
	Code       string
	Item       string
	Discount   float64
	Support    float64
	Confidence float64
	RangeDate  string
	IsActive   bool
	CreatedAt  time.Time
}
