package rest

import (
	"errors"
	"github.com/arvians-id/apriori/helper"
	response2 "github.com/arvians-id/apriori/http/controller/rest/response"
	"github.com/arvians-id/apriori/http/middleware"
	"github.com/arvians-id/apriori/service"
	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	NotificationService service.NotificationService
}

func NewNotificationController(notificationService *service.NotificationService) *NotificationController {
	return &NotificationController{
		NotificationService: *notificationService,
	}
}

func (controller *NotificationController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.GET("/notifications", controller.FindAll)
		authorized.GET("/notifications/user", controller.FindAllByUserId)
		authorized.PATCH("/notifications/mark", controller.MarkAll)
		authorized.PATCH("/notifications/mark/:id", controller.Mark)
	}

	return router
}

func (controller *NotificationController) FindAll(c *gin.Context) {
	notifications, err := controller.NotificationService.FindAll(c.Request.Context())
	if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response2.ReturnSuccessOK(c, "OK", notifications)
}

func (controller *NotificationController) FindAllByUserId(c *gin.Context) {
	id, isExist := c.Get("id_user")
	if !isExist {
		response2.ReturnErrorUnauthorized(c, errors.New("unauthorized"), nil)
		return
	}

	notifications, err := controller.NotificationService.FindAllByUserId(c.Request.Context(), int(id.(float64)))
	if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response2.ReturnSuccessOK(c, "OK", notifications)
}

func (controller *NotificationController) MarkAll(c *gin.Context) {
	id, isExist := c.Get("id_user")
	if !isExist {
		response2.ReturnErrorUnauthorized(c, errors.New("unauthorized"), nil)
		return
	}

	err := controller.NotificationService.MarkAll(c.Request.Context(), int(id.(float64)))
	if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response2.ReturnSuccessOK(c, "OK", nil)
}

func (controller *NotificationController) Mark(c *gin.Context) {
	idParam := helper.StrToInt(c.Param("id"))
	err := controller.NotificationService.Mark(c.Request.Context(), idParam)
	if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response2.ReturnSuccessOK(c, "OK", nil)
}
