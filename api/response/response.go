package response

import (
	"apriori/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

func ReturnSuccessOK(c *gin.Context, status string, data interface{}) {
	c.JSON(http.StatusOK, model.WebResponse{
		Code:   http.StatusOK,
		Status: status,
		Data:   data,
	})
}
