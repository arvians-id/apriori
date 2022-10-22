package mysql

import (
	"context"
	"database/sql"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/arvians-id/apriori/internal/repository"
	"log"
)

type CommentRepositoryImpl struct {
}

func NewCommentRepository() repository.CommentRepository {
	return &CommentRepositoryImpl{}
}

func (repository *CommentRepositoryImpl) FindAllRatingByProductCode(ctx context.Context, tx *sql.Tx, productCode string) ([]*model.RatingFromComment, error) {
	query := `SELECT rating, rating * COUNT(rating) as result_rating, SUM(CASE WHEN description != '' THEN 1 ELSE 0 END) as result_comment 
			  FROM comments 
			  WHERE product_code = ? 
			  GROUP BY rating 
			  ORDER BY rating DESC`
	rows, err := tx.QueryContext(ctx, query, productCode)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(rows)

	var ratings []*model.RatingFromComment
	for rows.Next() {
		var rating model.RatingFromComment
		err := rows.Scan(
			&rating.Rating,
			&rating.ResultRating,
			&rating.ResultComment,
		)
		if err != nil {
			return nil, err
		}

		ratings = append(ratings, &rating)
	}

	return ratings, nil
}

func (repository *CommentRepositoryImpl) FindAllByProductCode(ctx context.Context, tx *sql.Tx, productCode string, rating string, tags string) ([]*model.Comment, error) {
	query := `SELECT c.*,u.id_user,u.name 
			  FROM comments c 
				LEFT JOIN user_orders uo ON uo.id_order = c.user_order_id 
				LEFT JOIN payloads p ON p.id_payload = uo.payload_id 
				LEFT JOIN users u ON u.id_user = p.user_id 
			  WHERE c.product_code = ? AND CAST(c.rating as TEXT) LIKE ? AND c.tag SIMILAR TO ?
			  ORDER BY c.created_at DESC`
	rows, err := tx.QueryContext(ctx, query, productCode, "%"+rating+"%", "%("+tags+")%")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(rows)

	var comments []*model.Comment
	for rows.Next() {
		comment := model.Comment{
			UserOrder: &model.UserOrder{
				Payment: &model.Payment{
					User: &model.User{},
				},
			},
		}
		err := rows.Scan(
			&comment.IdComment,
			&comment.UserOrderId,
			&comment.ProductCode,
			&comment.Description,
			&comment.Tag,
			&comment.Rating,
			&comment.CreatedAt,
			&comment.UserOrder.Payment.User.IdUser,
			&comment.UserOrder.Payment.User.Name,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, &comment)
	}

	return comments, nil
}

func (repository *CommentRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (*model.Comment, error) {
	query := `SELECT c.*,u.id_user,u.name 
			  FROM comments c 
				LEFT JOIN user_orders uo ON uo.id_order = c.user_order_id 
				LEFT JOIN payloads p ON p.id_payload = uo.payload_id 
				LEFT JOIN users u ON u.id_user = p.user_id 
			  WHERE c.id_comment = ?`
	row := tx.QueryRowContext(ctx, query, id)

	comment := model.Comment{
		UserOrder: &model.UserOrder{
			Payment: &model.Payment{
				User: &model.User{},
			},
		},
	}
	err := row.Scan(
		&comment.IdComment,
		&comment.UserOrderId,
		&comment.ProductCode,
		&comment.Description,
		&comment.Tag,
		&comment.Rating,
		&comment.CreatedAt,
		&comment.UserOrder.Payment.User.IdUser,
		&comment.UserOrder.Payment.User.Name,
	)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (repository *CommentRepositoryImpl) FindByUserOrderId(ctx context.Context, tx *sql.Tx, userOrderId int) (*model.Comment, error) {
	query := `SELECT c.*,u.id_user,u.name 
			  FROM comments c 
				LEFT JOIN user_orders uo ON uo.id_order = c.user_order_id 
				LEFT JOIN payloads p ON p.id_payload = uo.payload_id 
				LEFT JOIN users u ON u.id_user = p.user_id 
			  WHERE c.user_order_id = ?`
	row := tx.QueryRowContext(ctx, query, userOrderId)

	comment := model.Comment{
		UserOrder: &model.UserOrder{
			Payment: &model.Payment{
				User: &model.User{},
			},
		},
	}
	err := row.Scan(
		&comment.IdComment,
		&comment.UserOrderId,
		&comment.ProductCode,
		&comment.Description,
		&comment.Tag,
		&comment.Rating,
		&comment.CreatedAt,
		&comment.UserOrder.Payment.User.IdUser,
		&comment.UserOrder.Payment.User.Name,
	)
	if err != nil {
		return nil, err
	}

	return &comment, nil

}

func (repository *CommentRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, comment *model.Comment) (*model.Comment, error) {
	query := `INSERT INTO comments (user_order_id, product_code, description, tag, rating, created_at) 
			  VALUES (?, ?, ?, ?, ?, ?)`
	row, err := tx.ExecContext(
		ctx,
		query,
		comment.UserOrderId,
		comment.ProductCode,
		comment.Description,
		comment.Tag,
		comment.Rating,
		comment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return nil, err
	}

	comment.IdComment = int(id)

	return comment, nil
}
