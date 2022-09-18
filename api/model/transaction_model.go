package model

type CreateTransactionRequest struct {
	ProductName  string `json:"product_name" binding:"required,max=256"`
	CustomerName string `json:"customer_name" binding:"required,max=100"`
}

type CreateTransactionFromFileRequest struct {
	File uint64 `json:"file" binding:"required"`
}

type UpdateTransactionRequest struct {
	ProductName   string `json:"product_name" binding:"required,max=256"`
	CustomerName  string `json:"customer_name" binding:"required,max=100"`
	NoTransaction string
}

type GetProductNameTransactionResponse struct {
	ProductName []string `json:"product_name"`
}
