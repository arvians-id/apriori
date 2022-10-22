package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/arvians-id/apriori/internal/http/middleware"
	response2 "github.com/arvians-id/apriori/internal/http/presenter/response"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/arvians-id/apriori/internal/service"
	"github.com/arvians-id/apriori/util"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type UserOrderController struct {
	PaymentService   service.PaymentService
	UserOrderService service.UserOrderService
	CacheService     service.CacheService
}

func NewUserOrderController(
	paymentService *service.PaymentService,
	UserOrderService *service.UserOrderService,
	cacheService *service.CacheService,
) *UserOrderController {
	return &UserOrderController{
		PaymentService:   *paymentService,
		UserOrderService: *UserOrderService,
		CacheService:     *cacheService,
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
		response2.ReturnErrorUnauthorized(c, errors.New("unauthorized"), nil)
		return
	}

	key := fmt.Sprintf("user-order-payment-%v", int(id.(float64)))
	paymentsCache, err := controller.CacheService.Get(c, key)
	if err == redis.Nil {
		payments, err := controller.PaymentService.FindAllByUserId(c.Request.Context(), int(id.(float64)))
		if err != nil {
			response2.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		err = controller.CacheService.Set(c.Request.Context(), key, payments)
		if err != nil {
			response2.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		response2.ReturnSuccessOK(c, "OK", payments)
		return
	} else if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	var paymentCacheResponses []model.Payment
	err = json.Unmarshal(paymentsCache, &paymentCacheResponses)
	if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response2.ReturnSuccessOK(c, "OK", paymentCacheResponses)
}

func (controller *UserOrderController) FindAllByUserId(c *gin.Context) {
	id, isExist := c.Get("id_user")
	if !isExist {
		response2.ReturnErrorUnauthorized(c, errors.New("unauthorized"), nil)
		return
	}

	key := fmt.Sprintf("user-order-rate-%v", int(id.(float64)))
	userOrdersCache, err := controller.CacheService.Get(c, key)
	if err == redis.Nil {
		userOrders, err := controller.UserOrderService.FindAllByUserId(c.Request.Context(), int(id.(float64)))
		if err != nil {
			response2.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		err = controller.CacheService.Set(c.Request.Context(), key, userOrders)
		if err != nil {
			response2.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		response2.ReturnSuccessOK(c, "OK", userOrders)
		return
	} else if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	var userOrderCacheResponses []model.UserOrder
	err = json.Unmarshal(userOrdersCache, &userOrderCacheResponses)
	if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response2.ReturnSuccessOK(c, "OK", userOrderCacheResponses)
}

func (controller *UserOrderController) FindAllById(c *gin.Context) {
	orderIdParam := c.Param("order_id")
	key := fmt.Sprintf("user-order-id-%v", orderIdParam)
	userOrdersCache, err := controller.CacheService.Get(c, key)
	if err == redis.Nil {
		payment, err := controller.PaymentService.FindByOrderId(c.Request.Context(), orderIdParam)
		if err != nil {
			if err.Error() == response2.ErrorNotFound {
				response2.ReturnErrorNotFound(c, err, nil)
				return
			}

			response2.ReturnErrorInternalServerError(c, err, nil)
			return
		}
		userOrder, err := controller.UserOrderService.FindAllByPayloadId(c.Request.Context(), payment.IdPayload)
		if err != nil {
			response2.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		err = controller.CacheService.Set(c.Request.Context(), key, userOrder)
		if err != nil {
			response2.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		response2.ReturnSuccessOK(c, "OK", userOrder)
		return
	} else if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	var userOrderCacheResponses []model.UserOrder
	err = json.Unmarshal(userOrdersCache, &userOrderCacheResponses)
	if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response2.ReturnSuccessOK(c, "OK", userOrderCacheResponses)
}

func (controller *UserOrderController) FindById(c *gin.Context) {
	IdParam := util.StrToInt(c.Param("id"))
	userOrder, err := controller.UserOrderService.FindById(c.Request.Context(), IdParam)
	if err != nil {
		if err.Error() == response2.ErrorNotFound {
			response2.ReturnErrorNotFound(c, err, nil)
			return
		}

		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response2.ReturnSuccessOK(c, "OK", userOrder)
}
