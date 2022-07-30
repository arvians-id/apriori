package controller

import (
	"apriori/api/middleware"
	"apriori/api/response"
	"apriori/cache"
	"apriori/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type UserOrderController struct {
	PaymentService   service.PaymentService
	UserOrderService service.UserOrderService
	UserOrderCache   cache.UserOrderCache
	PaymentCache     cache.PaymentCache
}

func NewUserOrderController(paymentService *service.PaymentService, UserOrderService *service.UserOrderService, userOrderCache *cache.UserOrderCache, PaymentCache *cache.PaymentCache) *UserOrderController {
	return &UserOrderController{
		PaymentService:   *paymentService,
		UserOrderService: *UserOrderService,
		UserOrderCache:   *userOrderCache,
		PaymentCache:     *PaymentCache,
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

	key := fmt.Sprintf("user-order-payment-%v", int(id.(float64)))
	paymentsCache, err := controller.PaymentCache.Get(c, key)
	if err == redis.Nil {
		payments, err := controller.PaymentService.FindAllByUserId(c.Request.Context(), int(id.(float64)))
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		err = controller.PaymentCache.Set(c.Request.Context(), key, payments)
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		response.ReturnSuccessOK(c, "OK", payments)
		return
	} else if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", paymentsCache)
}

func (controller *UserOrderController) FindById(c *gin.Context) {
	orderId := c.Param("order_id")

	key := fmt.Sprintf("user-order-id-%v", orderId)
	userOrdersCache, err := controller.UserOrderCache.Get(c, key)
	if err == redis.Nil {
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

		err = controller.UserOrderCache.Set(c.Request.Context(), key, userOrder)
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		response.ReturnSuccessOK(c, "OK", userOrder)
		return
	} else if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", userOrdersCache)
}
