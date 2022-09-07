package entity

import (
	"apriori/model"
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

func (comment *Comment) ToCommentResponse() *model.GetCommentResponse {
	return &model.GetCommentResponse{
		IdComment:   comment.IdComment,
		UserOrderId: comment.UserOrderId,
		ProductCode: comment.ProductCode,
		Description: comment.Description.String,
		Tag:         comment.Tag.String,
		Rating:      comment.Rating,
		CreatedAt:   comment.CreatedAt.String(),
		UserId:      comment.UserId,
		UserName:    comment.UserName,
	}
}

type RatingFromComment struct {
	Rating        int
	ResultRating  int
	ResultComment int
}

func (rating *RatingFromComment) ToRatingResponse() *model.GetRatingResponse {
	return &model.GetRatingResponse{
		Rating:        rating.Rating,
		ResultRating:  rating.ResultRating,
		ResultComment: rating.ResultComment,
	}
}
