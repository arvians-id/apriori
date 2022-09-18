package rest

import (
	"errors"
	"github.com/arvians-id/apriori/helper"
	"github.com/arvians-id/apriori/http/middleware"
	"github.com/arvians-id/apriori/http/request"
	"github.com/arvians-id/apriori/http/response"
	"github.com/arvians-id/apriori/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		UserService: *userService,
	}
}

func (controller *UserController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.GET("/users", controller.FindAll)
		authorized.GET("/users/:id", controller.FindById)
		authorized.POST("/users", controller.Create)
		authorized.PATCH("/users/:id", controller.Update)
		authorized.DELETE("/users/:id", controller.Delete)
		authorized.GET("/profile", controller.Profile)
		authorized.PATCH("/profile/update", controller.UpdateProfile)
	}

	return router
}

func (controller *UserController) Profile(c *gin.Context) {
	id, isExist := c.Get("id_user")
	if !isExist {
		response.ReturnErrorUnauthorized(c, errors.New("unauthorized"), nil)
		return
	}

	user, err := controller.UserService.FindById(c.Request.Context(), int(id.(float64)))
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", user)
}

func (controller *UserController) UpdateProfile(c *gin.Context) {
	var requestUpdate request.UpdateUserRequest
	err := c.ShouldBindJSON(&requestUpdate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	id, isExist := c.Get("id_user")
	if !isExist {
		response.ReturnErrorUnauthorized(c, errors.New("unauthorized"), nil)
		return
	}

	requestUpdate.IdUser = int(id.(float64))
	user, err := controller.UserService.Update(c.Request.Context(), &requestUpdate)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "updated", user)
}

func (controller *UserController) FindAll(c *gin.Context) {
	users, err := controller.UserService.FindAll(c.Request.Context())
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", users)
}

func (controller *UserController) FindById(c *gin.Context) {
	idParam := helper.StrToInt(c.Param("id"))
	user, err := controller.UserService.FindById(c.Request.Context(), idParam)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", user)
}

func (controller *UserController) Create(c *gin.Context) {
	var requestCreate request.CreateUserRequest
	err := c.ShouldBindJSON(&requestCreate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	user, err := controller.UserService.Create(c.Request.Context(), &requestCreate)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "created", user)
}

func (controller *UserController) Update(c *gin.Context) {
	var requestUpdate request.UpdateUserRequest
	err := c.ShouldBindJSON(&requestUpdate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	requestUpdate.IdUser = helper.StrToInt(c.Param("id"))
	user, err := controller.UserService.Update(c.Request.Context(), &requestUpdate)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "updated", user)
}

func (controller *UserController) Delete(c *gin.Context) {
	idParam := helper.StrToInt(c.Param("id"))
	err := controller.UserService.Delete(c.Request.Context(), idParam)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "deleted", nil)
}