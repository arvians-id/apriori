package repository

import (
	"context"
	"database/sql"
	"github.com/arvians-id/apriori/model"
)

type ProductRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx, search string, category string) ([]*model.Product, error)
	FindAllByAdmin(ctx context.Context, tx *sql.Tx) ([]*model.Product, error)
	FindAllBySimilarCategory(ctx context.Context, tx *sql.Tx, category string) ([]*model.Product, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (*model.Product, error)
	FindByCode(ctx context.Context, tx *sql.Tx, code string) (*model.Product, error)
	FindByName(ctx context.Context, tx *sql.Tx, name string) (*model.Product, error)
	Create(ctx context.Context, tx *sql.Tx, product *model.Product) (*model.Product, error)
	Update(ctx context.Context, tx *sql.Tx, product *model.Product) (*model.Product, error)
	Delete(ctx context.Context, tx *sql.Tx, code string) error
}

type AprioriRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]*model.Apriori, error)
	FindAllByActive(ctx context.Context, tx *sql.Tx) ([]*model.Apriori, error)
	FindAllByCode(ctx context.Context, tx *sql.Tx, code string) ([]*model.Apriori, error)
	FindByCodeAndId(ctx context.Context, tx *sql.Tx, code string, id int) (*model.Apriori, error)
	Create(ctx context.Context, tx *sql.Tx, apriories []*model.Apriori) error
	Update(ctx context.Context, tx *sql.Tx, apriori *model.Apriori) (*model.Apriori, error)
	Delete(ctx context.Context, tx *sql.Tx, code string) error
	UpdateAllStatus(ctx context.Context, tx *sql.Tx, status bool) error
	UpdateStatusByCode(ctx context.Context, tx *sql.Tx, code string, status bool) error
}

type PasswordResetRepository interface {
	FindByEmail(ctx context.Context, tx *sql.Tx, email string) (*model.PasswordReset, error)
	FindByEmailAndToken(ctx context.Context, tx *sql.Tx, passwordReset *model.PasswordReset) (*model.PasswordReset, error)
	Create(ctx context.Context, tx *sql.Tx, passwordReset *model.PasswordReset) (*model.PasswordReset, error)
	Update(ctx context.Context, tx *sql.Tx, passwordReset *model.PasswordReset) (*model.PasswordReset, error)
	Delete(ctx context.Context, tx *sql.Tx, email string) error
	Truncate(ctx context.Context, tx *sql.Tx) error
}

type PaymentRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]*model.Payment, error)
	FindAllByUserId(ctx context.Context, tx *sql.Tx, userId int) ([]*model.Payment, error)
	FindByOrderId(ctx context.Context, tx *sql.Tx, orderId string) (*model.Payment, error)
	Create(ctx context.Context, tx *sql.Tx, payment *model.Payment) (*model.Payment, error)
	Update(ctx context.Context, tx *sql.Tx, payment *model.Payment) error
	UpdateReceiptNumber(ctx context.Context, tx *sql.Tx, payment *model.Payment) error
	Delete(ctx context.Context, tx *sql.Tx, orderId *string) error
}

type TransactionRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]*model.Transaction, error)
	FindAllItemSet(ctx context.Context, tx *sql.Tx, startDate string, endDate string) ([]*model.Transaction, error)
	FindByNoTransaction(ctx context.Context, tx *sql.Tx, noTransaction string) (*model.Transaction, error)
	CreateByCsv(ctx context.Context, tx *sql.Tx, transaction []*model.Transaction) error
	Create(ctx context.Context, tx *sql.Tx, transaction *model.Transaction) (*model.Transaction, error)
	Update(ctx context.Context, tx *sql.Tx, transaction *model.Transaction) (*model.Transaction, error)
	Delete(ctx context.Context, tx *sql.Tx, noTransaction string) error
	Truncate(ctx context.Context, tx *sql.Tx) error
}

type UserOrderRepository interface {
	FindAllByPayloadId(ctx context.Context, tx *sql.Tx, payloadId string) ([]*model.UserOrder, error)
	FindAllByUserId(ctx context.Context, tx *sql.Tx, userId int) ([]*model.UserOrder, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (*model.UserOrder, error)
	Create(ctx context.Context, tx *sql.Tx, userOrder *model.UserOrder) (*model.UserOrder, error)
}

type UserRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]*model.User, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (*model.User, error)
	FindByEmail(ctx context.Context, tx *sql.Tx, email string) (*model.User, error)
	Create(ctx context.Context, tx *sql.Tx, user *model.User) (*model.User, error)
	Update(ctx context.Context, tx *sql.Tx, user *model.User) (*model.User, error)
	UpdatePassword(ctx context.Context, tx *sql.Tx, user *model.User) error
	Delete(ctx context.Context, tx *sql.Tx, id int) error
}

type CategoryRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]*model.Category, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (*model.Category, error)
	Create(ctx context.Context, tx *sql.Tx, category *model.Category) (*model.Category, error)
	Update(ctx context.Context, tx *sql.Tx, category *model.Category) (*model.Category, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
}

type CommentRepository interface {
	FindAllRatingByProductCode(ctx context.Context, tx *sql.Tx, productCode string) ([]*model.RatingFromComment, error)
	FindAllByProductCode(ctx context.Context, tx *sql.Tx, productCode string, rating string, tags string) ([]*model.Comment, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (*model.Comment, error)
	FindByUserOrderId(ctx context.Context, tx *sql.Tx, userOrderId int) (*model.Comment, error)
	Create(ctx context.Context, tx *sql.Tx, comment *model.Comment) (*model.Comment, error)
}

type NotificationRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]*model.Notification, error)
	FindAllByUserId(ctx context.Context, tx *sql.Tx, userId int) ([]*model.Notification, error)
	Create(ctx context.Context, tx *sql.Tx, notification *model.Notification) (*model.Notification, error)
	Mark(ctx context.Context, tx *sql.Tx, id int) error
	MarkAll(ctx context.Context, tx *sql.Tx, userId int) error
}
