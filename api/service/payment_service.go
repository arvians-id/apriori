package service

import (
	"apriori/config"
	"apriori/utils"
	"github.com/veritrans/go-midtrans"
)

type PaymentService interface {
	GetClient()
	GetToken(amount int64) (map[string]interface{}, error)
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

func (service *paymentService) GetToken(amount int64) (map[string]interface{}, error) {
	service.GetClient()

	orderID := utils.RandomString(20)
	token, err := service.SnapGateway.GetTokenQuick(orderID, amount)
	if err != nil {
		return nil, err
	}

	var data = make(map[string]interface{})
	data["clientKey"] = service.ClientKey
	data["token"] = token.Token

	return data, nil
}
