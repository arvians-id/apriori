package utils

import (
	"apriori/entity"
	"apriori/model"
)

func ToUserResponse(user entity.User) model.GetUserResponse {
	return model.GetUserResponse{
		IdUser:    user.IdUser,
		Role:      user.Role,
		Name:      user.Name,
		Email:     user.Email,
		Address:   user.Address,
		Phone:     user.Phone,
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
		Price:       user.Price,
		Image:       user.Image,
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

func ToAprioriResponse(apriori entity.Apriori) model.GetAprioriResponse {
	return model.GetAprioriResponse{
		IdApriori:   apriori.IdApriori,
		Code:        apriori.Code,
		Item:        apriori.Item,
		Discount:    apriori.Discount,
		Support:     apriori.Support,
		Confidence:  apriori.Confidence,
		RangeDate:   apriori.RangeDate,
		IsActive:    apriori.IsActive,
		Description: apriori.Description,
		Image:       apriori.Image,
		CreatedAt:   apriori.CreatedAt.String(),
	}
}

func ToPaymentNullableResponse(payment entity.PaymentNullable) model.GetPaymentNullableResponse {
	return model.GetPaymentNullableResponse{
		IdPayload:         payment.IdPayload,
		UserId:            payment.UserId,
		OrderId:           payment.OrderId,
		TransactionTime:   payment.TransactionTime,
		TransactionStatus: payment.TransactionStatus,
		TransactionId:     payment.TransactionId,
		StatusCode:        payment.StatusCode,
		SignatureKey:      payment.SignatureKey,
		SettlementTime:    payment.SettlementTime,
		PaymentType:       payment.PaymentType,
		MerchantId:        payment.MerchantId,
		GrossAmount:       payment.GrossAmount,
		FraudStatus:       payment.FraudStatus,
		BankType:          payment.BankType,
		VANumber:          payment.VANumber,
		BillerCode:        payment.BillerCode,
		BillKey:           payment.BillKey,
	}
}

func ToPaymentRelationResponse(payment entity.PaymentRelation) model.GetPaymentRelationResponse {
	return model.GetPaymentRelationResponse{
		IdPayload:         payment.IdPayload,
		UserId:            payment.UserId,
		OrderId:           payment.OrderId,
		TransactionTime:   payment.TransactionTime,
		TransactionStatus: payment.TransactionStatus,
		TransactionId:     payment.TransactionId,
		StatusCode:        payment.StatusCode,
		SignatureKey:      payment.SignatureKey,
		SettlementTime:    payment.SettlementTime,
		PaymentType:       payment.PaymentType,
		MerchantId:        payment.MerchantId,
		GrossAmount:       payment.GrossAmount,
		FraudStatus:       payment.FraudStatus,
		BankType:          payment.BankType,
		VANumber:          payment.VANumber,
		BillerCode:        payment.BillerCode,
		BillKey:           payment.BillKey,
		UserName:          payment.UserName,
	}
}

func ToUserOrderResponse(userOrder entity.UserOrder) model.GetUserOrderResponse {
	return model.GetUserOrderResponse{
		IdOrder:        userOrder.IdOrder,
		PayloadId:      userOrder.PayloadId,
		Code:           userOrder.Code,
		Name:           userOrder.Name,
		Price:          userOrder.Price,
		Image:          userOrder.Image,
		Quantity:       userOrder.Quantity,
		TotalPriceItem: userOrder.TotalPriceItem,
	}
}

func ToCategoryResponse(category entity.Category) model.GetCategoryResponse {
	return model.GetCategoryResponse{
		IdCategory: category.IdCategory,
		Name:       category.Name,
		CreatedAt:  category.CreatedAt.String(),
		UpdatedAt:  category.UpdatedAt.String(),
	}
}
