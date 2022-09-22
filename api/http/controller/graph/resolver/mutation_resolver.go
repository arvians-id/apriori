package resolver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/arvians-id/apriori/helper"
	"github.com/arvians-id/apriori/http/controller/graph/generated"
	"github.com/arvians-id/apriori/http/controller/rest/request"
	"github.com/arvians-id/apriori/http/controller/rest/response"
	"github.com/arvians-id/apriori/http/middleware"
	"github.com/arvians-id/apriori/model"
	"github.com/go-redis/redis/v8"
	"github.com/veritrans/go-midtrans"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct {
	*Resolver
}

func (r *mutationResolver) AuthLogin(ctx context.Context, input model.GetUserCredentialRequest) (*model.TokenJwt, error) {
	ginContext, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := r.UserService.FindByEmail(ctx, (*request.GetUserCredentialRequest)(&input))
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return nil, errors.New(response.ResponseErrorNotFound)
		}

		if err.Error() == response.WrongPassword {
			return nil, errors.New(response.WrongPassword)
		}

		return nil, err
	}

	expiredTimeAccess, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRED_TIME"))
	expirationTime := time.Now().Add(time.Duration(expiredTimeAccess) * 24 * time.Hour)
	token, err := r.JwtService.GenerateToken(user.IdUser, expirationTime)
	if err != nil {
		return nil, err
	}

	http.SetCookie(ginContext.Writer, &http.Cookie{
		Name:     "token",
		Value:    url.QueryEscape(token.AccessToken),
		Expires:  expirationTime,
		Path:     "/",
		HttpOnly: true,
	})

	return &model.TokenJwt{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (r *mutationResolver) AuthRegister(ctx context.Context, input model.CreateUserRequest) (*model.User, error) {
	user, err := r.UserService.Create(ctx, (*request.CreateUserRequest)(&input))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *mutationResolver) AuthRefresh(ctx context.Context, input model.GetRefreshTokenRequest) (*model.TokenJwt, error) {
	ginContext, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	token, err := r.JwtService.RefreshToken(input.RefreshToken)
	if err != nil {
		return nil, err
	}

	expiredTimeAccess, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRED_TIME"))
	expirationTime := time.Now().Add(time.Duration(expiredTimeAccess) * 24 * time.Hour)
	http.SetCookie(ginContext.Writer, &http.Cookie{
		Name:     "token",
		Value:    url.QueryEscape(token.AccessToken),
		Expires:  expirationTime,
		Path:     "/",
		HttpOnly: true,
	})

	return &model.TokenJwt{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (r *mutationResolver) AuthForgotPassword(ctx context.Context, input model.CreatePasswordResetRequest) (string, error) {
	result, err := r.PasswordResetService.CreateOrUpdateByEmail(ctx, input.Email)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return "", errors.New(response.ResponseErrorNotFound)
		}

		return "", err
	}

	message := fmt.Sprintf("%s/auth/reset-password?signature=%v", os.Getenv("APP_URL_FE"), result.Token)
	err = r.EmailService.SendEmailWithText(result.Email, "Forgot Password", &message)
	if err != nil {
		return "", err
	}

	return result.Token, nil
}

func (r *mutationResolver) AuthVerifyResetPassword(ctx context.Context, input model.UpdateResetPasswordUserRequest) (bool, error) {
	err := r.PasswordResetService.Verify(ctx, (*request.UpdateResetPasswordUserRequest)(&input))
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return false, errors.New(response.ResponseErrorNotFound)
		}

		if err.Error() == response.VerificationExpired {
			return false, errors.New(response.VerificationExpired)
		}

		return false, err
	}

	return true, nil
}

func (r *mutationResolver) AuthLogout(ctx context.Context) (bool, error) {
	ginContext, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}

	http.SetCookie(ginContext.Writer, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	return true, nil
}

func (r *mutationResolver) UserCreate(ctx context.Context, input model.CreateUserRequest) (*model.User, error) {
	user, err := r.UserService.Create(ctx, (*request.CreateUserRequest)(&input))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *mutationResolver) UserUpdate(ctx context.Context, input model.UpdateUserRequest) (*model.User, error) {
	user, err := r.UserService.Update(ctx, &request.UpdateUserRequest{
		IdUser:   *input.IdUser,
		Name:     input.Name,
		Email:    input.Email,
		Address:  input.Address,
		Phone:    input.Phone,
		Password: input.Password,
	})
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return nil, errors.New(response.ResponseErrorNotFound)
		}

		return nil, err
	}

	return user, nil
}

