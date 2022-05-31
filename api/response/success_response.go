package response

import (
	"apriori/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ReturnSuccessOK(c *gin.Context, status string, data interface{}) {
	c.JSON(http.StatusOK, model.WebResponse{
		Code:   http.StatusOK,
		Status: status,
		Data:   data,
	})
}
