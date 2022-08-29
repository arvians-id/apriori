package controller

import (
	"apriori/app/middleware"
	"apriori/app/response"
	"apriori/model"
	"apriori/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type AuthController struct {
	UserService          service.UserService
	JwtService           service.JwtService
	EmailService         service.EmailService
	PasswordResetService service.PasswordResetService
}

func NewAuthController(
	userService *service.UserService,
	jwtService service.JwtService,
	emailService service.EmailService,
	passwordResetService *service.PasswordResetService,
) *AuthController {
	return &AuthController{
		UserService:          *userService,
		JwtService:           jwtService,
		EmailService:         emailService,
		PasswordResetService: *passwordResetService,
	}
}

func (controller *AuthController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api/auth", middleware.AuthJwtMiddleware())
	{
		authorized.GET("/token", controller.Token)
	}

	unauthorized := router.Group("/api/auth")
	{
		unauthorized.POST("/login", controller.Login)
		unauthorized.POST("/refresh", controller.Refresh)
		unauthorized.POST("/forgot-password", controller.ForgotPassword)
		unauthorized.POST("/verify", controller.VerifyResetPassword)
		unauthorized.POST("/register", controller.Register)
		unauthorized.DELETE("/logout", controller.Logout)
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

	user, err := controller.UserService.FindByEmail(c.Request.Context(), &request)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	expiredTimeAccess, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRED_TIME"))
	expirationTime := time.Now().Add(time.Duration(expiredTimeAccess) * 24 * time.Hour)
	token, err := controller.JwtService.GenerateToken(user.IdUser, expirationTime)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    url.QueryEscape(token.AccessToken),
		Expires:  expirationTime,
		Path:     "/api",
		HttpOnly: true,
	})

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

	expiredTimeAccess, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRED_TIME"))
	expirationTime := time.Now().Add(time.Duration(expiredTimeAccess) * 24 * time.Hour)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    url.QueryEscape(token.AccessToken),
		Expires:  expirationTime,
		Path:     "/api",
		HttpOnly: true,
	})

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

	user, err := controller.UserService.Create(c.Request.Context(), &request)
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
	result, err := controller.PasswordResetService.CreateOrUpdateByEmail(c.Request.Context(), request.Email)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// Send email to user
	message := fmt.Sprintf("%s/auth/reset-password?signature=%v", os.Getenv("APP_URL_FE"), result.Token)
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
	err = controller.PasswordResetService.Verify(c.Request.Context(), &request)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "updated", nil)
}

func (controller *AuthController) Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/api",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	response.ReturnSuccessOK(c, "OK", nil)
}

func (controller *AuthController) Token(c *gin.Context) {
	_, isExist := c.Get("id_user")
	if !isExist {
		response.ReturnErrorUnauthorized(c, errors.New("unauthorized"), nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)
}
