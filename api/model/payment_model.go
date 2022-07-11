package model

type PayPaymentRequest struct {
	IdProduct         string      `json:"id_product,omitempty"`
	Code              interface{} `json:"code,omitempty"`
	Name              string      `json:"name,omitempty"`
	Price             int         `json:"price,omitempty"`
	Image             string      `json:"image,omitempty"`
	Qty               int         `json:"qty,omitempty"`
	TotalPricePerItem int         `json:"total_price_per_item,omitempty"`
}
