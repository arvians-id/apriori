package controller

import (
	"apriori/app/middleware"
	"apriori/app/response"
	"apriori/helper"
	"apriori/service"
	"errors"
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
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", notifications)
}

func (controller *NotificationController) FindAllByUserId(c *gin.Context) {
	id, isExist := c.Get("id_user")
	if !isExist {
		response.ReturnErrorUnauthorized(c, errors.New("unauthorized"), nil)
		return
	}

	notifications, err := controller.NotificationService.FindAllByUserId(c.Request.Context(), int(id.(float64)))
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", notifications)
}

func (controller *NotificationController) MarkAll(c *gin.Context) {
	id, isExist := c.Get("id_user")
	if !isExist {
		response.ReturnErrorUnauthorized(c, errors.New("unauthorized"), nil)
		return
	}

	err := controller.NotificationService.MarkAll(c.Request.Context(), int(id.(float64)))
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)
}

func (controller *NotificationController) Mark(c *gin.Context) {
	id := helper.StrToInt(c.Param("id"))
	err := controller.NotificationService.Mark(c.Request.Context(), id)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)
}
