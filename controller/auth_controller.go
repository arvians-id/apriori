package controller

import (
	"apriori/model"
	"apriori/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type AuthController struct {
	AuthService          service.AuthService
	UserService          service.UserService
	JwtService           service.JwtService
	EmailService         service.EmailService
	PasswordResetService service.PasswordResetService
}

func NewAuthController(authService service.AuthService, userService service.UserService, jwtService service.JwtService, emailService service.EmailService, passwordResetService service.PasswordResetService) *AuthController {
	return &AuthController{
		AuthService:          authService,
		UserService:          userService,
		JwtService:           jwtService,
		EmailService:         emailService,
		PasswordResetService: passwordResetService,
	}
}

func (controller *AuthController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api/auth")
	{
		authorized.POST("/login", controller.login)
		authorized.POST("/forgot-password", controller.forgotPassword)
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

func (controller *AuthController) forgotPassword(c *gin.Context) {
	// Check email if exists
	var request model.CreatePasswordResetRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	user, err := controller.UserService.FindByEmail(c.Request.Context(), request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	// Insert data token to database
	times := time.Now()
	timestamp := times.Unix()

	timestampString := strconv.Itoa(int(timestamp))
	token := controller.EmailService.MakeTokenVerificationEmail(user.Email, timestampString)

	result, err := controller.PasswordResetService.Create(c.Request.Context(), request, token, int32(timestamp))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: err.Error(),
			Data:   nil,
		})
		return
	}

	// Send email to user
	message := "http://localhost:8080/verify?signature=" + token + "&expired=" + timestampString + ""
	err = controller.EmailService.SendEmailWithText(user.Email, message)
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
		Status: "Mail Success sent!",
		Data: gin.H{
			"email":   user.Email,
			"token":   result.Token,
			"expired": result.Expired,
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
