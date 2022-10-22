package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ReturnSuccessOK(c *gin.Context, status string, data interface{}) {
	c.JSON(http.StatusOK, WebResponse{
		Code:   http.StatusOK,
		Status: status,
		Data:   data,
	})
}
