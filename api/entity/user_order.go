package entity

type UserOrder struct {
	IdOrder        int      `json:"id_order"`
	PayloadId      int      `json:"payload_id"`
	Code           *string  `json:"code"`
	Name           *string  `json:"name"`
	Price          *int64   `json:"price"`
	Image          *string  `json:"image"`
	Quantity       *int     `json:"quantity"`
	TotalPriceItem *int64   `json:"total_price_item"`
	Payment        *Payment `json:"payment"`
}
