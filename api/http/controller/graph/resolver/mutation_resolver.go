package resolver

import (
	"context"
	"errors"
	"fmt"
	"github.com/arvians-id/apriori/http/controller/graph/generated"
	"github.com/arvians-id/apriori/http/controller/rest/request"
	"github.com/arvians-id/apriori/http/controller/rest/response"
	"github.com/arvians-id/apriori/model"
	"os"
)

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct {
	*Resolver
}

func (r *mutationResolver) AuthLogin(ctx context.Context, input model.GetUserCredentialRequest) (*model.TokenJwt, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) AuthRegister(ctx context.Context, input model.CreateUserRequest) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) AuthRefresh(ctx context.Context, input model.GetRefreshTokenRequest) (*model.TokenJwt, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) AuthForgotPassword(ctx context.Context, input model.CreatePasswordResetRequest) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) AuthVerifyResetPassword(ctx context.Context, input model.UpdateResetPasswordUserRequest) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) AuthLogout(ctx context.Context) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) UserCreate(ctx context.Context, input model.CreateUserRequest) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) UserUpdate(ctx context.Context, input model.UpdateUserRequest) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) UserDelete(ctx context.Context, id int) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) UpdateProfile(ctx context.Context, input model.UpdateUserRequest) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) CategoryCreate(ctx context.Context, input model.CreateCategoryRequest) (*model.Category, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) CategoryUpdate(ctx context.Context, input model.UpdateCategoryRequest) (*model.Category, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) CategoryDelete(ctx context.Context, id int) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) TransactionCreate(ctx context.Context, input model.CreateTransactionRequest) (*model.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) TransactionCreateByCSV(ctx context.Context, input model.CreateTransactionRequest) (*model.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) TransactionUpdate(ctx context.Context, input model.UpdateTransactionRequest) (*model.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) TransactionDelete(ctx context.Context, numberTransaction string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) TransactionTruncate(ctx context.Context) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) PaymentUpdateReceiptNumber(ctx context.Context, input model.AddReceiptNumberRequest) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) PaymentPay(ctx context.Context, input model.GetPaymentTokenRequest) (*model.Payment, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) PaymentNotification(ctx context.Context) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) PaymentDelete(ctx context.Context, orderID string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) CommentCreate(ctx context.Context, input model.CreateCommentRequest) (*model.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) RajaOngkirCost(ctx context.Context, input model.GetDeliveryRequest) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) NotificationMarkAll(ctx context.Context) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) NotificationMark(ctx context.Context, id int) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) AprioriCreate(ctx context.Context, input model.GenerateCreateAprioriRequest) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) AprioriUpdate(ctx context.Context, input model.UpdateAprioriRequest) (*model.Apriori, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) AprioriDelete(ctx context.Context, code string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) AprioriGenerate(ctx context.Context, input model.GenerateAprioriRequest) (*model.GenerateApriori, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) AprioriUpdateStatus(ctx context.Context, code string) (bool, error) {
	//TODO implement me
	panic("implement me")
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
