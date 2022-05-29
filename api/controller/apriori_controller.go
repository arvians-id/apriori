package controller

import (
	"apriori/api/middleware"
	"apriori/api/response"
	"apriori/model"
	"apriori/service"
	"github.com/gin-gonic/gin"
	"strings"
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
		authorized.GET("/apriori", controller.FindAll)
		authorized.GET("/apriori/:code", controller.FindByCode)
		authorized.PATCH("/apriori/:code", controller.ChangeActive)
		authorized.POST("/apriori", controller.Create)
		authorized.DELETE("/apriori/:code", controller.Delete)
		authorized.POST("/apriori/generate", controller.Generate)
	}

	return router
}

func (controller *AprioriController) FindAll(c *gin.Context) {
	apriories, err := controller.AprioriService.FindAll(c.Request.Context())
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", apriories)
}

func (controller *AprioriController) FindByCode(c *gin.Context) {
	code := c.Param("code")
	apriori, err := controller.AprioriService.FindByCode(c.Request.Context(), code)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", apriori)
}

func (controller *AprioriController) ChangeActive(c *gin.Context) {
	code := c.Param("code")
	err := controller.AprioriService.ChangeActive(c.Request.Context(), code)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)
}

func (controller *AprioriController) Create(c *gin.Context) {
	var requestGenerate []model.GetGenerateAprioriResponse
	err := c.ShouldBindJSON(&requestGenerate)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	var request []model.CreateAprioriRequest
	for _, property := range requestGenerate {
		if property.Description == "Eligible" {
			ItemSet := strings.Join(property.ItemSet, ", ")

			request = append(request, model.CreateAprioriRequest{
				Item:       ItemSet,
				Discount:   property.Discount,
				Support:    property.Support,
				Confidence: property.Confidence,
				RangeDate:  property.RangeDate,
			})
		}
	}

	err = controller.AprioriService.Create(c.Request.Context(), request)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)
}
func (controller *AprioriController) Delete(c *gin.Context) {
	code := c.Param("code")
	err := controller.AprioriService.Delete(c.Request.Context(), code)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)

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
