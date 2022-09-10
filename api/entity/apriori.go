package entity

import (
	"apriori/model"
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
	IsActive    bool
	Description sql.NullString
	Image       string
	CreatedAt   time.Time
}

func (apriori *Apriori) ToAprioriResponse() *model.GetAprioriResponse {
	return &model.GetAprioriResponse{
		IdApriori:   apriori.IdApriori,
		Code:        apriori.Code,
		Item:        apriori.Item,
		Discount:    apriori.Discount,
		Support:     apriori.Support,
		Confidence:  apriori.Confidence,
		RangeDate:   apriori.RangeDate,
		IsActive:    apriori.IsActive,
		Description: apriori.Description.String,
		Image:       apriori.Image,
		CreatedAt:   apriori.CreatedAt.String(),
	}
}
