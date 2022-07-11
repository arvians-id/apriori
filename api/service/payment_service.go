package service

import (
	"apriori/config"
	"apriori/utils"
	"encoding/json"
	"github.com/veritrans/go-midtrans"
	"reflect"
)

type PaymentService interface {
	GetClient()
	GetToken(amount int64, userId int, items []string) (map[string]interface{}, error)
}

type paymentService struct {
	MidClient   midtrans.Client
	CoreGateway midtrans.CoreGateway
	SnapGateway midtrans.SnapGateway
	ServerKey   string
	ClientKey   string
}

func NewPaymentService(configuration config.Config) PaymentService {
	midClient := midtrans.NewClient()
	midClient.ServerKey = configuration.Get("MIDTRANS_SERVER_KEY")
	midClient.ClientKey = configuration.Get("MIDTRANS_CLIENT_KEY")
	midClient.APIEnvType = midtrans.Sandbox

	return &paymentService{
		MidClient: midClient,
		ServerKey: midClient.ServerKey,
		ClientKey: midClient.ClientKey,
	}
}

func (service *paymentService) GetClient() {
	service.CoreGateway = midtrans.CoreGateway{
		Client: service.MidClient,
	}
	service.SnapGateway = midtrans.SnapGateway{
		Client: service.MidClient,
	}
}

func (service *paymentService) GetToken(amount int64, userId int, items []string) (map[string]interface{}, error) {
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
		checkCode := reflect.ValueOf(value["code"]).Kind()
		var code string
		if checkCode == reflect.Float64 {
			code = utils.IntToStr(int(value["code"].(float64)))
		} else if checkCode == reflect.String {
			code = value["code"].(string)
		}
		request = append(request, midtrans.ItemDetail{
			ID:    code,
			Name:  value["name"].(string),
			Price: int64(value["price"].(float64)),
			Qty:   int32(value["quantity"].(float64)),
		})
	}

	orderID := utils.RandomString(20)
	var snapRequest midtrans.SnapReq
	snapRequest.TransactionDetails.OrderID = orderID
	snapRequest.TransactionDetails.GrossAmt = amount
	snapRequest.Items = &request
	snapRequest.CustomerDetail = &midtrans.CustDetail{
		FName: "Widdy Arfiansyah",
		Email: "widdyarfiansyah@ummi.ac.id",
	}

	token, err := service.SnapGateway.GetToken(&snapRequest)
	if err != nil {
		return nil, err
	}

	var data = make(map[string]interface{})
	data["clientKey"] = service.ClientKey
	data["token"] = token.Token

	return data, nil
}
