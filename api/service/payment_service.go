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
	DB                    *sql.DB
}

func NewPaymentService(
	configuration config.Config,
	paymentRepository *repository.PaymentRepository,
	userOrderRepository *repository.UserOrderRepository,
	transactionRepository *repository.TransactionRepository,
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
		DB:                    db,
	}
}

func (service *PaymentServiceImpl) GetClient() {
	service.SnapGateway = midtrans.SnapGateway{
		Client: service.MidClient,
	}
}

func (service *PaymentServiceImpl) FindAll(ctx context.Context) ([]*model.GetPaymentRelationResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	payments, err := service.PaymentRepository.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	var paymentResponses []*model.GetPaymentRelationResponse
	for _, payment := range payments {
		paymentResponses = append(paymentResponses, payment.ToPaymentRelationResponse())
	}

	return paymentResponses, nil
}

func (service *PaymentServiceImpl) FindAllByUserId(ctx context.Context, userId int) ([]*model.GetPaymentResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	payments, err := service.PaymentRepository.FindAllByUserId(ctx, tx, userId)
	if err != nil {
		return nil, err
	}

	var paymentResponses []*model.GetPaymentResponse
	for _, payment := range payments {
		paymentResponses = append(paymentResponses, payment.ToPaymentResponse())
	}

	return paymentResponses, nil
}

func (service *PaymentServiceImpl) FindByOrderId(ctx context.Context, orderId string) (*model.GetPaymentResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return &model.GetPaymentResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	paymentResponse, err := service.PaymentRepository.FindByOrderId(ctx, tx, orderId)
	if err != nil {
		return &model.GetPaymentResponse{}, err
	}

	return paymentResponse.ToPaymentResponse(), nil
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

	paymentRequest := entity.Payment{
		UserId: sql.NullString{
			String: request["custom_field1"].(string),
			Valid:  true,
		},
		OrderId: sql.NullString{
			String: request["order_id"].(string),
			Valid:  true,
		},
		TransactionTime: sql.NullString{
			String: request["transaction_time"].(string),
			Valid:  true,
		},
		TransactionStatus: sql.NullString{
			String: request["transaction_status"].(string),
			Valid:  true,
		},
		TransactionId: sql.NullString{
			String: request["transaction_id"].(string),
			Valid:  true,
		},
		StatusCode: sql.NullString{
			String: request["status_code"].(string),
			Valid:  true,
		},
		SignatureKey: sql.NullString{
			String: request["signature_key"].(string),
			Valid:  true,
		},
		SettlementTime: sql.NullString{
			String: settlementTime,
			Valid:  true,
		},
		PaymentType: sql.NullString{
			String: request["payment_type"].(string),
			Valid:  true,
		},
		MerchantId: sql.NullString{
			String: request["merchant_id"].(string),
			Valid:  true,
		},
		GrossAmount: sql.NullString{
			String: request["gross_amount"].(string),
			Valid:  true,
		},
		FraudStatus: sql.NullString{
			String: request["fraud_status"].(string),
			Valid:  true,
		},
		BankType: sql.NullString{
			String: bankType,
			Valid:  true,
		},
		VANumber: sql.NullString{
			String: vaNumber,
			Valid:  true,
		},
		BillerCode: sql.NullString{
			String: billerCode,
			Valid:  true,
		},
		BillKey: sql.NullString{
			String: billKey,
			Valid:  true,
		},
	}

	checkTransaction, _ := service.PaymentRepository.FindByOrderId(ctx, tx, request["order_id"].(string))
	if checkTransaction.OrderId.Valid {
		err := service.PaymentRepository.Update(ctx, tx, &paymentRequest)
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
				productName = append(productName, item.Name)
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
		}
	}

	return nil
}

func (service *PaymentServiceImpl) UpdateReceiptNumber(ctx context.Context, request *model.AddReceiptNumberRequest) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	payment, err := service.PaymentRepository.FindByOrderId(ctx, tx, request.OrderId)
	if err != nil {
		return err
	}

	paymentRequest := entity.Payment{
		OrderId: payment.OrderId,
		ReceiptNumber: sql.NullString{
			String: request.ReceiptNumber,
			Valid:  true,
		},
	}
	err = service.PaymentRepository.UpdateReceiptNumber(ctx, tx, &paymentRequest)
	if err != nil {
		return err
	}

	return nil
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

	err = service.PaymentRepository.Delete(ctx, tx, payment.OrderId.String)
	if err != nil {
		return err
	}

	return nil
}
func (service *PaymentServiceImpl) GetToken(ctx context.Context, amount int64, userId int, customerName string, items []string, rajaShipping *model.GetRajaOngkirResponse) (map[string]interface{}, error) {
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

	request = append(request, midtrans.ItemDetail{
		ID:    request[len(request)-1].ID,
		Name:  "Ongkos Kirim",
		Price: rajaShipping.ShippingCost,
		Qty:   1,
	})

	orderID := helper.RandomString(20)
	var snapRequest midtrans.SnapReq
	snapRequest.TransactionDetails.OrderID = orderID
	snapRequest.TransactionDetails.GrossAmt = amount
	snapRequest.Items = &request
	snapRequest.CustomerDetail = &midtrans.CustDetail{
		FName: "Widdy Arfiansyah",
		Email: "widdyarfiansyah@ummi.ac.id",
	}
	snapRequest.CustomField1 = helper.IntToStr(userId)

	// Save to database
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	paymentRequest := entity.Payment{
		UserId: sql.NullString{
			String: helper.IntToStr(userId),
			Valid:  true,
		},
		OrderId: sql.NullString{
			String: orderID,
			Valid:  true,
		},
		TransactionStatus: sql.NullString{
			String: "canceled",
			Valid:  true,
		},
		TransactionTime: sql.NullString{
			String: time.Now().Format("2006-01-02 15:04:05"),
			Valid:  true,
		},
		Address: sql.NullString{
			String: rajaShipping.Address,
			Valid:  true,
		},
		Courier: sql.NullString{
			String: rajaShipping.Courier,
			Valid:  true,
		},
		CourierService: sql.NullString{
			String: rajaShipping.CourierService,
			Valid:  true,
		},
	}
	payment, err := service.PaymentRepository.Create(ctx, tx, &paymentRequest)
	if err != nil {
		return nil, err
	}

	// Send id payload
	snapRequest.CustomField2 = helper.IntToStr(payment.IdPayload)
	snapRequest.CustomField3 = customerName

	token, err := service.SnapGateway.GetToken(&snapRequest)
	if err != nil {
		return nil, err
	}

	var tokenResponse = make(map[string]interface{})
	tokenResponse["clientKey"] = service.ClientKey
	tokenResponse["token"] = token.Token

	for _, value := range test {
		code := CheckCode(value)
		itemRequest := entity.UserOrder{
			PayloadId:      payment.IdPayload,
			Code:           code,
			Name:           value["name"].(string),
			Price:          int64(value["price"].(float64)),
			Image:          value["image"].(string),
			Quantity:       int(value["quantity"].(float64)),
			TotalPriceItem: int64(value["totalPricePerItem"].(float64)),
		}
		err := service.UserOrderRepository.Create(ctx, tx, &itemRequest)
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
