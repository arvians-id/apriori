package controller

import (
	"apriori/api/middleware"
	"apriori/api/response"
	"apriori/service"
	"errors"
	"github.com/gin-gonic/gin"
)

type UserOrderController struct {
	PaymentService   service.PaymentService
	UserOrderService service.UserOrderService
}

func NewUserOrderController(paymentService *service.PaymentService, UserOrderService *service.UserOrderService) *UserOrderController {
	return &UserOrderController{
		PaymentService:   *paymentService,
		UserOrderService: *UserOrderService,
	}
}

func (controller *UserOrderController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.GET("/user-order", controller.FindAll)
		authorized.GET("/user-order/:order_id", controller.FindById)
	}

	return router
}

func (controller *UserOrderController) FindAll(c *gin.Context) {
	id, isExist := c.Get("id_user")
	if !isExist {
		response.ReturnErrorUnauthorized(c, errors.New("unauthorized"), nil)
		return
	}

	payments, err := controller.PaymentService.FindAllByUserId(c.Request.Context(), int(id.(float64)))
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", payments)
}

func (controller *UserOrderController) FindById(c *gin.Context) {
	orderId := c.Param("order_id")

	payment, err := controller.PaymentService.FindByOrderId(c.Request.Context(), orderId)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	userOrder, err := controller.UserOrderService.FindAllByPayload(c.Request.Context(), payment.IdPayload)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", userOrder)
}
