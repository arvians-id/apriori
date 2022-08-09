package repository

import (
	"apriori/entity"
	"context"
	"database/sql"
)

type ProductRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx, search string) ([]entity.Product, error)
	FindById(ctx context.Context, tx *sql.Tx, productId uint64) (entity.Product, error)
	FindByCode(ctx context.Context, tx *sql.Tx, code string) (entity.Product, error)
	FindByName(ctx context.Context, tx *sql.Tx, name string) (entity.Product, error)
	Create(ctx context.Context, tx *sql.Tx, product entity.Product) (entity.Product, error)
	Update(ctx context.Context, tx *sql.Tx, product entity.Product) (entity.Product, error)
	Delete(ctx context.Context, tx *sql.Tx, code string) error
}

type AprioriRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]entity.Apriori, error)
	FindByActive(ctx context.Context, tx *sql.Tx) ([]entity.Apriori, error)
	FindByCode(ctx context.Context, tx *sql.Tx, code string) ([]entity.Apriori, error)
	FindByCodeAndId(ctx context.Context, tx *sql.Tx, code string, id int) (entity.Apriori, error)
	UpdateApriori(ctx context.Context, tx *sql.Tx, apriori entity.Apriori) (entity.Apriori, error)
	ChangeAllStatus(ctx context.Context, tx *sql.Tx, status int) error
	ChangeStatusByCode(ctx context.Context, tx *sql.Tx, code string, status int) error
	Create(ctx context.Context, tx *sql.Tx, apriories []entity.Apriori) error
	Delete(ctx context.Context, tx *sql.Tx, code string) error
}

type AuthRepository interface {
	VerifyCredential(ctx context.Context, tx *sql.Tx, email string, password string) (entity.User, error)
}

type PasswordResetRepository interface {
	FindByEmail(ctx context.Context, tx *sql.Tx, email string) (entity.PasswordReset, error)
	FindByEmailAndToken(ctx context.Context, tx *sql.Tx, reset entity.PasswordReset) (entity.PasswordReset, error)
	Create(ctx context.Context, tx *sql.Tx, reset entity.PasswordReset) (entity.PasswordReset, error)
	Update(ctx context.Context, tx *sql.Tx, reset entity.PasswordReset) (entity.PasswordReset, error)
	Delete(ctx context.Context, tx *sql.Tx, email string) error
}

type PaymentRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]entity.PaymentRelation, error)
	FindAllByUserId(ctx context.Context, tx *sql.Tx, userId int) ([]entity.PaymentNullable, error)
	Create(ctx context.Context, tx *sql.Tx, payment entity.Payment) (entity.Payment, error)
	FindByOrderId(ctx context.Context, tx *sql.Tx, orderId string) (entity.PaymentNullable, error)
	Update(ctx context.Context, tx *sql.Tx, payment entity.Payment) error
	Delete(ctx context.Context, tx *sql.Tx, orderId string) error
}

type TransactionRepository interface {
	FindItemSet(ctx context.Context, tx *sql.Tx, startDate string, endDate string) ([]entity.Transaction, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]entity.Transaction, error)
	FindByTransaction(ctx context.Context, tx *sql.Tx, noTransaction string) (entity.Transaction, error)
	Create(ctx context.Context, tx *sql.Tx, transaction entity.Transaction) (entity.Transaction, error)
	CreateFromCsv(ctx context.Context, tx *sql.Tx, transaction []entity.Transaction) error
	Update(ctx context.Context, tx *sql.Tx, transaction entity.Transaction) (entity.Transaction, error)
	Delete(ctx context.Context, tx *sql.Tx, noTransaction string) error
	Truncate(ctx context.Context, tx *sql.Tx) error
}

type UserOrderRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx, payloadId string) ([]entity.UserOrder, error)
	Create(ctx context.Context, tx *sql.Tx, userOrder entity.UserOrder) error
}

type UserRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]entity.User, error)
	FindById(ctx context.Context, tx *sql.Tx, userId uint64) (entity.User, error)
	FindByEmail(ctx context.Context, tx *sql.Tx, email string) (entity.User, error)
	Create(ctx context.Context, tx *sql.Tx, user entity.User) (entity.User, error)
	Update(ctx context.Context, tx *sql.Tx, user entity.User) (entity.User, error)
	UpdatePassword(ctx context.Context, tx *sql.Tx, user entity.User) error
	Delete(ctx context.Context, tx *sql.Tx, userId uint64) error
}

type CategoryRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]entity.Category, error)
	FindById(ctx context.Context, tx *sql.Tx, categoryId int) (entity.Category, error)
	Create(ctx context.Context, tx *sql.Tx, category entity.Category) (entity.Category, error)
	Update(ctx context.Context, tx *sql.Tx, category entity.Category) (entity.Category, error)
	Delete(ctx context.Context, tx *sql.Tx, categoryId int) error
}
