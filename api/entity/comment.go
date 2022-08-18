package entity

import (
	"database/sql"
	"time"
)

type Comment struct {
	IdComment   int
	UserOrderId int
	ProductCode string
	Description sql.NullString
	Tag         sql.NullString
	Rating      int
	CreatedAt   time.Time
	UserId      int
	UserName    string
}

type RatingFromComment struct {
	Rating        int
	ResultRating  int
	ResultComment int
}
