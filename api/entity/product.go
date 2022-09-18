package entity

import (
	"time"
)

type Product struct {
	IdProduct   int       `json:"id_product"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Price       int       `json:"price"`
	Category    string    `json:"category"`
	IsEmpty     bool      `json:"is_empty"`
	Mass        int       `json:"mass"`
	Image       *string   `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductRecommendation struct {
	AprioriId          int     `json:"apriori_id"`
	AprioriCode        string  `json:"apriori_code"`
	AprioriItem        string  `json:"apriori_item"`
	AprioriDiscount    float64 `json:"apriori_discount"`
	ProductTotalPrice  int     `json:"product_total_price"`
	PriceAfterDiscount int     `json:"price_discount"`
	Image              *string `json:"apriori_image"`
	Mass               int     `json:"mass,omitempty"`
	Description        *string `json:"apriori_description"`
}
