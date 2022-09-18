package service

import (
	"context"
	"github.com/arvians-id/apriori/entity"
	request2 "github.com/arvians-id/apriori/http/request"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"mime/multipart"
	"sync"
	"time"
)

type AprioriService interface {
	FindAll(ctx context.Context) ([]*entity.Apriori, error)
	FindAllByActive(ctx context.Context) ([]*entity.Apriori, error)
	FindAllByCode(ctx context.Context, code string) ([]*entity.Apriori, error)
	FindByCodeAndId(ctx context.Context, code string, id int) (*entity.ProductRecommendation, error)
	Create(ctx context.Context, requests []*request2.CreateAprioriRequest) error
	Update(ctx context.Context, request *request2.UpdateAprioriRequest) (*entity.Apriori, error)
	UpdateStatus(ctx context.Context, code string) error
	Delete(ctx context.Context, code string) error
	Generate(ctx context.Context, request *request2.GenerateAprioriRequest) ([]*entity.GenerateApriori, error)
}

type CacheService interface {
	GetClient() (*redis.Client, error)
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value interface{}) error
	Del(ctx context.Context, key ...string) error
	FlushDB(ctx context.Context) error
}

type CategoryService interface {
	FindAll(ctx context.Context) ([]*entity.Category, error)
	FindById(ctx context.Context, id int) (*entity.Category, error)
	Create(ctx context.Context, request *request2.CreateCategoryRequest) (*entity.Category, error)
	Update(ctx context.Context, request *request2.UpdateCategoryRequest) (*entity.Category, error)
	Delete(ctx context.Context, id int) error
}

type CommentService interface {
	FindAllRatingByProductCode(ctx context.Context, productCode string) ([]*entity.RatingFromComment, error)
	FindAllByProductCode(ctx context.Context, productCode string, rating string, tags string) ([]*entity.Comment, error)
	FindById(ctx context.Context, id int) (*entity.Comment, error)
	FindByUserOrderId(ctx context.Context, userOrderId int) (*entity.Comment, error)
	Create(ctx context.Context, request *request2.CreateCommentRequest) (*entity.Comment, error)
}

type EmailService interface {
	SendEmailWithText(toEmail string, subject string, message *string) error
}

type JwtService interface {
	GenerateToken(IdUser int, expirationTime time.Time) (*TokenDetails, error)
	RefreshToken(refreshToken string) (*TokenDetails, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type PasswordResetService interface {
	CreateOrUpdateByEmail(ctx context.Context, email string) (*entity.PasswordReset, error)
	Verify(ctx context.Context, request *request2.UpdateResetPasswordUserRequest) error
}

type PaymentService interface {
	GetClient()
	FindAll(ctx context.Context) ([]*entity.Payment, error)
	FindAllByUserId(ctx context.Context, userId int) ([]*entity.Payment, error)
	FindByOrderId(ctx context.Context, orderId string) (*entity.Payment, error)
	CreateOrUpdate(ctx context.Context, request map[string]interface{}) error
	UpdateReceiptNumber(ctx context.Context, request *request2.AddReceiptNumberRequest) (*entity.Payment, error)
	Delete(ctx context.Context, orderId string) error
	GetToken(ctx context.Context, request *request2.GetPaymentTokenRequest) (map[string]interface{}, error)
}

type ProductService interface {
	FindAllByAdmin(ctx context.Context) ([]*entity.Product, error)
	FindAll(ctx context.Context, search string, category string) ([]*entity.Product, error)
	FindAllBySimilarCategory(ctx context.Context, code string) ([]*entity.Product, error)
	FindAllRecommendation(ctx context.Context, code string) ([]*entity.ProductRecommendation, error)
	FindByCode(ctx context.Context, code string) (*entity.Product, error)
	Create(ctx context.Context, request *request2.CreateProductRequest) (*entity.Product, error)
	Update(ctx context.Context, request *request2.UpdateProductRequest) (*entity.Product, error)
	Delete(ctx context.Context, code string) error
}

type StorageService interface {
	UploadFile(c *gin.Context, image *multipart.FileHeader) (chan string, error)
	UploadFileS3(file multipart.File, header *multipart.FileHeader) (string, error)
	WaitUploadFileS3(file multipart.File, header *multipart.FileHeader, wg *sync.WaitGroup) (string, error)
	//DeleteFileS3(fileName string) error
	//WaitDeleteFileS3(fileName string, wg *sync.WaitGroup) error
}

type TransactionService interface {
	FindAll(ctx context.Context) ([]*entity.Transaction, error)
	FindByNoTransaction(ctx context.Context, noTransaction string) (*entity.Transaction, error)
	Create(ctx context.Context, request *request2.CreateTransactionRequest) (*entity.Transaction, error)
	CreateByCsv(ctx context.Context, data [][]string) error
	Update(ctx context.Context, request *request2.UpdateTransactionRequest) (*entity.Transaction, error)
	Delete(ctx context.Context, noTransaction string) error
	Truncate(ctx context.Context) error
}

type UserOrderService interface {
	FindAllByPayloadId(ctx context.Context, payloadId int) ([]*entity.UserOrder, error)
	FindAllByUserId(ctx context.Context, userId int) ([]*entity.UserOrder, error)
	FindById(ctx context.Context, id int) (*entity.UserOrder, error)
}

type UserService interface {
	FindAll(ctx context.Context) ([]*entity.User, error)
	FindById(ctx context.Context, id int) (*entity.User, error)
	FindByEmail(ctx context.Context, request *request2.GetUserCredentialRequest) (*entity.User, error)
	Create(ctx context.Context, request *request2.CreateUserRequest) (*entity.User, error)
	Update(ctx context.Context, request *request2.UpdateUserRequest) (*entity.User, error)
	Delete(ctx context.Context, id int) error
}

type NotificationService interface {
	FindAll(ctx context.Context) ([]*entity.Notification, error)
	FindAllByUserId(ctx context.Context, userId int) ([]*entity.Notification, error)
	Create(ctx context.Context, request *request2.CreateNotificationRequest) *NotificationServiceImpl
	MarkAll(ctx context.Context, userId int) error
	Mark(ctx context.Context, id int) error
	WithSendMail() error
}
