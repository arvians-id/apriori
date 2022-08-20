package entity

import (
	"database/sql"
	"time"
)

type Apriori struct {
	IdApriori   int
	Code        string
	Item        string
	Discount    float64
	Support     float64
	Confidence  float64
	RangeDate   string
	IsActive    int
	Description sql.NullString
	Image       string
	CreatedAt   time.Time
}
