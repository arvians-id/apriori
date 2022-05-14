package model

type CreateProductRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UpdateProductRequest struct {
	IdProduct   uint64 `json:"id_product"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UpdatedAt   string `json:"updated_at"`
}

type GetProductResponse struct {
	IdProduct   uint64 `json:"id_product"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
