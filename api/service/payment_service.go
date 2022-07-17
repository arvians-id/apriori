package service

import (
	"apriori/config"
	"apriori/entity"
	"apriori/model"
	"apriori/repository"
	"apriori/utils"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/veritrans/go-midtrans"
	"reflect"
	"strings"
	"time"
)

type PaymentService interface {
	FindAll(ctx context.Context) ([]model.GetPaymentRelationResponse, error)
	FindAllByUserId(ctx context.Context, userId int) ([]model.GetPaymentNullableResponse, error)
	FindByOrderId(ctx context.Context, orderId string) (model.GetPaymentNullableResponse, error)
	Delete(ctx context.Context, orderId string) error
	GetClient()
	GetToken(ctx context.Context, amount int64, userId int, customerName string, items []string) (map[string]interface{}, error)
	CreateOrUpdate(ctx context.Context, request map[string]interface{}) error
}

type paymentService struct {
	DB                    *sql.DB
	MidClient             midtrans.Client
	CoreGateway           midtrans.CoreGateway
	SnapGateway           midtrans.SnapGateway
	ServerKey             string
	ClientKey             string
	PaymentRepository     repository.PaymentRepository
	UserOrderRepository   repository.UserOrderRepository
	TransactionRepository repository.TransactionRepository
	date                  string
}

func NewPaymentService(configuration config.Config, paymentRepository *repository.PaymentRepository, userOrderRepository *repository.UserOrderRepository, transactionRepository *repository.TransactionRepository, db *sql.DB) PaymentService {
	midClient := midtrans.NewClient()
	midClient.ServerKey = configuration.Get("MIDTRANS_SERVER_KEY")
	midClient.ClientKey = configuration.Get("MIDTRANS_CLIENT_KEY")
	midClient.APIEnvType = midtrans.Sandbox

	return &paymentService{
		MidClient:             midClient,
		ServerKey:             midClient.ServerKey,
		ClientKey:             midClient.ClientKey,
		PaymentRepository:     *paymentRepository,
		UserOrderRepository:   *userOrderRepository,
		TransactionRepository: *transactionRepository,
		DB:                    db,
		date:                  "2006-01-02 15:04:05",
	}
}

func (service *paymentService) FindAll(ctx context.Context) ([]model.GetPaymentRelationResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	payments, err := service.PaymentRepository.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	var paymentResponse []model.GetPaymentRelationResponse
	for _, payment := range payments {
		paymentResponse = append(paymentResponse, utils.ToPaymentRelationResponse(payment))
	}

	return paymentResponse, nil
}

func (service *paymentService) FindAllByUserId(ctx context.Context, userId int) ([]model.GetPaymentNullableResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	payments, err := service.PaymentRepository.FindAllByUserId(ctx, tx, userId)
	if err != nil {
		return nil, err
	}

	var paymentResponse []model.GetPaymentNullableResponse
	for _, payment := range payments {
		paymentResponse = append(paymentResponse, utils.ToPaymentNullableResponse(payment))
	}

	return paymentResponse, nil
}

func (service *paymentService) FindByOrderId(ctx context.Context, orderId string) (model.GetPaymentNullableResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return model.GetPaymentNullableResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	payment, err := service.PaymentRepository.FindByOrderId(ctx, tx, orderId)
	if err != nil {
		return model.GetPaymentNullableResponse{}, err
	}

	return utils.ToPaymentNullableResponse(payment), nil
}

func (service *paymentService) Delete(ctx context.Context, orderId string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	err = service.PaymentRepository.Delete(ctx, tx, orderId)
	if err != nil {
		return err
	}

	return nil
}

func (service *paymentService) GetClient() {
	service.CoreGateway = midtrans.CoreGateway{
		Client: service.MidClient,
	}
	service.SnapGateway = midtrans.SnapGateway{
		Client: service.MidClient,
	}
}

