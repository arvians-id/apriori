package controller

import (
	"apriori/api/response"
	"apriori/model"
	"apriori/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"time"
)

type AuthController struct {
	AuthService          service.AuthService
	UserService          service.UserService
	JwtService           service.JwtService
	EmailService         service.EmailService
	PasswordResetService service.PasswordResetService
}

func NewAuthController(authService *service.AuthService, userService *service.UserService, jwtService service.JwtService, emailService service.EmailService, passwordResetService *service.PasswordResetService) *AuthController {
	return &AuthController{
		AuthService:          *authService,
		UserService:          *userService,
		JwtService:           jwtService,
		EmailService:         emailService,
		PasswordResetService: *passwordResetService,
	}
}

func (controller *AuthController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api/auth")
	{
		authorized.POST("/login", controller.login)
		authorized.POST("/forgot-password", controller.forgotPassword)
		authorized.POST("/verify", controller.verifyResetPassword)
		authorized.DELETE("/logout", controller.logout)
	}

	return router
}

func (controller *AuthController) login(c *gin.Context) {
	var request model.GetUserCredentialRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	user, err := controller.AuthService.VerifyCredential(c.Request.Context(), request)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	token, err := controller.JwtService.GenerateToken(user.IdUser, expirationTime)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    url.QueryEscape(token),
		Expires:  expirationTime,
		Path:     "/api",
		HttpOnly: true,
	})

	response.ReturnSuccessOK(c, "OK", gin.H{
		"token":     token,
		"expiresAt": expirationTime,
	})
}

func (controller *AuthController) forgotPassword(c *gin.Context) {
	// Check email if exists
	var request model.CreatePasswordResetRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	// Insert or update data token into database
	result, err := controller.PasswordResetService.CreateOrUpdate(c.Request.Context(), request.Email)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// Send email to user
	message := fmt.Sprintf("http://localhost:8080/verify?signature=%v", result.Token)
	err = controller.EmailService.SendEmailWithText(result.Email, message)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "mail sent successfully", nil)
}

func (controller *AuthController) verifyResetPassword(c *gin.Context) {
	// Check email if exists
	var request model.UpdateResetPasswordUserRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	request.Token = c.Query("signature")
	err = controller.PasswordResetService.Verify(c.Request.Context(), request)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "updated", nil)
}

func (controller *AuthController) logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/api",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	response.ReturnSuccessOK(c, "OK", nil)
}
