package helper

import (
	"apriori/entity"
	"apriori/model"
)

func ToUserResponse(user entity.User) model.GetUserResponse {
	return model.GetUserResponse{
		IdUser:    user.IdUser,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}
}

func ToPasswordResetResponse(reset entity.PasswordReset) model.GetPasswordResetResponse {
	return model.GetPasswordResetResponse{
		Email:   reset.Email,
		Token:   reset.Token,
		Expired: reset.Expired,
	}
}

func ToProductResponse(user entity.Product) model.GetProductResponse {
	return model.GetProductResponse{
		IdProduct:   user.IdProduct,
		Code:        user.Code,
		Name:        user.Name,
		Description: user.Description,
		CreatedAt:   user.CreatedAt.String(),
		UpdatedAt:   user.UpdatedAt.String(),
	}
}

func ToTransactionProductResponse(transaction entity.TransactionProduct) model.GetTransactionProductResponse {
	return model.GetTransactionProductResponse{
		TransactionId:           transaction.TransactionId,
		TransactionCustomerName: transaction.TransactionCustomerName,
		TransactionNo:           transaction.TransactionNo,
		TransactionQuantity:     transaction.TransactionQuantity,
		TransactionCreatedAt:    transaction.TransactionCreatedAt.String(),
		ProductId:               transaction.ProductId,
		ProductCode:             transaction.ProductCode,
		ProductName:             transaction.ProductName,
		ProductDescription:      transaction.ProductDescription,
	}
}

func ToProductTransactionResponse(product entity.ProductTransaction) model.GetProductTransactionResponse {
	return model.GetProductTransactionResponse{
		Code:        product.Code,
		ProductName: product.ProductName,
		Transaction: product.Transaction,
		Support:     product.Support,
	}
}
