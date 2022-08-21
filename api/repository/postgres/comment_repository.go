package postgres

import (
	"apriori/entity"
	"apriori/repository"
	"context"
	"database/sql"
	"log"
)

type commentRepository struct {
}

func NewCommentRepository() repository.CommentRepository {
	return &commentRepository{}
}

func (repository *commentRepository) FindAllRatingByProductCode(ctx context.Context, tx *sql.Tx, productCode string) ([]*entity.RatingFromComment, error) {
	query := `SELECT rating, rating * COUNT(rating) as result_rating, SUM(CASE WHEN description != '' THEN 1 ELSE 0 END) as result_comment 
			  FROM comments 
			  WHERE product_code = $1 
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

	var ratings []*entity.RatingFromComment
	for rows.Next() {
		var rating entity.RatingFromComment
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

func (repository *commentRepository) FindAllByProductCode(ctx context.Context, tx *sql.Tx, productCode string, rating string, tags string) ([]*entity.Comment, error) {
	query := `SELECT c.*,u.id_user,u.name 
			  FROM comments c 
				LEFT JOIN user_orders uo ON uo.id_order = c.user_order_id 
				LEFT JOIN payloads p ON p.id_payload = uo.payload_id 
				LEFT JOIN users u ON u.id_user = p.user_id 
			  WHERE c.product_code = $1 AND CAST(c.rating as TEXT) LIKE $2 AND c.tag SIMILAR TO $3
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

	var comments []*entity.Comment
	for rows.Next() {
		var comment entity.Comment
		err := rows.Scan(
			&comment.IdComment,
			&comment.UserOrderId,
			&comment.ProductCode,
			&comment.Description,
			&comment.Tag,
			&comment.Rating,
			&comment.CreatedAt,
			&comment.UserId,
			&comment.UserName,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, &comment)
	}

	return comments, nil
}

func (repository *commentRepository) FindById(ctx context.Context, tx *sql.Tx, id int) (*entity.Comment, error) {
	query := `SELECT c.*,u.id_user,u.name 
			  FROM comments c 
				LEFT JOIN user_orders uo ON uo.id_order = c.user_order_id 
				LEFT JOIN payloads p ON p.id_payload = uo.payload_id 
				LEFT JOIN users u ON u.id_user = p.user_id 
			  WHERE c.id_comment = $1`
	row := tx.QueryRowContext(ctx, query, id)

	var comment entity.Comment
	err := row.Scan(
		&comment.IdComment,
		&comment.UserOrderId,
		&comment.ProductCode,
		&comment.Description,
		&comment.Tag,
		&comment.Rating,
		&comment.CreatedAt,
		&comment.UserId,
		&comment.UserName,
	)
	if err != nil {
		return &entity.Comment{}, err
	}

	return &comment, nil
}

func (repository *commentRepository) FindByUserOrderId(ctx context.Context, tx *sql.Tx, userOrderId int) (*entity.Comment, error) {
	query := `SELECT c.*,u.id_user,u.name 
			  FROM comments c 
				LEFT JOIN user_orders uo ON uo.id_order = c.user_order_id 
				LEFT JOIN payloads p ON p.id_payload = uo.payload_id 
				LEFT JOIN users u ON u.id_user = p.user_id 
			  WHERE c.user_order_id = $1`
	row := tx.QueryRowContext(ctx, query, userOrderId)

	var comment entity.Comment
	err := row.Scan(
		&comment.IdComment,
		&comment.UserOrderId,
		&comment.ProductCode,
		&comment.Description,
		&comment.Tag,
		&comment.Rating,
		&comment.CreatedAt,
		&comment.UserId,
		&comment.UserName,
	)
	if err != nil {
		return &entity.Comment{}, err
	}

	return &comment, nil

}

func (repository *commentRepository) Create(ctx context.Context, tx *sql.Tx, comment *entity.Comment) (*entity.Comment, error) {
	id := 0
	query := `INSERT INTO comments (user_order_id, product_code, description, tag, rating, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id_comment`
	queryContext := tx.QueryRowContext(
		ctx,
		query,
		comment.UserOrderId,
		comment.ProductCode,
		comment.Description,
		comment.Tag,
		comment.Rating,
		comment.CreatedAt,
	)
	err := queryContext.Scan(&id)
	if err != nil {
		return &entity.Comment{}, err
	}

	comment.IdComment = id

	return comment, nil
}
