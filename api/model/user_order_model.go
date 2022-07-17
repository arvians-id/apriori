package model

type GetUserOrderResponse struct {
	IdOrder        uint64 `json:"id_order,omitempty"`
	PayloadId      uint64 `json:"payload_id,omitempty"`
	Code           string `json:"code,omitempty"`
	Name           string `json:"name,omitempty"`
	Price          int64  `json:"price,omitempty"`
	Image          string `json:"image,omitempty"`
	Quantity       int    `json:"quantity,omitempty"`
	TotalPriceItem int64  `json:"total_price_item,omitempty"`
}
