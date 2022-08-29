package controller

import (
	"apriori/app/middleware"
	"apriori/app/response"
	"apriori/helper"
	"apriori/model"
	"apriori/service"
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

func NewPaymentController(
	paymentService *service.PaymentService,
	userOrderService *service.UserOrderService,
	emailService service.EmailService,
	cacheService *service.CacheService,
) *PaymentController {
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
	var request model.AddReceiptNumberRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	request.OrderId = c.Param("order_id")
	err = controller.PaymentService.UpdateReceiptNumber(c.Request.Context(), &request)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)
}

func (controller *PaymentController) Pay(c *gin.Context) {
	grossAmount := int64(helper.StrToInt(c.PostForm("gross_amount")))
	items := c.PostFormArray("items")
	userId := helper.StrToInt(c.PostForm("user_id"))
	customerName := c.PostForm("customer_name")

	var rajaShipping model.GetRajaOngkirResponse
	rajaShipping.Address = c.PostForm("address")
	rajaShipping.Courier = c.PostForm("courier")
	rajaShipping.CourierService = c.PostForm("courier_service")
	rajaShipping.ShippingCost = int64(helper.StrToInt(c.PostForm("shipping_cost")))

	data, err := controller.PaymentService.GetToken(c.Request.Context(), grossAmount, userId, customerName, items, &rajaShipping)
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
