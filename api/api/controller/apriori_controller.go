package controller

import (
	"apriori/api/middleware"
	"apriori/api/response"
	"apriori/model"
	"apriori/service"
	"apriori/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

type AprioriController struct {
	AprioriService service.AprioriService
	StorageService service.StorageService
}

func NewAprioriController(aprioriService service.AprioriService, storageService *service.StorageService) *AprioriController {
	return &AprioriController{
		AprioriService: aprioriService,
		StorageService: *storageService,
	}
}

func (controller *AprioriController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.PATCH("/apriori/:code", controller.ChangeActive)
		authorized.POST("/apriori", controller.Create)
		authorized.DELETE("/apriori/:code", controller.Delete)
		authorized.PATCH("/apriori/:code/update/:id", controller.UpdateApriori)
		authorized.POST("/apriori/generate", controller.Generate)
	}

	router.GET("/api/apriori", controller.FindAll)
	router.GET("/api/apriori/:code", controller.FindByCode)
	router.GET("/api/apriori/:code/detail/:id", controller.FindAprioriById)
	router.GET("/api/apriori/actives", controller.FindByActive)

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

func (controller *AprioriController) FindByActive(c *gin.Context) {
	apriories, err := controller.AprioriService.FindByActive(c.Request.Context())
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

func (controller *AprioriController) FindAprioriById(c *gin.Context) {
	code := c.Param("code")
	id := utils.StrToInt(c.Param("id"))
	apriori, err := controller.AprioriService.FindAprioriById(c.Request.Context(), code, id)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", apriori)
}

func (controller *AprioriController) UpdateApriori(c *gin.Context) {
	code := c.Param("code")
	id := utils.StrToInt(c.Param("id"))
	description := c.PostForm("description")

	file, header, err := c.Request.FormFile("image")
	filePath := ""
	if err == nil {
		pathName, err := controller.StorageService.UploadFileS3(file, header)
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}
		filePath = pathName
	}

	var request model.UpdateAprioriRequest

	request.IdApriori = uint64(id)
	request.Code = code
	request.Description = description
	request.Image = filePath
	apriories, err := controller.AprioriService.UpdateApriori(c.Request.Context(), request)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", apriories)
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
		ItemSet := strings.Join(property.ItemSet, ", ")

		request = append(request, model.CreateAprioriRequest{
			Item:       ItemSet,
			Discount:   property.Discount,
			Support:    property.Support,
			Confidence: property.Confidence,
			RangeDate:  property.RangeDate,
		})
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
