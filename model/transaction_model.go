package model

type CreateTransactionRequest struct {
	IdProduct     uint64 `json:"id_product" binding:"required"`
	CustomerName  string `json:"customer_name" binding:"required,max=100"`
	NoTransaction string `json:"no_transaction" binding:"required,max=100"`
	Quantity      int32  `json:"quantity" binding:"required,numeric"`
}

type CreateTransactionFromFileRequest struct {
	File uint64 `json:"file" binding:"required"`
}

type UpdateTransactionRequest struct {
	IdProduct     uint64 `json:"id_product" binding:"required"`
	CustomerName  string `json:"customer_name" binding:"required,max=100"`
	NoTransaction string `json:"no_transaction" binding:"required,max=100"`
	Quantity      int32  `json:"quantity" binding:"required,numeric"`
}

type GetTransactionProductResponse struct {
	TransactionId           uint64 `json:"transaction_id"`
	TransactionCustomerName string `json:"transaction_customer_name"`
	TransactionNo           string `json:"transaction_no"`
	TransactionQuantity     int32  `json:"transaction_quantity"`
	TransactionCreatedAt    string `json:"transaction_created_at"`
	ProductId               uint64 `json:"product_id"`
	ProductCode             string `json:"product_code"`
	ProductName             string `json:"product_name"`
	ProductDescription      string `json:"product_description"`
}
