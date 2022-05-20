package controller

import (
	"apriori/helper"
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
		helper.ReturnErrorBadRequest(c, err, nil)
		return
	}

	user, err := controller.AuthService.VerifyCredential(c.Request.Context(), request)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	token, err := controller.JwtService.GenerateToken(user.IdUser)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
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

	helper.ReturnSuccessOK(c, "OK", gin.H{
		"token":     token,
		"expiresAt": expirationTime,
	})
}

func (controller *AuthController) forgotPassword(c *gin.Context) {
	// Check email if exists
	var request model.CreatePasswordResetRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		helper.ReturnErrorBadRequest(c, err, nil)
		return
	}

	// Insert or update data token into database
	result, err := controller.PasswordResetService.CreateOrUpdate(c.Request.Context(), request.Email)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// Send email to user
	message := fmt.Sprintf("http://localhost:8080/verify?signature=%v", result.Token)
	err = controller.EmailService.SendEmailWithText(result.Email, message)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "mail sent successfully", result)
}

func (controller *AuthController) verifyResetPassword(c *gin.Context) {
	// Check email if exists
	var request model.UpdateResetPasswordUserRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		helper.ReturnErrorBadRequest(c, err, nil)
		return
	}

	request.Token = c.Query("signature")
	user, err := controller.PasswordResetService.Verify(c.Request.Context(), request)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "successfully updated", user)
}

func (controller *AuthController) logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/api",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	helper.ReturnSuccessOK(c, "OK", nil)
}
