package utils

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

func ToTransactionResponse(transaction entity.Transaction) model.GetTransactionResponse {
	return model.GetTransactionResponse{
		IdTransaction: transaction.IdTransaction,
		ProductName:   transaction.ProductName,
		CustomerName:  transaction.CustomerName,
		NoTransaction: transaction.NoTransaction,
		CreatedAt:     transaction.CreatedAt.String(),
		UpdatedAt:     transaction.UpdatedAt.String(),
	}
}
