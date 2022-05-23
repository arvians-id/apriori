package model

type CreateTransactionRequest struct {
	ProductName   string `json:"product_name" binding:"required,max=256"`
	CustomerName  string `json:"customer_name" binding:"required,max=100"`
	NoTransaction string `json:"no_transaction" binding:"required,max=100"`
}

type CreateTransactionFromFileRequest struct {
	File uint64 `json:"file" binding:"required"`
}

type UpdateTransactionRequest struct {
	ProductName   string `json:"product_name" binding:"required,max=256"`
	CustomerName  string `json:"customer_name" binding:"required,max=100"`
	NoTransaction string `json:"no_transaction" binding:"required,max=100"`
}

type GetTransactionResponse struct {
	IdTransaction uint64 `json:"id_transaction"`
	ProductName   string `json:"product_name"`
	CustomerName  string `json:"customer_name"`
	NoTransaction string `json:"no_transaction"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}
