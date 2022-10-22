package resolver

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/arvians-id/apriori/cmd/library/cache"
	"github.com/arvians-id/apriori/internal/service"
)

type Resolver struct {
	AprioriService       service.AprioriService
	Redis                *cache.Redis
	CategoryService      service.CategoryService
	CommentService       service.CommentService
	EmailService         service.EmailService
	JwtService           service.JwtService
	NotificationService  service.NotificationService
	PasswordResetService service.PasswordResetService
	PaymentService       service.PaymentService
	ProductService       service.ProductService
	StorageService       service.StorageService
	TransactionService   service.TransactionService
	UserOrderService     service.UserOrderService
	UserService          service.UserService
}
