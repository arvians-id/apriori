package repository

import (
	"apriori/entity"
	"context"
	"database/sql"
)

type ProductRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx, search string, category string) ([]*entity.Product, error)
	FindAllByAdmin(ctx context.Context, tx *sql.Tx) ([]*entity.Product, error)
	FindAllBySimilarCategory(ctx context.Context, tx *sql.Tx, category string) ([]*entity.Product, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (*entity.Product, error)
	FindByCode(ctx context.Context, tx *sql.Tx, code string) (*entity.Product, error)
	FindByName(ctx context.Context, tx *sql.Tx, name string) (*entity.Product, error)
	Create(ctx context.Context, tx *sql.Tx, product *entity.Product) (*entity.Product, error)
	Update(ctx context.Context, tx *sql.Tx, product *entity.Product) (*entity.Product, error)
	Delete(ctx context.Context, tx *sql.Tx, code string) error
}

type AprioriRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]*entity.Apriori, error)
	FindAllByActive(ctx context.Context, tx *sql.Tx) ([]*entity.Apriori, error)
	FindAllByCode(ctx context.Context, tx *sql.Tx, code string) ([]*entity.Apriori, error)
	FindByCodeAndId(ctx context.Context, tx *sql.Tx, code string, id int) (*entity.Apriori, error)
	Create(ctx context.Context, tx *sql.Tx, apriories []*entity.Apriori) error
	Update(ctx context.Context, tx *sql.Tx, apriori *entity.Apriori) (*entity.Apriori, error)
	Delete(ctx context.Context, tx *sql.Tx, code string) error
	UpdateAllStatus(ctx context.Context, tx *sql.Tx, status bool) error
	UpdateStatusByCode(ctx context.Context, tx *sql.Tx, code string, status bool) error
}

type PasswordResetRepository interface {
	FindByEmail(ctx context.Context, tx *sql.Tx, email string) (*entity.PasswordReset, error)
	FindByEmailAndToken(ctx context.Context, tx *sql.Tx, passwordReset *entity.PasswordReset) (*entity.PasswordReset, error)
	Create(ctx context.Context, tx *sql.Tx, passwordReset *entity.PasswordReset) (*entity.PasswordReset, error)
	Update(ctx context.Context, tx *sql.Tx, passwordReset *entity.PasswordReset) (*entity.PasswordReset, error)
	Delete(ctx context.Context, tx *sql.Tx, email string) error
	Truncate(ctx context.Context, tx *sql.Tx) error
}

type PaymentRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]*entity.PaymentRelation, error)
	FindAllByUserId(ctx context.Context, tx *sql.Tx, userId int) ([]*entity.Payment, error)
	FindByOrderId(ctx context.Context, tx *sql.Tx, orderId string) (*entity.Payment, error)
	Create(ctx context.Context, tx *sql.Tx, payment *entity.Payment) (*entity.Payment, error)
	Update(ctx context.Context, tx *sql.Tx, payment *entity.Payment) error
	UpdateReceiptNumber(ctx context.Context, tx *sql.Tx, payment *entity.Payment) error
	Delete(ctx context.Context, tx *sql.Tx, orderId string) error
}

type TransactionRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]*entity.Transaction, error)
	FindAllItemSet(ctx context.Context, tx *sql.Tx, startDate string, endDate string) ([]*entity.Transaction, error)
	FindByNoTransaction(ctx context.Context, tx *sql.Tx, noTransaction string) (*entity.Transaction, error)
	CreateByCsv(ctx context.Context, tx *sql.Tx, transaction []*entity.Transaction) error
	Create(ctx context.Context, tx *sql.Tx, transaction *entity.Transaction) (*entity.Transaction, error)
	Update(ctx context.Context, tx *sql.Tx, transaction *entity.Transaction) (*entity.Transaction, error)
	Delete(ctx context.Context, tx *sql.Tx, noTransaction string) error
	Truncate(ctx context.Context, tx *sql.Tx) error
}

type UserOrderRepository interface {
	FindAllByPayloadId(ctx context.Context, tx *sql.Tx, payloadId string) ([]*entity.UserOrder, error)
	FindAllByUserId(ctx context.Context, tx *sql.Tx, userId int) ([]*entity.UserOrderRelationByUserId, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (*entity.UserOrder, error)
	Create(ctx context.Context, tx *sql.Tx, userOrder *entity.UserOrder) error
}

type UserRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]*entity.User, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (*entity.User, error)
	FindByEmail(ctx context.Context, tx *sql.Tx, email string) (*entity.User, error)
	Create(ctx context.Context, tx *sql.Tx, user *entity.User) (*entity.User, error)
	Update(ctx context.Context, tx *sql.Tx, user *entity.User) (*entity.User, error)
	UpdatePassword(ctx context.Context, tx *sql.Tx, user *entity.User) error
	Delete(ctx context.Context, tx *sql.Tx, id int) error
}

type CategoryRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]*entity.Category, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (*entity.Category, error)
	Create(ctx context.Context, tx *sql.Tx, category *entity.Category) (*entity.Category, error)
	Update(ctx context.Context, tx *sql.Tx, category *entity.Category) (*entity.Category, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
}

type CommentRepository interface {
	FindAllRatingByProductCode(ctx context.Context, tx *sql.Tx, productCode string) ([]*entity.RatingFromComment, error)
	FindAllByProductCode(ctx context.Context, tx *sql.Tx, productCode string, rating string, tags string) ([]*entity.Comment, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (*entity.Comment, error)
	FindByUserOrderId(ctx context.Context, tx *sql.Tx, userOrderId int) (*entity.Comment, error)
	Create(ctx context.Context, tx *sql.Tx, comment *entity.Comment) (*entity.Comment, error)
}

type NotificationRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]*entity.NotificationRelation, error)
	FindAllByUserId(ctx context.Context, tx *sql.Tx, userId int) ([]*entity.Notification, error)
	Create(ctx context.Context, tx *sql.Tx, notification *entity.Notification) (*entity.Notification, error)
	Mark(ctx context.Context, tx *sql.Tx, id int) error
	MarkAll(ctx context.Context, tx *sql.Tx, userId int) error
}
