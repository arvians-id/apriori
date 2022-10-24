package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/arvians-id/apriori/cmd/config"
	request2 "github.com/arvians-id/apriori/internal/http/presenter/request"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/arvians-id/apriori/internal/repository"
	"github.com/arvians-id/apriori/util"
	"github.com/veritrans/go-midtrans"
	"log"
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

func (service *PaymentServiceImpl) FindAll(ctx context.Context) ([]*model.Payment, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[PaymentService][FindAll] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	payments, err := service.PaymentRepository.FindAll(ctx, tx)
	if err != nil {
		log.Println("[PaymentService][FindAll][FindAll] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return payments, nil
}

func (service *PaymentServiceImpl) FindAllByUserId(ctx context.Context, userId int) ([]*model.Payment, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[PaymentService][FindAllByUserId] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	payments, err := service.PaymentRepository.FindAllByUserId(ctx, tx, userId)
	if err != nil {
		log.Println("[PaymentService][FindAllByUserId][FindAllByUserId] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return payments, nil
}

func (service *PaymentServiceImpl) FindByOrderId(ctx context.Context, orderId string) (*model.Payment, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[PaymentService][FindByOrderId] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	payment, err := service.PaymentRepository.FindByOrderId(ctx, tx, orderId)
	if err != nil {
		log.Println("[PaymentService][FindByOrderId][FindByOrderId] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return payment, nil
}

func (service *PaymentServiceImpl) CreateOrUpdate(ctx context.Context, requestPayment map[string]interface{}) error {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[PaymentService][CreateOrUpdate] problem in db transaction, err: ", err.Error())
		return err
	}
	defer util.CommitOrRollback(tx)

	var bankType, vaNumber, billerCode, billKey, settlementTime string

	if requestPayment["va_numbers"] != nil {
		bankType = requestPayment["va_numbers"].([]interface{})[0].(map[string]interface{})["bank"].(string)
		vaNumber = requestPayment["va_numbers"].([]interface{})[0].(map[string]interface{})["va_number"].(string)
	} else if requestPayment["permata_va_number"] != nil {
		bankType = "permata bank"
		vaNumber = requestPayment["permata_va_number"].(string)
	} else if requestPayment["biller_code"] != nil && requestPayment["bill_key"] != nil {
		billerCode = requestPayment["biller_code"].(string)
		billKey = requestPayment["bill_key"].(string)
		bankType = "mandiri"
	}

	setTime, ok := requestPayment["settlement_time"]
	if ok {
		settlementTime = setTime.(string)
	} else {
		settlementTime = ""
	}

	orderID := requestPayment["order_id"].(string)
	transactionTime := requestPayment["transaction_time"].(string)
	transactionStatus := requestPayment["transaction_status"].(string)
	transactionId := requestPayment["transaction_id"].(string)
	statusCode := requestPayment["status_code"].(string)
	signatureKey := requestPayment["signature_key"].(string)
	paymentType := requestPayment["payment_type"].(string)
	merchantId := requestPayment["merchant_id"].(string)
	grossAmount := requestPayment["gross_amount"].(string)
	fraudStatus := requestPayment["fraud_status"].(string)

	checkTransaction, err := service.PaymentRepository.FindByOrderId(ctx, tx, requestPayment["order_id"].(string))
	if err != nil {
		log.Println("[PaymentService][CreateOrUpdate][FindByOrderId] problem in getting from repository, err: ", err.Error())
		return err
	}

	checkTransaction.UserId = util.StrToInt(requestPayment["custom_field1"].(string))
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
			log.Println("[PaymentService][CreateOrUpdate][Update] problem in getting from repository, err: ", err.Error())
			return err
		}

		if requestPayment["transaction_status"].(string) == "settlement" {
			timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
			if err != nil {
				log.Println("[PaymentService][CreateOrUpdate] problem in parsing to time, err: ", err.Error())
				return err
			}

			userOrder, err := service.UserOrderRepository.FindAllByPayloadId(ctx, tx, requestPayment["custom_field2"].(string))
			if err != nil {
				log.Println("[PaymentService][CreateOrUpdate][FindAllByPayloadId] problem in getting from repository, err: ", err.Error())
				return err
			}

			var productName []string
			for _, item := range userOrder {
				productName = append(productName, *item.Name)
			}

			transaction := model.Transaction{
				ProductName:   strings.ToLower(strings.Join(productName, ", ")),
				CustomerName:  requestPayment["custom_field3"].(string),
				NoTransaction: requestPayment["order_id"].(string),
				CreatedAt:     timeNow,
				UpdatedAt:     timeNow,
			}

			_, err = service.TransactionRepository.Create(ctx, tx, &transaction)
			if err != nil {
				log.Println("[PaymentService][CreateOrUpdate][TransactionCreate] problem in getting from repository, err: ", err.Error())
				return err
			}

			var notificationRequest request2.CreateNotificationRequest
			notificationRequest.UserId = checkTransaction.UserId
			notificationRequest.Title = "Transaction Successfully"
			notificationRequest.Description = "You have successfully made a payment. Thank you for shopping at Ryzy Shop"
			notificationRequest.URL = "product"
			err = service.NotificationService.Create(ctx, &notificationRequest).WithSendMail()
			if err != nil {
				log.Println("[PaymentService][CreateOrUpdate][NotificationCreate] problem in getting from repository, err: ", err.Error())
				return err
			}
		}
	}

	return nil
}

func (service *PaymentServiceImpl) UpdateReceiptNumber(ctx context.Context, request *request2.AddReceiptNumberRequest) (*model.Payment, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[PaymentService][UpdateReceiptNumber] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	payment, err := service.PaymentRepository.FindByOrderId(ctx, tx, request.OrderId)
	if err != nil {
		log.Println("[PaymentService][UpdateReceiptNumber][FindByOrderId] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	payment.ReceiptNumber = &request.ReceiptNumber

	err = service.PaymentRepository.UpdateReceiptNumber(ctx, tx, payment)
	if err != nil {
		log.Println("[PaymentService][UpdateReceiptNumber][UpdateReceiptNumber] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return payment, nil
}

func (service *PaymentServiceImpl) Delete(ctx context.Context, orderId string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[PaymentService][Delete] problem in db transaction, err: ", err.Error())
		return err
	}
	defer util.CommitOrRollback(tx)

	payment, err := service.PaymentRepository.FindByOrderId(ctx, tx, orderId)
	if err != nil {
		log.Println("[PaymentService][Delete][FindByOrderId] problem in getting from repository, err: ", err.Error())
		return err
	}

	err = service.PaymentRepository.Delete(ctx, tx, payment.OrderId)
	if err != nil {
		log.Println("[PaymentService][Delete][FindByOrderId] problem in getting from repository, err: ", err.Error())
		return err
	}

	return nil
}
func (service *PaymentServiceImpl) GetToken(ctx context.Context, request *request2.GetPaymentTokenRequest) (map[string]interface{}, error) {
	service.GetClient()

	var items []map[string]interface{}
	for _, item := range request.Items {
		err := json.Unmarshal([]byte(item), &items)
		if err != nil {
			log.Println("[PaymentService][GetToken] unable to unmarshal json, err: ", err.Error())
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

	orderID := util.RandomString(20)
	var snapRequest midtrans.SnapReq
	snapRequest.TransactionDetails.OrderID = orderID
	snapRequest.TransactionDetails.GrossAmt = request.GrossAmount
	snapRequest.Items = &itemDetails
	snapRequest.CustomerDetail = &midtrans.CustDetail{
		FName: request.CustomerName,
	}
	snapRequest.CustomField1 = util.IntToStr(request.UserId)

	// Save to database
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[PaymentService][GetToken] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	canceled := "canceled"
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	paymentRequest := model.Payment{
		UserId:            request.UserId,
		OrderId:           &orderID,
		TransactionStatus: &canceled,
		TransactionTime:   &timeNow,
		Address:           &request.Address,
		Courier:           &request.Courier,
		CourierService:    &request.CourierService,
	}
	payment, err := service.PaymentRepository.Create(ctx, tx, &paymentRequest)
	if err != nil {
		log.Println("[PaymentService][GetToken][Create] problem in getting from repository, err: ", err.Error())
		return nil, err
	}
	// Send id payload
	snapRequest.CustomField2 = util.IntToStr(payment.IdPayload)
	snapRequest.CustomField3 = request.CustomerName

	token, err := service.SnapGateway.GetToken(&snapRequest)
	if err != nil {
		log.Println("[PaymentService][GetToken][GetToken] problem in getting from repository, err: ", err.Error())
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
		userOrder := model.UserOrder{
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
			log.Println("[PaymentService][GetToken][Create] problem in getting from repository, err: ", err.Error())
			return nil, err
		}
	}

	return tokenResponse, nil
}

func CheckCode(value map[string]interface{}) string {
	checkCode := reflect.ValueOf(value["code"]).Kind()
	var code string
	if checkCode == reflect.Float64 {
		code = util.IntToStr(int(value["code"].(float64)))
	} else if checkCode == reflect.String {
		code = value["code"].(string)
	}

	return code
}
