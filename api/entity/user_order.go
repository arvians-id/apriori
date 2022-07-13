package entity

type UserOrder struct {
	IdOrder        uint64
	PayloadId      uint64
	Code           string
	Name           string
	Price          int64
	Image          string
	Quantity       int
	TotalPriceItem int64
}
