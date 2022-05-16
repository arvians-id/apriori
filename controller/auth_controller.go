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
	authorized := router.Group("/api/auth")
	{
		authorized.POST("/login", controller.login)
		authorized.DELETE("/logout", controller.logout)
	}

	return router
}

func (controller *AuthController) login(c *gin.Context) {
	var request model.GetUserCredentialRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	user, err := controller.AuthService.VerifyCredential(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	token, err := controller.JwtService.GenerateToken(user.IdUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse{
			Code:   http.StatusInternalServerError,
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
		Path:     "/api",
		HttpOnly: true,
	})

	c.JSON(http.StatusOK, model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data: gin.H{
			"token":     token,
			"expiresAt": expirationTime,
		},
	})
}

func (controller *AuthController) logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/api",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	c.JSON(http.StatusOK, model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   nil,
	})
}
