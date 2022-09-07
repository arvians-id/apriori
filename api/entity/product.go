package entity

import (
	"apriori/model"
	"time"
)

type Product struct {
	IdProduct   int
	Code        string
	Name        string
	Description string
	Price       int
	Category    string
	IsEmpty     int
	Mass        int
	Image       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (product *Product) ToProductResponse() *model.GetProductResponse {
	return &model.GetProductResponse{
		IdProduct:   product.IdProduct,
		Code:        product.Code,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category,
		IsEmpty:     product.IsEmpty,
		Mass:        product.Mass,
		Image:       product.Image,
		CreatedAt:   product.CreatedAt.String(),
		UpdatedAt:   product.UpdatedAt.String(),
	}
}
