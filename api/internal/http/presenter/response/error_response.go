package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ErrorNotFound         = "sql: no rows in result set"
	WrongPassword         = "wrong password"
	VerificationExpired   = "reset password verification is expired"
	ResponseErrorNotFound = "data not found"
)

func ReturnErrorNotFound(c *gin.Context, err error, data interface{}) {
	c.JSON(http.StatusNotFound, WebResponse{
		Code:   http.StatusNotFound,
		Status: "data not found",
		Data:   data,
	})
}

func ReturnErrorInternalServerError(c *gin.Context, err error, data interface{}) {
	c.JSON(http.StatusInternalServerError, WebResponse{
		Code:   http.StatusInternalServerError,
		Status: err.Error(),
		Data:   data,
	})
}

func ReturnErrorBadRequest(c *gin.Context, err error, data interface{}) {
	c.JSON(http.StatusBadRequest, WebResponse{
		Code:   http.StatusBadRequest,
		Status: err.Error(),
		Data:   data,
	})
}

func ReturnErrorUnauthorized(c *gin.Context, err error, data interface{}) {
	c.JSON(http.StatusUnauthorized, WebResponse{
		Code:   http.StatusUnauthorized,
		Status: err.Error(),
		Data:   data,
	})
}
