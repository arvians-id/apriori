package controller

import (
	"apriori/api/middleware"
	"apriori/api/response"
	"apriori/model"
	"apriori/service"
	"github.com/gin-gonic/gin"
)

type AprioriController struct {
	AprioriService service.AprioriService
}

func NewAprioriController(aprioriService service.AprioriService) *AprioriController {
	return &AprioriController{
		AprioriService: aprioriService,
	}
}

func (controller *AprioriController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.POST("/apriori", controller.Generate)
	}

	return router
}

func (controller *AprioriController) Generate(c *gin.Context) {
	var request model.GenerateAprioriRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	transactions, err := controller.AprioriService.Generate(c.Request.Context(), request)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", transactions)
}
