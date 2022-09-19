package entity

import (
	"time"
)

type Apriori struct {
	IdApriori   int        `json:"id_apriori"`
	Code        string     `json:"code"`
	Item        string     `json:"item"`
	Discount    float64    `json:"discount"`
	Support     float64    `json:"support"`
	Confidence  float64    `json:"confidence"`
	RangeDate   string     `json:"range_date"`
	IsActive    bool       `json:"is_active"`
	Description *string    `json:"description"`
	Mass        int        `json:"mass"`
	Image       *string    `json:"image"`
	CreatedAt   time.Time  `json:"created_at"`
	UserOrder   *UserOrder `json:"user_order"`
}

type GenerateApriori struct {
	ItemSet     []string `json:"item_set"`
	Support     float64  `json:"support"`
	Iterate     int      `json:"iterate"`
	Transaction int      `json:"transaction"`
	Confidence  float64  `json:"confidence"`
	Discount    float64  `json:"discount"`
	Description string   `json:"description"`
	RangeDate   string   `json:"range_date"`
}
