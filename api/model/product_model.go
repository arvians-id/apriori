package model

type CreateProductRequest struct {
	Code        string `form:"code" binding:"required,min=2,max=10"`
	Name        string `form:"name" binding:"required,min=6,max=100"`
	Description string `form:"description" binding:"omitempty,max=100"`
	Price       int    `form:"price" binding:"min=0"`
	Category    string `form:"category" binding:"required,max=100"`
	Mass        int    `form:"mass"`
	Image       string `form:"-"`
}

type UpdateProductRequest struct {
	IdProduct   int    `form:"-"`
	Code        string `form:"-"`
	Name        string `form:"name" binding:"required,min=6,max=100"`
	Description string `form:"description" binding:"omitempty,max=100"`
	Price       int    `form:"price" binding:"min=0"`
	Category    string `form:"category" binding:"required,max=100"`
	IsEmpty     bool   `form:"is_empty"`
	Mass        int    `form:"mass"`
	Image       string `form:"-"`
}

type GetProductTransactionResponse struct {
	Code        string  `json:"code"`
	ProductName string  `json:"product_name"`
	Transaction int32   `json:"transaction"`
	Support     float32 `json:"support"`
}
