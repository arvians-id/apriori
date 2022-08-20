package entity

import "time"

type Apriori struct {
	IdApriori   int
	Code        string
	Item        string
	Discount    float64
	Support     float64
	Confidence  float64
	RangeDate   string
	IsActive    int
	Description *string
	Image       string
	CreatedAt   time.Time
}
