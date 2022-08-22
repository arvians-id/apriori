package helper

import (
	"apriori/entity"
	"apriori/model"
)

func ToUserResponse(user *entity.User) *model.GetUserResponse {
	return &model.GetUserResponse{
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

func ToPasswordResetResponse(reset *entity.PasswordReset) *model.GetPasswordResetResponse {
	return &model.GetPasswordResetResponse{
		Email:   reset.Email,
		Token:   reset.Token,
		Expired: reset.Expired,
	}
}

func ToProductResponse(product *entity.Product) *model.GetProductResponse {
	return &model.GetProductResponse{
		IdProduct:   product.IdProduct,
		Code:        product.Code,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category,
		IsEmpty:     product.IsEmpty,
		Mass:        product.Mass,
		Image:       product.Image,
		CreatedAt:   product.CreatedAt.String(),
		UpdatedAt:   product.UpdatedAt.String(),
	}
}

func ToTransactionResponse(transaction *entity.Transaction) *model.GetTransactionResponse {
	return &model.GetTransactionResponse{
		IdTransaction: transaction.IdTransaction,
		ProductName:   transaction.ProductName,
		CustomerName:  transaction.CustomerName,
		NoTransaction: transaction.NoTransaction,
		CreatedAt:     transaction.CreatedAt.String(),
		UpdatedAt:     transaction.UpdatedAt.String(),
	}
}

func ToAprioriResponse(apriori *entity.Apriori) *model.GetAprioriResponse {
	return &model.GetAprioriResponse{
		IdApriori:   apriori.IdApriori,
		Code:        apriori.Code,
		Item:        apriori.Item,
		Discount:    apriori.Discount,
		Support:     apriori.Support,
		Confidence:  apriori.Confidence,
		RangeDate:   apriori.RangeDate,
		IsActive:    apriori.IsActive,
		Description: apriori.Description.String,
		Image:       apriori.Image,
		CreatedAt:   apriori.CreatedAt.String(),
	}
}

func ToPaymentResponse(payment *entity.Payment) *model.GetPaymentResponse {
	return &model.GetPaymentResponse{
		IdPayload:         payment.IdPayload,
		UserId:            payment.UserId.String,
		OrderId:           payment.OrderId.String,
		TransactionTime:   payment.TransactionTime.String,
		TransactionStatus: payment.TransactionStatus.String,
		TransactionId:     payment.TransactionId.String,
		StatusCode:        payment.StatusCode.String,
		SignatureKey:      payment.SignatureKey.String,
		SettlementTime:    payment.SettlementTime.String,
		PaymentType:       payment.PaymentType.String,
		MerchantId:        payment.MerchantId.String,
		GrossAmount:       payment.GrossAmount.String,
		FraudStatus:       payment.FraudStatus.String,
		BankType:          payment.BankType.String,
		VANumber:          payment.VANumber.String,
		BillerCode:        payment.BillerCode.String,
		BillKey:           payment.BillKey.String,
		ReceiptNumber:     payment.ReceiptNumber.String,
		Address:           payment.Address.String,
		Courier:           payment.Courier.String,
		CourierService:    payment.CourierService.String,
	}
}

func ToPaymentRelationResponse(payment *entity.PaymentRelation) *model.GetPaymentRelationResponse {
	return &model.GetPaymentRelationResponse{
		IdPayload:         payment.Payment.IdPayload,
		UserId:            payment.Payment.UserId.String,
		OrderId:           payment.Payment.OrderId.String,
		TransactionTime:   payment.Payment.TransactionTime.String,
		TransactionStatus: payment.Payment.TransactionStatus.String,
		TransactionId:     payment.Payment.TransactionId.String,
		StatusCode:        payment.Payment.StatusCode.String,
		SignatureKey:      payment.Payment.SignatureKey.String,
		SettlementTime:    payment.Payment.SettlementTime.String,
		PaymentType:       payment.Payment.PaymentType.String,
		MerchantId:        payment.Payment.MerchantId.String,
		GrossAmount:       payment.Payment.GrossAmount.String,
		FraudStatus:       payment.Payment.FraudStatus.String,
		BankType:          payment.Payment.BankType.String,
		VANumber:          payment.Payment.VANumber.String,
		BillerCode:        payment.Payment.BillerCode.String,
		BillKey:           payment.Payment.BillKey.String,
		ReceiptNumber:     payment.Payment.ReceiptNumber.String,
		Address:           payment.Payment.Address.String,
		Courier:           payment.Payment.Courier.String,
		CourierService:    payment.Payment.CourierService.String,
		UserName:          payment.UserName.String,
	}
}

func ToUserOrderResponse(userOrder *entity.UserOrder) *model.GetUserOrderResponse {
	return &model.GetUserOrderResponse{
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

func ToUserOrderRelationByUserIdResponse(userOrder *entity.UserOrderRelationByUserId) *model.GetUserOrderRelationByUserIdResponse {
	return &model.GetUserOrderRelationByUserIdResponse{
		IdOrder:           userOrder.UserOrder.IdOrder,
		PayloadId:         userOrder.UserOrder.PayloadId,
		Code:              userOrder.UserOrder.Code,
		Name:              userOrder.UserOrder.Name,
		Price:             userOrder.UserOrder.Price,
		Image:             userOrder.UserOrder.Image,
		Quantity:          userOrder.UserOrder.Quantity,
		TotalPriceItem:    userOrder.UserOrder.TotalPriceItem,
		OrderId:           userOrder.OrderId,
		TransactionStatus: userOrder.TransactionStatus,
	}
}

func ToCategoryResponse(category *entity.Category) *model.GetCategoryResponse {
	return &model.GetCategoryResponse{
		IdCategory: category.IdCategory,
		Name:       category.Name,
		CreatedAt:  category.CreatedAt.String(),
		UpdatedAt:  category.UpdatedAt.String(),
	}
}

func ToCommentResponse(comment *entity.Comment) *model.GetCommentResponse {
	return &model.GetCommentResponse{
		IdComment:   comment.IdComment,
		UserOrderId: comment.UserOrderId,
		ProductCode: comment.ProductCode,
		Description: comment.Description.String,
		Tag:         comment.Tag.String,
		Rating:      comment.Rating,
		CreatedAt:   comment.CreatedAt.String(),
		UserId:      comment.UserId,
		UserName:    comment.UserName,
	}
}

func ToRatingResponse(rating *entity.RatingFromComment) *model.GetRatingResponse {
	return &model.GetRatingResponse{
		Rating:        rating.Rating,
		ResultRating:  rating.ResultRating,
		ResultComment: rating.ResultComment,
	}
}
