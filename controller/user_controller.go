package controller

import (
	"apriori/helper"
	"apriori/middleware"
	"apriori/model"
	"apriori/service"
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
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
		authorized.GET("/users/:userId", controller.FindById)
		authorized.POST("/users", controller.Create)
		authorized.PATCH("/users/:userId", controller.Update)
		authorized.DELETE("/users/:userId", controller.Delete)

		authorized.GET("/profile", controller.Profile)
	}

	return router
}

func (controller *UserController) Profile(c *gin.Context) {
	id, isExist := c.Get("id_user")
	if !isExist {
		helper.ReturnErrorUnauthorized(c, errors.New("unauthorized"), nil)
		return
	}

	user, err := controller.UserService.FindById(c.Request.Context(), uint64(id.(float64)))
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "OK", user)
}

func (controller *UserController) FindAll(c *gin.Context) {
	users, err := controller.UserService.FindAll(c.Request.Context())
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "OK", users)
}

func (controller *UserController) FindById(c *gin.Context) {
	params := c.Param("userId")
	id, err := strconv.Atoi(params)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	user, err := controller.UserService.FindById(c.Request.Context(), uint64(id))
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "OK", user)
}

func (controller *UserController) Create(c *gin.Context) {
	var request model.CreateUserRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		helper.ReturnErrorBadRequest(c, err, nil)
		return
	}

	err = controller.UserService.Create(c.Request.Context(), request)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "created", nil)
}

func (controller *UserController) Update(c *gin.Context) {
	var request model.UpdateUserRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		helper.ReturnErrorBadRequest(c, err, nil)
		return
	}

	params := c.Param("userId")
	id, err := strconv.Atoi(params)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	request.IdUser = uint64(id)

	err = controller.UserService.Update(c.Request.Context(), request)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "updated", nil)
}

func (controller *UserController) Delete(c *gin.Context) {
	params := c.Param("userId")
	id, err := strconv.Atoi(params)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	err = controller.UserService.Delete(c.Request.Context(), uint64(id))
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "deleted", nil)
}
