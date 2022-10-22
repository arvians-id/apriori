package rest

import (
	"errors"
	"github.com/arvians-id/apriori/internal/http/middleware"
	"github.com/arvians-id/apriori/internal/http/presenter/response"
	"github.com/arvians-id/apriori/internal/service"
	"github.com/arvians-id/apriori/util"
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
		authorized.GET("/user-order/user", controller.FindAllByUserId)
		authorized.GET("/user-order/:order_id", controller.FindAllById)
		authorized.GET("/user-order/single/:id", controller.FindById)
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

func (controller *UserOrderController) FindAllByUserId(c *gin.Context) {
	id, isExist := c.Get("id_user")
	if !isExist {
		response.ReturnErrorUnauthorized(c, errors.New("unauthorized"), nil)
		return
	}

	userOrders, err := controller.UserOrderService.FindAllByUserId(c.Request.Context(), int(id.(float64)))
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", userOrders)
}

func (controller *UserOrderController) FindAllById(c *gin.Context) {
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

	userOrder, err := controller.UserOrderService.FindAllByPayloadId(c.Request.Context(), payment.IdPayload)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", userOrder)
}

func (controller *UserOrderController) FindById(c *gin.Context) {
	IdParam := util.StrToInt(c.Param("id"))
	userOrder, err := controller.UserOrderService.FindById(c.Request.Context(), IdParam)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", userOrder)
}
