package service

import (
	"apriori/model"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"mime/multipart"
	"sync"
	"time"
)

type AprioriService interface {
	FindAll(ctx context.Context) ([]*model.GetAprioriResponse, error)
	FindAllByActive(ctx context.Context) ([]*model.GetAprioriResponse, error)
	FindAllByCode(ctx context.Context, code string) ([]*model.GetAprioriResponse, error)
	FindByCodeAndId(ctx context.Context, code string, id int) (*model.GetProductRecommendationResponse, error)
	Create(ctx context.Context, requests []*model.CreateAprioriRequest) error
	Update(ctx context.Context, request *model.UpdateAprioriRequest) (*model.GetAprioriResponse, error)
	UpdateStatus(ctx context.Context, code string) error
	Delete(ctx context.Context, code string) error
	Generate(ctx context.Context, request *model.GenerateAprioriRequest) ([]*model.GetGenerateAprioriResponse, error)
}

type CacheService interface {
	GetClient() (*redis.Client, error)
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}) error
	Del(ctx context.Context, key ...string) error
	Subscribe(ctx context.Context) (string, error)
	Publish(ctx context.Context) error
	FlushDB(ctx context.Context) error
}

type CategoryService interface {
	FindAll(ctx context.Context) ([]*model.GetCategoryResponse, error)
	FindById(ctx context.Context, id int) (*model.GetCategoryResponse, error)
	Create(ctx context.Context, request *model.CreateCategoryRequest) (*model.GetCategoryResponse, error)
	Update(ctx context.Context, request *model.UpdateCategoryRequest) (*model.GetCategoryResponse, error)
	Delete(ctx context.Context, id int) error
}

type CommentService interface {
	FindAllRatingByProductCode(ctx context.Context, productCode string) ([]*model.GetRatingResponse, error)
	FindAllByProductCode(ctx context.Context, productCode string, rating string, tags string) ([]*model.GetCommentResponse, error)
	FindById(ctx context.Context, id int) (*model.GetCommentResponse, error)
	FindByUserOrderId(ctx context.Context, userOrderId int) (*model.GetCommentResponse, error)
	Create(ctx context.Context, request *model.CreateCommentRequest) (*model.GetCommentResponse, error)
}

type EmailService interface {
	SendEmailWithText(toEmail string, message string) error
}

type JwtService interface {
	GenerateToken(IdUser int, expirationTime time.Time) (*TokenDetails, error)
	RefreshToken(refreshToken string) (*TokenDetails, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type PasswordResetService interface {
	CreateOrUpdateByEmail(ctx context.Context, email string) (*model.GetPasswordResetResponse, error)
	Verify(ctx context.Context, request *model.UpdateResetPasswordUserRequest) error
}

type PaymentService interface {
	GetClient()
	FindAll(ctx context.Context) ([]*model.GetPaymentRelationResponse, error)
	FindAllByUserId(ctx context.Context, userId int) ([]*model.GetPaymentResponse, error)
	FindByOrderId(ctx context.Context, orderId string) (*model.GetPaymentResponse, error)
	CreateOrUpdate(ctx context.Context, request map[string]interface{}) error
	UpdateReceiptNumber(ctx context.Context, request *model.AddReceiptNumberRequest) error
	Delete(ctx context.Context, orderId string) error
	GetToken(ctx context.Context, amount int64, userId int, customerName string, items []string, rajaShipping *model.GetRajaOngkirResponse) (map[string]interface{}, error)
}

type ProductService interface {
	FindAllByAdmin(ctx context.Context) ([]*model.GetProductResponse, error)
	FindAll(ctx context.Context, search string, category string) ([]*model.GetProductResponse, error)
	FindAllBySimilarCategory(ctx context.Context, code string) ([]*model.GetProductResponse, error)
	FindAllRecommendation(ctx context.Context, code string) ([]*model.GetProductRecommendationResponse, error)
	FindByCode(ctx context.Context, code string) (*model.GetProductResponse, error)
	Create(ctx context.Context, request *model.CreateProductRequest) (*model.GetProductResponse, error)
	Update(ctx context.Context, request *model.UpdateProductRequest) (*model.GetProductResponse, error)
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
	FindAll(ctx context.Context) ([]*model.GetTransactionResponse, error)
	FindByNoTransaction(ctx context.Context, noTransaction string) (*model.GetTransactionResponse, error)
	Create(ctx context.Context, request *model.CreateTransactionRequest) (*model.GetTransactionResponse, error)
	CreateByCsv(ctx context.Context, data [][]string) error
	Update(ctx context.Context, request *model.UpdateTransactionRequest) (*model.GetTransactionResponse, error)
	Delete(ctx context.Context, noTransaction string) error
	Truncate(ctx context.Context) error
}

type UserOrderService interface {
	FindAllByPayloadId(ctx context.Context, payloadId int) ([]*model.GetUserOrderResponse, error)
	FindAllByUserId(ctx context.Context, userId int) ([]*model.GetUserOrderRelationByUserIdResponse, error)
	FindById(ctx context.Context, id int) (*model.GetUserOrderResponse, error)
}

type UserService interface {
	FindAll(ctx context.Context) ([]*model.GetUserResponse, error)
	FindById(ctx context.Context, id int) (*model.GetUserResponse, error)
	FindByEmail(ctx context.Context, request *model.GetUserCredentialRequest) (*model.GetUserResponse, error)
	Create(ctx context.Context, request *model.CreateUserRequest) (*model.GetUserResponse, error)
	Update(ctx context.Context, request *model.UpdateUserRequest) (*model.GetUserResponse, error)
	Delete(ctx context.Context, id int) error
}
