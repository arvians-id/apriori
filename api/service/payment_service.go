package service

import (
	"apriori/config"
	"apriori/entity"
	"apriori/helper"
	"apriori/model"
	"apriori/repository"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/veritrans/go-midtrans"
	"reflect"
	"strings"
	"time"
)

type PaymentServiceImpl struct {
	MidClient             midtrans.Client
	SnapGateway           midtrans.SnapGateway
	ServerKey             string
	ClientKey             string
	PaymentRepository     repository.PaymentRepository
	UserOrderRepository   repository.UserOrderRepository
	TransactionRepository repository.TransactionRepository
	NotificationService   NotificationService
	DB                    *sql.DB
}

func NewPaymentService(
	configuration config.Config,
	paymentRepository *repository.PaymentRepository,
	userOrderRepository *repository.UserOrderRepository,
	transactionRepository *repository.TransactionRepository,
	notificationService *NotificationService,
	db *sql.DB,
) PaymentService {
	midClient := midtrans.NewClient()
	midClient.ServerKey = configuration.Get("MIDTRANS_SERVER_KEY")
	midClient.ClientKey = configuration.Get("MIDTRANS_CLIENT_KEY")
	midClient.APIEnvType = midtrans.Sandbox

	return &PaymentServiceImpl{
		MidClient:             midClient,
		ServerKey:             midClient.ServerKey,
		ClientKey:             midClient.ClientKey,
		PaymentRepository:     *paymentRepository,
		UserOrderRepository:   *userOrderRepository,
		TransactionRepository: *transactionRepository,
		NotificationService:   *notificationService,
		DB:                    db,
	}
}

func (service *PaymentServiceImpl) GetClient() {
	service.SnapGateway = midtrans.SnapGateway{
		Client: service.MidClient,
	}
}

