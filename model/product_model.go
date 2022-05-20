package model

type CreateProductRequest struct {
	Code        string `json:"code" binding:"required,min=2,max=10"`
	Name        string `json:"name" binding:"required,min=6,max=100"`
	Description string `json:"description" binding:"omitempty,max=100"`
}

type UpdateProductRequest struct {
	IdProduct   uint64 `json:"id_product"`
	Code        string `json:"code" binding:"required,min=2,max=10"`
	Name        string `json:"name" binding:"required,min=6,max=100"`
	Description string `json:"description" binding:"omitempty,max=100"`
}

type GetProductResponse struct {
	IdProduct   uint64 `json:"id_product"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
