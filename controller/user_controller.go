package controller

import (
	"apriori/middleware"
	"apriori/model"
	"apriori/service"
	"github.com/gin-gonic/gin"
	"net/http"
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
		c.JSON(http.StatusUnauthorized, model.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "unauthorized",
			Data:   nil,
		})
		return
	}

	user, err := controller.UserService.FindById(c.Request.Context(), uint64(id.(float64)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   user,
	})
}

func (controller *UserController) FindAll(c *gin.Context) {
	users, err := controller.UserService.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   users,
	})
}

func (controller *UserController) FindById(c *gin.Context) {
	params := c.Param("userId")
	id, err := strconv.Atoi(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	user, err := controller.UserService.FindById(c.Request.Context(), uint64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   user,
	})
}

func (controller *UserController) Create(c *gin.Context) {
	var request model.CreateUserRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	user, err := controller.UserService.Create(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.WebResponse{
		Code:   http.StatusOK,
		Status: "successfully created",
		Data:   user,
	})
}

func (controller *UserController) Update(c *gin.Context) {
	params := c.Param("userId")
	id, err := strconv.Atoi(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	var request model.UpdateUserRequest
	err = c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	request.IdUser = uint64(id)

	user, err := controller.UserService.Update(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.WebResponse{
		Code:   http.StatusOK,
		Status: "successfully updated",
		Data:   user,
	})
}

func (controller *UserController) Delete(c *gin.Context) {
	params := c.Param("userId")
	id, err := strconv.Atoi(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	err = controller.UserService.Delete(c.Request.Context(), uint64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.WebResponse{
		Code:   http.StatusOK,
		Status: "successfully deleted",
		Data:   nil,
	})
}