func (r *mutationResolver) UserDelete(ctx context.Context, id int) (bool, error) {
	err := r.UserService.Delete(ctx, id)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return false, errors.New(response.ResponseErrorNotFound)
		}

		return false, err
	}

	return true, nil
}

func (r *mutationResolver) UpdateProfile(ctx context.Context, input model.UpdateUserRequest) (*model.User, error) {
	ginContext, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	id, isExist := ginContext.Get("id_user")
	if !isExist {
		return nil, errors.New("unauthorized")
	}

	user, err := r.UserService.Update(ctx, &request.UpdateUserRequest{
		IdUser:   int(id.(float64)),
		Name:     input.Name,
		Email:    input.Email,
		Address:  input.Address,
		Phone:    input.Phone,
		Password: input.Password,
	})
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return nil, errors.New(response.ResponseErrorNotFound)
		}

		return nil, err
	}

	return user, nil
}

func (r *mutationResolver) CategoryCreate(ctx context.Context, input model.CreateCategoryRequest) (*model.Category, error) {
	category, err := r.CategoryService.Create(ctx, (*request.CreateCategoryRequest)(&input))
	if err != nil {
		return nil, err
	}

	// delete previous cache
	_ = r.CacheService.Del(ctx, fmt.Sprintf("categories"))

	return category, nil
}

func (r *mutationResolver) CategoryUpdate(ctx context.Context, input model.UpdateCategoryRequest) (*model.Category, error) {
	category, err := r.CategoryService.Update(ctx, (*request.UpdateCategoryRequest)(&input))
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return nil, errors.New(response.ResponseErrorNotFound)
		}

		return nil, err
	}

	// delete previous cache
	_ = r.CacheService.Del(ctx, fmt.Sprintf("categories"))

	return category, nil
}

func (r *mutationResolver) CategoryDelete(ctx context.Context, id int) (bool, error) {
	err := r.CategoryService.Delete(ctx, id)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return false, errors.New(response.ResponseErrorNotFound)
		}

		return false, err
	}

	// delete previous cache
	_ = r.CacheService.Del(ctx, fmt.Sprintf("categories"))

	return true, nil
}

func (r *mutationResolver) TransactionCreate(ctx context.Context, input model.CreateTransactionRequest) (*model.Transaction, error) {
	transaction, err := r.TransactionService.Create(ctx, (*request.CreateTransactionRequest)(&input))
	if err != nil {
		return nil, err
	}

	// delete previous cache
	_ = r.CacheService.Del(ctx, "all-transaction")

	return transaction, nil
}

func (r *mutationResolver) TransactionCreateByCSV(ctx context.Context, file graphql.Upload) (bool, error) {
	initFileName := fmt.Sprintf("%v-%s", file.Size, file.Filename)
	fileNameGenerated, err := r.StorageService.UploadFileS3GraphQL(file, initFileName)
	if err != nil {
		return false, err
	}
	defer os.Remove(initFileName)

	data, err := helper.OpenCsvFile(fileNameGenerated)
	if err != nil {
		return false, err
	}

	err = r.TransactionService.CreateByCsv(ctx, data)
	if err != nil {
		return false, err
	}

	// delete previous cache
	_ = r.CacheService.Del(ctx, "all-transaction")

	return true, nil
}

func (r *mutationResolver) TransactionUpdate(ctx context.Context, input model.UpdateTransactionRequest) (*model.Transaction, error) {
	transaction, err := r.TransactionService.Update(ctx, (*request.UpdateTransactionRequest)(&input))
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return nil, errors.New(response.ResponseErrorNotFound)
		}

		return nil, err
	}

	// delete previous cache
	_ = r.CacheService.Del(ctx, "all-transaction")

	return transaction, nil
}

func (r *mutationResolver) TransactionDelete(ctx context.Context, numberTransaction string) (bool, error) {
	err := r.TransactionService.Delete(ctx, numberTransaction)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return false, errors.New(response.ResponseErrorNotFound)
		}

		return false, err
	}

	// delete previous cache
	_ = r.CacheService.Del(ctx, "all-transaction")

	return true, nil
}

func (r *mutationResolver) TransactionTruncate(ctx context.Context) (bool, error) {
	err := r.TransactionService.Truncate(ctx)
	if err != nil {
		return false, nil
	}

	// delete previous cache
	_ = r.CacheService.Del(ctx, "all-transaction")

	return true, nil
}

