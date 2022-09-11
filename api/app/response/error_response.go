package response

import (
	"apriori/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ErrorNotFound       = "sql: no rows in result set"
	WrongPassword       = "wrong password"
	VerificationExpired = "reset password verification is expired"
)

func ReturnErrorNotFound(c *gin.Context, err error, data interface{}) {
	c.JSON(http.StatusNotFound, model.WebResponse{
		Code:   http.StatusNotFound,
		Status: "data not found",
		Data:   data,
	})
}

func ReturnErrorInternalServerError(c *gin.Context, err error, data interface{}) {
	c.JSON(http.StatusInternalServerError, model.WebResponse{
		Code:   http.StatusInternalServerError,
		Status: err.Error(),
		Data:   data,
	})
}

func ReturnErrorBadRequest(c *gin.Context, err error, data interface{}) {
	c.JSON(http.StatusBadRequest, model.WebResponse{
		Code:   http.StatusBadRequest,
		Status: err.Error(),
		Data:   data,
	})
}

func ReturnErrorUnauthorized(c *gin.Context, err error, data interface{}) {
	c.JSON(http.StatusUnauthorized, model.WebResponse{
		Code:   http.StatusUnauthorized,
		Status: err.Error(),
		Data:   data,
	})
}
