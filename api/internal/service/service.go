package service

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	request2 "github.com/arvians-id/apriori/internal/http/presenter/request"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"mime/multipart"
	"sync"
	"time"
)

type AprioriService interface {
	FindAll(ctx context.Context) ([]*model.Apriori, error)
	FindAllByActive(ctx context.Context) ([]*model.Apriori, error)
	FindAllByCode(ctx context.Context, code string) ([]*model.Apriori, error)
	FindByCodeAndId(ctx context.Context, code string, id int) (*model.ProductRecommendation, error)
	Create(ctx context.Context, requests []*request2.CreateAprioriRequest) error
	Update(ctx context.Context, request *request2.UpdateAprioriRequest) (*model.Apriori, error)
	UpdateStatus(ctx context.Context, code string) error
	Delete(ctx context.Context, code string) error
	Generate(ctx context.Context, request *request2.GenerateAprioriRequest) ([]*model.GenerateApriori, error)
}

type CategoryService interface {
	FindAll(ctx context.Context) ([]*model.Category, error)
	FindById(ctx context.Context, id int) (*model.Category, error)
	Create(ctx context.Context, request *request2.CreateCategoryRequest) (*model.Category, error)
	Update(ctx context.Context, request *request2.UpdateCategoryRequest) (*model.Category, error)
	Delete(ctx context.Context, id int) error
}

type CommentService interface {
	FindAllRatingByProductCode(ctx context.Context, productCode string) ([]*model.RatingFromComment, error)
	FindAllByProductCode(ctx context.Context, productCode string, rating string, tags string) ([]*model.Comment, error)
	FindById(ctx context.Context, id int) (*model.Comment, error)
	FindByUserOrderId(ctx context.Context, userOrderId int) (*model.Comment, error)
	Create(ctx context.Context, request *request2.CreateCommentRequest) (*model.Comment, error)
}

type EmailService interface {
	SendEmailWithText(toEmail string, subject string, message *string) error
}

type JwtService interface {
	GenerateToken(IdUser int, RoleUser int, expirationTime time.Time) (*TokenDetails, error)
	RefreshToken(refreshToken string) (*TokenDetails, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type PasswordResetService interface {
	CreateOrUpdateByEmail(ctx context.Context, email string) (*model.PasswordReset, error)
	Verify(ctx context.Context, request *request2.UpdateResetPasswordUserRequest) error
}

type PaymentService interface {
	GetClient()
	FindAll(ctx context.Context) ([]*model.Payment, error)
	FindAllByUserId(ctx context.Context, userId int) ([]*model.Payment, error)
	FindByOrderId(ctx context.Context, orderId string) (*model.Payment, error)
	CreateOrUpdate(ctx context.Context, request map[string]interface{}) error
	UpdateReceiptNumber(ctx context.Context, request *request2.AddReceiptNumberRequest) (*model.Payment, error)
	Delete(ctx context.Context, orderId string) error
	GetToken(ctx context.Context, request *request2.GetPaymentTokenRequest) (map[string]interface{}, error)
}

type ProductService interface {
	FindAllByAdmin(ctx context.Context) ([]*model.Product, error)
	FindAll(ctx context.Context, search string, category string) ([]*model.Product, error)
	FindAllBySimilarCategory(ctx context.Context, code string) ([]*model.Product, error)
	FindAllRecommendation(ctx context.Context, code string) ([]*model.ProductRecommendation, error)
	FindByCode(ctx context.Context, code string) (*model.Product, error)
	Create(ctx context.Context, request *request2.CreateProductRequest) (*model.Product, error)
	Update(ctx context.Context, request *request2.UpdateProductRequest) (*model.Product, error)
	Delete(ctx context.Context, code string) error
}

type StorageService interface {
	UploadFile(c *gin.Context, image *multipart.FileHeader) (chan string, error)
	UploadFileS3(file multipart.File, header *multipart.FileHeader) (string, error)
	UploadFileS3GraphQL(file graphql.Upload, initFileName string) (string, error)
	WaitUploadFileS3(file multipart.File, header *multipart.FileHeader, wg *sync.WaitGroup) (string, error)
	//DeleteFileS3(fileName string) error
	//WaitDeleteFileS3(fileName string, wg *sync.WaitGroup) error
}

type TransactionService interface {
	FindAll(ctx context.Context) ([]*model.Transaction, error)
	FindByNoTransaction(ctx context.Context, noTransaction string) (*model.Transaction, error)
	Create(ctx context.Context, request *request2.CreateTransactionRequest) (*model.Transaction, error)
	CreateByCsv(ctx context.Context, data [][]string) error
	Update(ctx context.Context, request *request2.UpdateTransactionRequest) (*model.Transaction, error)
	Delete(ctx context.Context, noTransaction string) error
	Truncate(ctx context.Context) error
}

type UserOrderService interface {
	FindAllByPayloadId(ctx context.Context, payloadId int) ([]*model.UserOrder, error)
	FindAllByUserId(ctx context.Context, userId int) ([]*model.UserOrder, error)
	FindById(ctx context.Context, id int) (*model.UserOrder, error)
}

type UserService interface {
	FindAll(ctx context.Context) ([]*model.User, error)
	FindById(ctx context.Context, id int) (*model.User, error)
	FindByEmail(ctx context.Context, request *request2.GetUserCredentialRequest) (*model.User, error)
	Create(ctx context.Context, request *request2.CreateUserRequest) (*model.User, error)
	Update(ctx context.Context, request *request2.UpdateUserRequest) (*model.User, error)
	Delete(ctx context.Context, id int) error
}

type NotificationService interface {
	FindAll(ctx context.Context) ([]*model.Notification, error)
	FindAllByUserId(ctx context.Context, userId int) ([]*model.Notification, error)
	Create(ctx context.Context, request *request2.CreateNotificationRequest) *NotificationServiceImpl
	MarkAll(ctx context.Context, userId int) error
	Mark(ctx context.Context, id int) error
	WithSendMail() error
}
