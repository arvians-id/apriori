package resolver

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/arvians-id/apriori/cmd/library/auth"
	"github.com/arvians-id/apriori/cmd/library/aws"
	"github.com/arvians-id/apriori/cmd/library/messaging"
	"github.com/arvians-id/apriori/internal/service"
)

type Resolver struct {
	AprioriService       service.AprioriService
	CategoryService      service.CategoryService
	CommentService       service.CommentService
	EmailService         service.EmailService
	NotificationService  service.NotificationService
	PasswordResetService service.PasswordResetService
	PaymentService       service.PaymentService
	ProductService       service.ProductService
	TransactionService   service.TransactionService
	UserOrderService     service.UserOrderService
	UserService          service.UserService
	StorageS3            *aws.StorageS3
	Jwt                  *auth.JsonWebToken
	Producer             *messaging.Producer
}
