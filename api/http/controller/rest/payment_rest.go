package rest

import (
	"encoding/json"
	"fmt"
	"github.com/arvians-id/apriori/helper"
	"github.com/arvians-id/apriori/http/middleware"
	"github.com/arvians-id/apriori/http/request"
	"github.com/arvians-id/apriori/http/response"
	"github.com/arvians-id/apriori/service"
	"github.com/gin-gonic/gin"
	"github.com/veritrans/go-midtrans"
)

type PaymentController struct {
	PaymentService      service.PaymentService
	UserOrderService    service.UserOrderService
	EmailService        service.EmailService
	CacheService        service.CacheService
	NotificationService service.NotificationService
}

func NewPaymentController(
	paymentService *service.PaymentService,
	userOrderService *service.UserOrderService,
	emailService *service.EmailService,
	cacheService *service.CacheService,
	notificationService *service.NotificationService,
) *PaymentController {
	return &PaymentController{
		PaymentService:      *paymentService,
		UserOrderService:    *userOrderService,
		EmailService:        *emailService,
		CacheService:        *cacheService,
		NotificationService: *notificationService,
	}
}

func (controller *PaymentController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.GET("/payments", controller.FindAll)
		authorized.GET("/payments/:order_id", controller.FindByOrderId)
		authorized.PATCH("/payments/:order_id", controller.UpdateReceiptNumber)
	}

	unauthorized := router.Group("/api")
	{
		unauthorized.POST("/payments/pay", controller.Pay)
		unauthorized.POST("/payments/notification", controller.Notification)
		unauthorized.DELETE("/payments/:order_id", controller.Delete)
	}

	return router
}

func (controller *PaymentController) FindAll(c *gin.Context) {
	payments, err := controller.PaymentService.FindAll(c.Request.Context())
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", payments)
}

func (controller *PaymentController) FindByOrderId(c *gin.Context) {
	orderIdParam := c.Param("order_id")

	payment, err := controller.PaymentService.FindByOrderId(c.Request.Context(), orderIdParam)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", payment)
}

func (controller *PaymentController) UpdateReceiptNumber(c *gin.Context) {
	var requestPayment request.AddReceiptNumberRequest
	err := c.ShouldBindJSON(&requestPayment)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	requestPayment.OrderId = c.Param("order_id")
	payment, err := controller.PaymentService.UpdateReceiptNumber(c.Request.Context(), &requestPayment)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// Notification
	var notificationRequest request.CreateNotificationRequest
	notificationRequest.UserId = payment.UserId
	notificationRequest.Title = "Receipt number arrived"
	notificationRequest.Description = "Your receipt number has been entered by the admin"
	notificationRequest.URL = "product"
	err = controller.NotificationService.Create(c.Request.Context(), &notificationRequest).WithSendMail()
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)
}

func (controller *PaymentController) Pay(c *gin.Context) {
	var requestToken request.GetPaymentTokenRequest
	err := c.ShouldBind(&requestToken)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	data, err := controller.PaymentService.GetToken(c.Request.Context(), &requestToken)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	key := fmt.Sprintf("user-order-payment-%v", requestToken.UserId)
	_ = controller.CacheService.Del(c.Request.Context(), key)

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

	err = controller.PaymentService.CreateOrUpdate(c.Request.Context(), resArray)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	key := fmt.Sprintf("user-order-id-%v", helper.StrToInt(resArray["custom_field2"].(string)))
	key2 := fmt.Sprintf("user-order-payment-%v", helper.StrToInt(resArray["custom_field1"].(string)))
	key3 := fmt.Sprintf("user-order-rate-%v", helper.StrToInt(resArray["custom_field1"].(string)))
	_ = controller.CacheService.Del(c.Request.Context(), key, key2, key3)

	response.ReturnSuccessOK(c, "OK", nil)
}

func (controller *PaymentController) Delete(c *gin.Context) {
	orderIdParam := c.Param("order_id")
	err := controller.PaymentService.Delete(c.Request.Context(), orderIdParam)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)
}
