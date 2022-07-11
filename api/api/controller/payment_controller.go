package controller

import (
	"apriori/api/response"
	"apriori/service"
	"apriori/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/veritrans/go-midtrans"
)

type PaymentController struct {
	PaymentService service.PaymentService
	EmailService   service.EmailService
}

func NewPaymentController(paymentService *service.PaymentService, emailService service.EmailService) *PaymentController {
	return &PaymentController{
		PaymentService: *paymentService,
		EmailService:   emailService,
	}
}

func (controller *PaymentController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api")
	{
		authorized.POST("/pay", controller.Pay)
		authorized.POST("/notification", controller.Notification)
	}

	return router
}

func (controller *PaymentController) Pay(c *gin.Context) {
	grossAmount := int64(utils.StrToInt(c.PostForm("gross_amount")))
	items := c.PostFormArray("items")
	userId := 7
	data, err := controller.PaymentService.GetToken(grossAmount, userId, items)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", data)
}

func (controller *PaymentController) Notification(c *gin.Context) {
	var payload midtrans.ChargeReqWithMap
	err := c.BindJSON(&payload)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	encode, _ := json.Marshal(payload)
	resArray := make(map[string]interface{})
	err = json.Unmarshal(encode, &resArray)

	// Send email to user
	message := fmt.Sprintf("%s", resArray["signature_key"])
	err = controller.EmailService.SendEmailWithText("widdyarfiansyah@ummi.ac.id", message)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", resArray)
}
