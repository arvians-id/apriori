package rest

import (
	"errors"
	"fmt"
	"github.com/arvians-id/apriori/internal/http/middleware"
	"github.com/arvians-id/apriori/internal/http/presenter/request"
	"github.com/arvians-id/apriori/internal/http/presenter/response"
	"github.com/arvians-id/apriori/internal/service"
	"github.com/gin-gonic/gin"
	"log"
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
	jwtService *service.JwtService,
	emailService *service.EmailService,
	passwordResetService *service.PasswordResetService,
) *AuthController {
	return &AuthController{
		UserService:          *userService,
		JwtService:           *jwtService,
		EmailService:         *emailService,
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
	var requestCredential request.GetUserCredentialRequest
	err := c.ShouldBindJSON(&requestCredential)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	user, err := controller.UserService.FindByEmail(c.Request.Context(), &requestCredential)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		if err.Error() == response.WrongPassword {
			response.ReturnErrorBadRequest(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	expiredTimeAccess, err := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRED_TIME"))
	if err != nil {
		log.Println("[AuthController][Login] problem in conversion string to integer, err: ", err.Error())
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	expirationTime := time.Now().Add(time.Duration(expiredTimeAccess) * 24 * time.Hour)
	token, err := controller.JwtService.GenerateToken(user.IdUser, user.Role, expirationTime)
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
	var requestToken request.GetRefreshTokenRequest
	err := c.ShouldBindJSON(&requestToken)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	token, err := controller.JwtService.RefreshToken(requestToken.RefreshToken)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	expiredTimeAccess, err := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRED_TIME"))
	if err != nil {
		log.Println("[AuthController][Login] problem in conversion string to integer, err: ", err.Error())
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

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
	var requestCreate request.CreateUserRequest
	err := c.ShouldBindJSON(&requestCreate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	user, err := controller.UserService.Create(c.Request.Context(), &requestCreate)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "created", user)
}

func (controller *AuthController) ForgotPassword(c *gin.Context) {
	var requestCreate request.CreatePasswordResetRequest
	err := c.ShouldBindJSON(&requestCreate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	result, err := controller.PasswordResetService.CreateOrUpdateByEmail(c.Request.Context(), requestCreate.Email)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	message := fmt.Sprintf("%s/auth/reset-password?signature=%v", os.Getenv("APP_URL_FE"), result.Token)
	err = controller.EmailService.SendEmailWithText(result.Email, "Forgot Password", &message)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "mail sent successfully", gin.H{
		"signature": result.Token,
	})
}

func (controller *AuthController) VerifyResetPassword(c *gin.Context) {
	var requestUpdate request.UpdateResetPasswordUserRequest
	err := c.ShouldBindJSON(&requestUpdate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	requestUpdate.Token = c.Query("signature")
	err = controller.PasswordResetService.Verify(c.Request.Context(), &requestUpdate)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		if err.Error() == response.VerificationExpired {
			response.ReturnErrorBadRequest(c, err, nil)
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
