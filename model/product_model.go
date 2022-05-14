package model

import "time"

type CreateProductRequest struct {
	Code        string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UpdateProductRequest struct {
	IdProduct   uint64
	Code        string
	Name        string
	Description string
	UpdatedAt   time.Time
}

type GetProductResponse struct {
	IdProduct   uint64
	Code        string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
