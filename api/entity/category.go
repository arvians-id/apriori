package entity

import (
	"apriori/model"
	"time"
)

type Category struct {
	IdCategory int
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (category *Category) ToCategoryResponse() *model.GetCategoryResponse {
	return &model.GetCategoryResponse{
		IdCategory: category.IdCategory,
		Name:       category.Name,
		CreatedAt:  category.CreatedAt.String(),
		UpdatedAt:  category.UpdatedAt.String(),
	}
}
