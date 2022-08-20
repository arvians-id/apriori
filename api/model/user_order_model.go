package model

type GetUserOrderResponse struct {
	IdOrder        int    `json:"id_order,omitempty"`
	PayloadId      int    `json:"payload_id,omitempty"`
	Code           string `json:"code,omitempty"`
	Name           string `json:"name,omitempty"`
	Price          int64  `json:"price,omitempty"`
	Image          string `json:"image,omitempty"`
	Quantity       int    `json:"quantity,omitempty"`
	TotalPriceItem int64  `json:"total_price_item,omitempty"`
}

type GetUserOrderRelationByUserIdResponse struct {
	IdOrder           int    `json:"id_order,omitempty"`
	PayloadId         int    `json:"payload_id,omitempty"`
	Code              string `json:"code,omitempty"`
	Name              string `json:"name,omitempty"`
	Price             int64  `json:"price,omitempty"`
	Image             string `json:"image,omitempty"`
	Quantity          int    `json:"quantity,omitempty"`
	TotalPriceItem    int64  `json:"total_price_item,omitempty"`
	OrderId           string `json:"order_id,omitempty"`
	TransactionStatus string `json:"transaction_status,omitempty"`
}