func (r *mutationResolver) PaymentUpdateReceiptNumber(ctx context.Context, input model.AddReceiptNumberRequest) (bool, error) {
	payment, err := r.PaymentService.UpdateReceiptNumber(ctx, &request.AddReceiptNumberRequest{
		OrderId:       input.OrderId,
		ReceiptNumber: input.ReceiptNumber,
	})
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return false, errors.New(response.ResponseErrorNotFound)
		}

		return false, err
	}

	// Notification
	var notificationRequest model.CreateNotificationRequest
	notificationRequest.UserId = payment.UserId
	notificationRequest.Title = "Receipt number arrived"
	notificationRequest.Description = "Your receipt number has been entered by the admin"
	notificationRequest.URL = "product"
	err = r.NotificationService.Create(ctx, &request.CreateNotificationRequest{
		UserId:      notificationRequest.UserId,
		Title:       notificationRequest.Title,
		Description: notificationRequest.Description,
		URL:         notificationRequest.URL,
	}).WithSendMail()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) PaymentPay(ctx context.Context, input model.GetPaymentTokenRequest) (map[string]interface{}, error) {
	data, err := r.PaymentService.GetToken(ctx, (*request.GetPaymentTokenRequest)(&input))
	if err != nil {
		return nil, err
	}

	// delete previous cache
	key := fmt.Sprintf("user-order-payment-%v", input.UserId)
	_ = r.CacheService.Del(ctx, key)

	return data, nil
}

func (r *mutationResolver) PaymentNotification(ctx context.Context) (bool, error) {
	ginContext, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}

	var payload midtrans.ChargeReqWithMap
	err = ginContext.BindJSON(&payload)
	if err != nil {
		return false, err
	}

	encode, _ := json.Marshal(payload)
	resArray := make(map[string]interface{})
	err = json.Unmarshal(encode, &resArray)

	err = r.PaymentService.CreateOrUpdate(ctx, resArray)
	if err != nil {
		return false, err
	}

	// delete previous cache
	key := fmt.Sprintf("user-order-id-%v", helper.StrToInt(resArray["custom_field2"].(string)))
	key2 := fmt.Sprintf("user-order-payment-%v", helper.StrToInt(resArray["custom_field1"].(string)))
	key3 := fmt.Sprintf("user-order-rate-%v", helper.StrToInt(resArray["custom_field1"].(string)))
	_ = r.CacheService.Del(ctx, key, key2, key3)

	return true, nil
}

func (r *mutationResolver) PaymentDelete(ctx context.Context, orderID string) (bool, error) {
	err := r.PaymentService.Delete(ctx, orderID)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) CommentCreate(ctx context.Context, input model.CreateCommentRequest) (*model.Comment, error) {
	comment, err := r.CommentService.Create(ctx, (*request.CreateCommentRequest)(&input))
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *mutationResolver) RajaOngkirCost(ctx context.Context, input model.GetDeliveryRequest) (interface{}, error) {
	payload := fmt.Sprintf(
		"origin=%v&destination=%v&weight=%v&courier=%v",
		input.Origin,
		input.Destination,
		input.Weight,
		input.Courier,
	)
	data := strings.NewReader(payload)
	req, _ := http.NewRequest("POST", "https://api.rajaongkir.com/starter/cost", data)
	req.Header.Add("key", os.Getenv("RAJA_ONGKIR_SECRET_KEY"))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	var rajaOngkirModel interface{}
	err := json.NewDecoder(res.Body).Decode(&rajaOngkirModel)
	if err != nil {
		return nil, err
	}

	return rajaOngkirModel, nil
}

func (r *mutationResolver) NotificationMarkAll(ctx context.Context) (bool, error) {
	ginContext, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}

	id, isExist := ginContext.Get("id_user")
	if !isExist {
		return false, errors.New("unauthorized")
	}

	err = r.NotificationService.MarkAll(ctx, int(id.(float64)))
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) NotificationMark(ctx context.Context, id int) (bool, error) {
	err := r.NotificationService.Mark(ctx, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) AprioriCreate(ctx context.Context, input []*model.GenerateCreateAprioriRequest) (bool, error) {
	var aprioriRequests []*request.CreateAprioriRequest
	for _, generateRequest := range input {
		ItemSet := strings.Join(generateRequest.ItemSet, ", ")
		aprioriRequests = append(aprioriRequests, &request.CreateAprioriRequest{
			Item:       ItemSet,
			Discount:   generateRequest.Discount,
			Support:    generateRequest.Support,
			Confidence: generateRequest.Confidence,
			RangeDate:  generateRequest.RangeDate,
		})
	}

	err := r.AprioriService.Create(ctx, aprioriRequests)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) AprioriUpdate(ctx context.Context, input model.UpdateAprioriRequest) (*model.Apriori, error) {
	var fileName string
	if input.Image.Filename != "" {
		initFileName := fmt.Sprintf("%v-%s", input.Image.Size, input.Image.Filename)
		fileNameGenerated, err := r.StorageService.UploadFileS3GraphQL(input.Image, initFileName)
		if err != nil {
			return nil, err
		}
		defer os.Remove(initFileName)
		fileName = fileNameGenerated
	}

	apriories, err := r.AprioriService.Update(ctx, &request.UpdateAprioriRequest{
		IdApriori:   input.IdApriori,
		Code:        input.Code,
		Description: input.Description,
		Image:       fileName,
	})
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return nil, errors.New(response.ResponseErrorNotFound)
		}

		return nil, err
	}

	return apriories, nil
}

