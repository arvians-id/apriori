package controller

import (
	"apriori/model"
	"apriori/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"time"
)

type AuthController struct {
	AuthService service.AuthService
	UserService service.UserService
	JwtService  service.JwtService
}

func NewAuthController(authService service.AuthService, userService service.UserService, jwtService service.JwtService) *AuthController {
	return &AuthController{
		AuthService: authService,
		UserService: userService,
		JwtService:  jwtService,
	}
}

func (controller *AuthController) Route(router *gin.Engine) *gin.Engine {
	router.POST("/api/login", controller.login)
	router.POST("/api/register", controller.register)

	return router
}
func (controller *AuthController) login(c *gin.Context) {
	var request model.GetUserCredentialRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(500, model.WebResponse{
			Code:   500,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	user, err := controller.AuthService.VerifyCredential(c.Request.Context(), request)
	if err != nil {
		c.JSON(500, model.WebResponse{
			Code:   500,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	token, err := controller.JwtService.GenerateToken(user.IdUser)
	if err != nil {
		c.JSON(500, model.WebResponse{
			Code:   500,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    url.QueryEscape(token),
		Expires:  expirationTime,
		HttpOnly: true,
	})

	c.JSON(200, model.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   token,
	})
}

func (controller *AuthController) register(c *gin.Context) {
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