func (service *PaymentServiceImpl) FindAll(ctx context.Context) ([]*entity.Payment, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	payments, err := service.PaymentRepository.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (service *PaymentServiceImpl) FindAllByUserId(ctx context.Context, userId int) ([]*entity.Payment, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	payments, err := service.PaymentRepository.FindAllByUserId(ctx, tx, userId)
	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (service *PaymentServiceImpl) FindByOrderId(ctx context.Context, orderId string) (*entity.Payment, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	payment, err := service.PaymentRepository.FindByOrderId(ctx, tx, orderId)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (service *PaymentServiceImpl) CreateOrUpdate(ctx context.Context, request map[string]interface{}) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	var bankType, vaNumber, billerCode, billKey, settlementTime string

	if request["va_numbers"] != nil {
		bankType = request["va_numbers"].([]interface{})[0].(map[string]interface{})["bank"].(string)
		vaNumber = request["va_numbers"].([]interface{})[0].(map[string]interface{})["va_number"].(string)
	} else if request["permata_va_number"] != nil {
		bankType = "permata bank"
		vaNumber = request["permata_va_number"].(string)
	} else if request["biller_code"] != nil && request["bill_key"] != nil {
		billerCode = request["biller_code"].(string)
		billKey = request["bill_key"].(string)
		bankType = "mandiri"
	}

	setTime, ok := request["settlement_time"]
	if ok {
		settlementTime = setTime.(string)
	} else {
		settlementTime = ""
	}

	orderID := request["order_id"].(string)
	transactionTime := request["transaction_time"].(string)
	transactionStatus := request["transaction_status"].(string)
	transactionId := request["transaction_id"].(string)
	statusCode := request["status_code"].(string)
	signatureKey := request["signature_key"].(string)
	paymentType := request["payment_type"].(string)
	merchantId := request["merchant_id"].(string)
	grossAmount := request["gross_amount"].(string)
	fraudStatus := request["fraud_status"].(string)

	checkTransaction, _ := service.PaymentRepository.FindByOrderId(ctx, tx, request["order_id"].(string))
	checkTransaction.UserId = request["custom_field1"].(string)
	checkTransaction.OrderId = &orderID
	checkTransaction.TransactionTime = &transactionTime
	checkTransaction.TransactionStatus = &transactionStatus
	checkTransaction.TransactionId = &transactionId
	checkTransaction.StatusCode = &statusCode
	checkTransaction.SignatureKey = &signatureKey
	checkTransaction.SettlementTime = &settlementTime
	checkTransaction.PaymentType = &paymentType
	checkTransaction.MerchantId = &merchantId
	checkTransaction.GrossAmount = &grossAmount
	checkTransaction.FraudStatus = &fraudStatus
	checkTransaction.BankType = &bankType
	checkTransaction.VANumber = &vaNumber
	checkTransaction.BillerCode = &billerCode
	checkTransaction.BillKey = &billKey

	if checkTransaction.OrderId != nil {
		err := service.PaymentRepository.Update(ctx, tx, checkTransaction)
		if err != nil {
			return err
		}

		if request["transaction_status"].(string) == "settlement" {
			timeNow, err := time.Parse(helper.TimeFormat, time.Now().Format(helper.TimeFormat))
			if err != nil {
				return err
			}

			userOrder, err := service.UserOrderRepository.FindAllByPayloadId(ctx, tx, request["custom_field2"].(string))
			if err != nil {
				return err
			}

			var productName []string
			for _, item := range userOrder {
				productName = append(productName, *item.Name)
			}

			transaction := entity.Transaction{
				ProductName:   strings.ToLower(strings.Join(productName, ", ")),
				CustomerName:  request["custom_field3"].(string),
				NoTransaction: request["order_id"].(string),
				CreatedAt:     timeNow,
				UpdatedAt:     timeNow,
			}

			_, err = service.TransactionRepository.Create(ctx, tx, &transaction)
			if err != nil {
				return err
			}

			var notificationRequest model.CreateNotificationRequest
			notificationRequest.UserId = helper.StrToInt(checkTransaction.UserId)
			notificationRequest.Title = "Transaction Successfully"
			notificationRequest.Description = "You have successfully made a payment. Thank you for shopping at Ryzy Shop"
			notificationRequest.URL = "product"
			err = service.NotificationService.Create(ctx, &notificationRequest).WithSendMail()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (service *PaymentServiceImpl) UpdateReceiptNumber(ctx context.Context, request *model.AddReceiptNumberRequest) (*entity.Payment, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	payment, err := service.PaymentRepository.FindByOrderId(ctx, tx, request.OrderId)
	if err != nil {
		return nil, err
	}

	payment.ReceiptNumber = &request.ReceiptNumber

	err = service.PaymentRepository.UpdateReceiptNumber(ctx, tx, payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (service *PaymentServiceImpl) Delete(ctx context.Context, orderId string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	payment, err := service.PaymentRepository.FindByOrderId(ctx, tx, orderId)
	if err != nil {
		return err
	}

	err = service.PaymentRepository.Delete(ctx, tx, payment.OrderId)
	if err != nil {
		return err
	}

	return nil
}
func (service *PaymentServiceImpl) GetToken(ctx context.Context, request *model.GetPaymentTokenRequest) (map[string]interface{}, error) {
	service.GetClient()

	var items []map[string]interface{}
	for _, item := range request.Items {
		err := json.Unmarshal([]byte(item), &items)
		if err != nil {
			return nil, err
		}
	}

	var itemDetails []midtrans.ItemDetail
	for _, item := range items {
		code := CheckCode(item)
		itemDetails = append(itemDetails, midtrans.ItemDetail{
			ID:    code,
			Name:  item["name"].(string),
			Price: int64(item["price"].(float64)),
			Qty:   int32(item["quantity"].(float64)),
		})
	}

	itemDetails = append(itemDetails, midtrans.ItemDetail{
		ID:    itemDetails[len(itemDetails)-1].ID,
		Name:  "Pajak",
		Price: 5000,
		Qty:   1,
	})

	itemDetails = append(itemDetails, midtrans.ItemDetail{
		ID:    itemDetails[len(itemDetails)-1].ID,
		Name:  "Ongkos Kirim",
		Price: request.ShippingCost,
		Qty:   1,
	})

	orderID := helper.RandomString(20)
	userId := helper.IntToStr(request.UserId)
	var snapRequest midtrans.SnapReq
	snapRequest.TransactionDetails.OrderID = orderID
	snapRequest.TransactionDetails.GrossAmt = request.GrossAmount
	snapRequest.Items = &itemDetails
	snapRequest.CustomerDetail = &midtrans.CustDetail{
		FName: request.CustomerName,
	}
	snapRequest.CustomField1 = userId

	// Save to database
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	canceled := "canceled"
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	paymentRequest := entity.Payment{
		UserId:            userId,
		OrderId:           &orderID,
		TransactionStatus: &canceled,
		TransactionTime:   &timeNow,
		Address:           &request.Address,
		Courier:           &request.Courier,
		CourierService:    &request.CourierService,
	}
	payment, err := service.PaymentRepository.Create(ctx, tx, &paymentRequest)
	if err != nil {
		return nil, err
	}
	// Send id payload
	snapRequest.CustomField2 = helper.IntToStr(payment.IdPayload)
	snapRequest.CustomField3 = request.CustomerName

	token, err := service.SnapGateway.GetToken(&snapRequest)
	if err != nil {
		return nil, err
	}

	var tokenResponse = make(map[string]interface{})
	tokenResponse["clientKey"] = service.ClientKey
	tokenResponse["token"] = token.Token

	for _, item := range items {
		code := CheckCode(item)
		price := int64(item["price"].(float64))
		quantity := int(item["quantity"].(float64))
		totalPriceItem := int64(item["totalPricePerItem"].(float64))
		name := item["name"].(string)
		image := item["image"].(string)
		userOrder := entity.UserOrder{
			PayloadId:      payment.IdPayload,
			Code:           &code,
			Name:           &name,
			Price:          &price,
			Image:          &image,
			Quantity:       &quantity,
			TotalPriceItem: &totalPriceItem,
		}
		_, err := service.UserOrderRepository.Create(ctx, tx, &userOrder)
		if err != nil {
			return nil, err
		}
	}

	return tokenResponse, nil
}

func CheckCode(value map[string]interface{}) string {
	checkCode := reflect.ValueOf(value["code"]).Kind()
	var code string
	if checkCode == reflect.Float64 {
		code = helper.IntToStr(int(value["code"].(float64)))
	} else if checkCode == reflect.String {
		code = value["code"].(string)
	}

	return code
}
