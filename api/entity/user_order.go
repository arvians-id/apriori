package entity

type UserOrder struct {
	IdOrder        int
	PayloadId      int
	Code           string
	Name           string
	Price          int64
	Image          string
	Quantity       int
	TotalPriceItem int64
}

type UserOrderRelationByUserId struct {
	UserOrder         UserOrder
	OrderId           string
	TransactionStatus string
}
