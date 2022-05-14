package model

type CreateTransactionRequest struct {
	IdProduct     uint64 `json:"id_product"`
	CustomerName  string `json:"customer_name"`
	NoTransaction string `json:"no_transaction"`
	Quantity      int32  `json:"quantity"`
	CreatedAt     string `json:"created_at"`
}

type UpdateTransactionRequest struct {
	IdTransaction uint64 `json:"id_transaction"`
	IdProduct     uint64 `json:"id_product"`
	CustomerName  string `json:"customer_name"`
	Quantity      int32  `json:"quantity"`
}

type GetTransactionResponse struct {
	IdTransaction uint64 `json:"id_transaction"`
	IdProduct     uint64 `json:"id_product"`
	CustomerName  string `json:"customer_name"`
	NoTransaction string `json:"no_transaction"`
	Quantity      int32  `json:"quantity"`
	CreatedAt     string `json:"created_at"`
}
