package controller

import (
	"apriori/model"
	"apriori/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type UserController struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (controller *UserController) Route(router *gin.Engine) *gin.Engine {
	router.GET("/api/users", controller.FindAll)
	router.GET("/api/users/:userId", controller.FindById)
	router.POST("/api/users", controller.Create)
	router.PATCH("/api/users/:userId", controller.Update)
	router.DELETE("/api/users/:userId", controller.Delete)

	return router
}

func (controller *UserController) FindAll(c *gin.Context) {
	users, err := controller.UserService.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(500, model.WebResponse{
			Code:   500,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	c.JSON(200, model.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   users,
	})
}

func (controller *UserController) FindById(c *gin.Context) {
	params := c.Param("userId")
	id, err := strconv.Atoi(params)
	if err != nil {
		c.JSON(500, model.WebResponse{
			Code:   500,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	user, err := controller.UserService.FindById(c.Request.Context(), uint64(id))
	if err != nil {
		c.JSON(500, model.WebResponse{
			Code:   500,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	c.JSON(200, model.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   user,
	})
}

func (controller *UserController) Create(c *gin.Context) {
	var request model.CreateUserRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(500, model.WebResponse{
			Code:   500,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	user, err := controller.UserService.Create(c.Request.Context(), request)
	if err != nil {
		c.JSON(500, model.WebResponse{
			Code:   500,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	c.JSON(200, model.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   user,
	})
}

func (controller *UserController) Update(c *gin.Context) {
	params := c.Param("userId")
	id, err := strconv.Atoi(params)
	if err != nil {
		c.JSON(500, model.WebResponse{
			Code:   500,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	var request model.UpdateUserRequest
	err = c.BindJSON(&request)
	if err != nil {
		c.JSON(500, model.WebResponse{
			Code:   500,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	request.IdUser = uint64(id)

	user, err := controller.UserService.Update(c.Request.Context(), request)
	if err != nil {
		c.JSON(500, model.WebResponse{
			Code:   500,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	c.JSON(200, model.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   user,
	})
}

func (controller *UserController) Delete(c *gin.Context) {
	params := c.Param("userId")
	id, err := strconv.Atoi(params)
	if err != nil {
		c.JSON(500, model.WebResponse{
			Code:   500,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	err = controller.UserService.Delete(c.Request.Context(), uint64(id))
	if err != nil {
		c.JSON(500, model.WebResponse{
			Code:   500,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	c.JSON(200, model.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   nil,
	})
}