func (r *mutationResolver) AprioriDelete(ctx context.Context, code string) (bool, error) {
	err := r.AprioriService.Delete(ctx, code)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return false, errors.New(response.ResponseErrorNotFound)
		}

		return false, err
	}

	return true, nil
}

func (r *mutationResolver) AprioriGenerate(ctx context.Context, input model.GenerateAprioriRequest) ([]*model.GenerateApriori, error) {
	key := fmt.Sprintf(
		"%v%v%v%v%s%s",
		input.MinimumDiscount,
		input.MaximumDiscount,
		input.MinimumSupport,
		input.MinimumConfidence,
		input.StartDate,
		input.EndDate,
	)
	aprioriCache, err := r.CacheService.Get(ctx, key)
	if err == redis.Nil {
		apriori, err := r.AprioriService.Generate(ctx, (*request.GenerateAprioriRequest)(&input))
		if err != nil {
			return nil, err
		}

		err = r.CacheService.Set(ctx, key, apriori)
		if err != nil {
			return nil, err
		}

		return apriori, nil
	} else if err != nil {
		return nil, err
	}

	var aprioriCacheResponses []*model.GenerateApriori
	err = json.Unmarshal(aprioriCache, &aprioriCacheResponses)
	if err != nil {
		return nil, err
	}

	return aprioriCacheResponses, nil
}

func (r *mutationResolver) AprioriUpdateStatus(ctx context.Context, code string) (bool, error) {
	err := r.AprioriService.UpdateStatus(ctx, code)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return false, errors.New(response.ResponseErrorNotFound)
		}

		return false, err
	}

	return true, nil
}

func (r *mutationResolver) ProductCreate(ctx context.Context, input model.CreateProductRequest) (*model.Product, error) {
	fileName := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/assets/%s", os.Getenv("AWS_BUCKET"), os.Getenv("AWS_REGION"), "no-image.png")
	if input.Image.Filename != "" {
		initFileName := fmt.Sprintf("%v-%s", input.Image.Size, input.Image.Filename)
		fileNameGenerated, err := r.StorageService.UploadFileS3GraphQL(input.Image, initFileName)
		if err != nil {
			return nil, err
		}
		defer os.Remove(initFileName)
		fileName = fileNameGenerated
	}

	product, err := r.ProductService.Create(ctx, &request.CreateProductRequest{
		Code:        input.Code,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Category:    input.Category,
		Mass:        input.Mass,
		Image:       fileName,
	})
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *mutationResolver) ProductUpdate(ctx context.Context, input model.UpdateProductRequest) (*model.Product, error) {
	var fileName string
	if input.Image.Filename != "" {
		initFileName := fmt.Sprintf("%v-%s", input.Image.Size, input.Image.Filename)
		fileNameGenerated, err := r.StorageService.UploadFileS3GraphQL(input.Image, initFileName)
		if err != nil {
			return nil, err
		}
		defer os.Remove(initFileName)
		fileName = fileNameGenerated
	}

	product, err := r.ProductService.Update(ctx, &request.UpdateProductRequest{
		Code:        input.Code,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Category:    input.Category,
		IsEmpty:     input.IsEmpty,
		Mass:        input.Mass,
		Image:       fileName,
	})
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return nil, errors.New(response.ResponseErrorNotFound)
		}
		return nil, err
	}

	// delete previous cache
	_ = r.CacheService.Del(ctx, fmt.Sprintf("product-%s", product.Code))

	return product, nil
}

func (r *mutationResolver) ProductDelete(ctx context.Context, code string) (bool, error) {
	err := r.ProductService.Delete(ctx, code)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			return false, errors.New(response.ResponseErrorNotFound)
		}
		return false, err
	}

	// delete previous cache
	_ = r.CacheService.Del(ctx, fmt.Sprintf("product-%s", code))

	return true, nil
}
