package rest

import (
	"encoding/json"
	"github.com/arvians-id/apriori/cmd/library/messaging"
	"github.com/arvians-id/apriori/internal/http/middleware"
	"github.com/arvians-id/apriori/internal/http/presenter/request"
	"github.com/arvians-id/apriori/internal/http/presenter/response"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/arvians-id/apriori/internal/service"
	"github.com/arvians-id/apriori/util"
	"github.com/gin-gonic/gin"
	"github.com/veritrans/go-midtrans"
	"log"
)

type PaymentController struct {
	PaymentService      service.PaymentService
	UserOrderService    service.UserOrderService
	NotificationService service.NotificationService
	UserService         service.UserService
	Producer            messaging.Producer
}

func NewPaymentController(
	paymentService *service.PaymentService,
	userOrderService *service.UserOrderService,
	notificationService *service.NotificationService,
	userService *service.UserService,
	producer *messaging.Producer,
) *PaymentController {
	return &PaymentController{
		PaymentService:      *paymentService,
		UserOrderService:    *userOrderService,
		NotificationService: *notificationService,
		UserService:         *userService,
		Producer:            *producer,
	}
}

func (controller *PaymentController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.GET("/payments", middleware.SetupXApiKeyMiddleware(), controller.FindAll)
		authorized.GET("/payments/:order_id", middleware.SetupXApiKeyMiddleware(), controller.FindByOrderId)
		authorized.PATCH("/payments/:order_id", middleware.SetupXApiKeyMiddleware(), controller.UpdateReceiptNumber)
	}

	unauthorized := router.Group("/api")
	{
		unauthorized.POST("/payments/pay", middleware.SetupXApiKeyMiddleware(), controller.Pay)
		unauthorized.POST("/payments/notification", controller.Notification)
		unauthorized.DELETE("/payments/:order_id", middleware.SetupXApiKeyMiddleware(), controller.Delete)
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

	// Send Notification
	notification, err := controller.NotificationService.Create(c.Request.Context(), &request.CreateNotificationRequest{
		UserId:      payment.UserId,
		Title:       "Receipt number arrived",
		Description: "Your receipt number had been entered by admin",
		URL:         "product",
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	user, err := controller.UserService.FindById(c.Request.Context(), payment.UserId)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// Send Email
	emailService := model.EmailService{
		ToEmail: user.Email,
		Subject: notification.Title,
		Message: *notification.Description,
	}
	err = controller.Producer.Publish("mail_topic", emailService)
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

	response.ReturnSuccessOK(c, "OK", data)
}

func (controller *PaymentController) Notification(c *gin.Context) {
	var payload midtrans.ChargeReqWithMap
	err := c.BindJSON(&payload)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	encode, err := json.Marshal(payload)
	if err != nil {
		log.Println("[PaymentController][Notification] unable to marshal json, err: ", err.Error())
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	resArray := make(map[string]interface{})
	err = json.Unmarshal(encode, &resArray)
	if err != nil {
		log.Println("[PaymentController][Notification] unable to unmarshal json, err: ", err.Error())
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	isSettlement, err := controller.PaymentService.CreateOrUpdate(c.Request.Context(), resArray)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	if isSettlement {
		idUser := util.StrToInt(resArray["custom_field1"].(string))
		// Send Notification
		notification, err := controller.NotificationService.Create(c.Request.Context(), &request.CreateNotificationRequest{
			UserId:      idUser,
			Title:       "Transaction Successfully",
			Description: "You have successfully made a payment. Thank you for shopping at Ryzy Shop",
			URL:         "product",
		})
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		user, err := controller.UserService.FindById(c.Request.Context(), idUser)
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		// Send Email
		emailService := model.EmailService{
			ToEmail: user.Email,
			Subject: notification.Title,
			Message: *notification.Description,
		}
		err = controller.Producer.Publish("mail_topic", emailService)
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}
	}

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
