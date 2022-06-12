package controller

import (
	"apriori/api/response"
	"apriori/model"
	"apriori/service"
	"fmt"
	"github.com/gin-gonic/gin"
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
		authorized.POST("/login", controller.Login)
		authorized.POST("/refresh", controller.Refresh)
		authorized.POST("/forgot-password", controller.ForgotPassword)
		authorized.POST("/verify", controller.VerifyResetPassword)
		authorized.POST("/register", controller.Register)
		authorized.DELETE("/logout", controller.Logout)
	}

	return router
}

func (controller *AuthController) Login(c *gin.Context) {
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

	expirationTime := time.Now().Add(1 * time.Hour)
	token, err := controller.JwtService.GenerateToken(user.IdUser, expirationTime)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", gin.H{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}

func (controller *AuthController) Refresh(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	token, err := controller.JwtService.RefreshToken(request.RefreshToken)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", gin.H{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}

func (controller *AuthController) Register(c *gin.Context) {
	var request model.CreateUserRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	user, err := controller.UserService.Create(c.Request.Context(), request)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "created", user)
}

func (controller *AuthController) ForgotPassword(c *gin.Context) {
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

	response.ReturnSuccessOK(c, "mail sent successfully", gin.H{
		"signature": result.Token,
	})
}

func (controller *AuthController) VerifyResetPassword(c *gin.Context) {
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

func (controller *AuthController) Logout(c *gin.Context) {
	response.ReturnSuccessOK(c, "OK", nil)
}