func (service *paymentService) CreateOrUpdate(ctx context.Context, request map[string]interface{}) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

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

	paymentRequest := entity.Payment{
		UserId:            request["custom_field1"].(string),
		OrderId:           request["order_id"].(string),
		TransactionTime:   request["transaction_time"].(string),
		TransactionStatus: request["transaction_status"].(string),
		TransactionId:     request["transaction_id"].(string),
		StatusCode:        request["status_code"].(string),
		SignatureKey:      request["signature_key"].(string),
		SettlementTime:    settlementTime,
		PaymentType:       request["payment_type"].(string),
		MerchantId:        request["merchant_id"].(string),
		GrossAmount:       request["gross_amount"].(string),
		FraudStatus:       request["fraud_status"].(string),
		BankType:          bankType,
		VANumber:          vaNumber,
		BillerCode:        billerCode,
		BillKey:           billKey,
	}

	checkTransaction, _ := service.PaymentRepository.FindByOrderId(ctx, tx, request["order_id"].(string))
	if checkTransaction.OrderId != nil {
		err := service.PaymentRepository.Update(ctx, tx, paymentRequest)
		if err != nil {
			return err
		}
		if request["transaction_status"].(string) == "settlement" {
			timeNow, err := time.Parse(service.date, time.Now().Format(service.date))
			if err != nil {
				return err
			}
			userOrder, err := service.UserOrderRepository.FindAll(ctx, tx, request["custom_field2"].(string))
			if err != nil {
				return err
			}
			var productName []string
			for _, item := range userOrder {
				productName = append(productName, item.Name)
			}
			transaction := entity.Transaction{
				ProductName:   strings.ToLower(strings.Join(productName, ", ")),
				CustomerName:  request["custom_field3"].(string),
				NoTransaction: request["order_id"].(string),
				CreatedAt:     timeNow,
				UpdatedAt:     timeNow,
			}

			_, err = service.TransactionRepository.Create(ctx, tx, transaction)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (service *paymentService) GetToken(ctx context.Context, amount int64, userId int, customerName string, items []string) (map[string]interface{}, error) {
	service.GetClient()

	var test []map[string]interface{}
	for _, values := range items {
		err := json.Unmarshal([]byte(values), &test)
		if err != nil {
			return nil, err
		}
	}

	var request []midtrans.ItemDetail
	for _, value := range test {
		code := CheckCode(value)
		request = append(request, midtrans.ItemDetail{
			ID:    code,
			Name:  value["name"].(string),
			Price: int64(value["price"].(float64)),
			Qty:   int32(value["quantity"].(float64)),
		})
	}

	request = append(request, midtrans.ItemDetail{
		ID:    request[len(request)-1].ID,
		Name:  "Pajak",
		Price: 5000,
		Qty:   1,
	})

	orderID := utils.RandomString(20)
	var snapRequest midtrans.SnapReq
	snapRequest.TransactionDetails.OrderID = orderID
	snapRequest.TransactionDetails.GrossAmt = amount
	snapRequest.Items = &request
	snapRequest.CustomerDetail = &midtrans.CustDetail{
		FName: "Widdy Arfiansyah",
		Email: "widdyarfiansyah@ummi.ac.id",
	}
	snapRequest.CustomField1 = utils.IntToStr(userId)

	// Save to database
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	paymentRequest := entity.Payment{
		UserId:            utils.IntToStr(userId),
		OrderId:           orderID,
		TransactionStatus: "canceled",
		TransactionTime:   time.Now().Format("2006-01-02 15:04:05"),
	}
	payment, err := service.PaymentRepository.Create(ctx, tx, paymentRequest)
	if err != nil {
		return nil, err
	}

	// Send id payload
	snapRequest.CustomField2 = utils.IntToStr(payment.IdPayload)
	snapRequest.CustomField3 = customerName

	token, err := service.SnapGateway.GetToken(&snapRequest)
	if err != nil {
		return nil, err
	}

	var data = make(map[string]interface{})
	data["clientKey"] = service.ClientKey
	data["token"] = token.Token

	for _, value := range test {
		code := CheckCode(value)
		itemRequest := entity.UserOrder{
			PayloadId:      uint64(payment.IdPayload),
			Code:           code,
			Name:           value["name"].(string),
			Price:          int64(value["price"].(float64)),
			Image:          value["image"].(string),
			Quantity:       int(value["quantity"].(float64)),
			TotalPriceItem: int64(value["totalPricePerItem"].(float64)),
		}
		err := service.UserOrderRepository.Create(ctx, tx, itemRequest)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func CheckCode(value map[string]interface{}) string {
	checkCode := reflect.ValueOf(value["code"]).Kind()
	var code string
	if checkCode == reflect.Float64 {
		code = utils.IntToStr(int(value["code"].(float64)))
	} else if checkCode == reflect.String {
		code = value["code"].(string)
	}

	return code
}
