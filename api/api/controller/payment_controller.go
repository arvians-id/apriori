package controller

import (
	"apriori/api/middleware"
	"apriori/api/response"
	"apriori/service"
	"apriori/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/veritrans/go-midtrans"
)

type PaymentController struct {
	PaymentService   service.PaymentService
	UserOrderService service.UserOrderService
	EmailService     service.EmailService
	CacheService     service.CacheService
}

func NewPaymentController(paymentService *service.PaymentService, userOrderService *service.UserOrderService, emailService service.EmailService, cacheService *service.CacheService) *PaymentController {
	return &PaymentController{
		PaymentService:   *paymentService,
		UserOrderService: *userOrderService,
		EmailService:     emailService,
		CacheService:     *cacheService,
	}
}

func (controller *PaymentController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.GET("/payments", controller.FindAll)
		authorized.GET("/payments/:order_id", controller.FindByOrderId)
	}

	notAuthorized := router.Group("/api")
	{
		notAuthorized.POST("/payments/pay", controller.Pay)
		notAuthorized.POST("/payments/notification", controller.Notification)
		notAuthorized.DELETE("/payments/:order_id", controller.Delete)
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
	orderId := c.Param("order_id")

	payment, err := controller.PaymentService.FindByOrderId(c.Request.Context(), orderId)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", payment)
}

func (controller *PaymentController) Pay(c *gin.Context) {
	grossAmount := int64(utils.StrToInt(c.PostForm("gross_amount")))
	items := c.PostFormArray("items")
	userId := utils.StrToInt(c.PostForm("user_id"))
	customerName := c.PostForm("customer_name")

	data, err := controller.PaymentService.GetToken(c.Request.Context(), grossAmount, userId, customerName, items)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	// delete previous cache
	key := fmt.Sprintf("user-order-payment-%v", userId)
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

	//Save to database
	err = controller.PaymentService.CreateOrUpdate(c.Request.Context(), resArray)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	key := fmt.Sprintf("user-order-id-%v", utils.StrToInt(resArray["custom_field2"].(string)))
	key2 := fmt.Sprintf("user-order-payment-%v", utils.StrToInt(resArray["custom_field1"].(string)))
	_ = controller.CacheService.Del(c.Request.Context(), key, key2)

	response.ReturnSuccessOK(c, "OK", nil)
}

func (controller *PaymentController) Delete(c *gin.Context) {
	orderId := c.Param("order_id")

	err := controller.PaymentService.Delete(c.Request.Context(), orderId)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)
}
