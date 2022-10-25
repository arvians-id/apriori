package rest

import (
	"github.com/arvians-id/apriori/cmd/library/aws"
	"github.com/arvians-id/apriori/internal/http/middleware"
	"github.com/arvians-id/apriori/internal/http/presenter/request"
	"github.com/arvians-id/apriori/internal/http/presenter/response"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/arvians-id/apriori/internal/service"
	"github.com/arvians-id/apriori/util"
	"github.com/gin-gonic/gin"
	"strings"
)

type AprioriController struct {
	AprioriService service.AprioriService
	StorageS3      aws.StorageS3
}

func NewAprioriController(aprioriService service.AprioriService, storageS3 *aws.StorageS3) *AprioriController {
	return &AprioriController{
		AprioriService: aprioriService,
		StorageS3:      *storageS3,
	}
}

func (controller *AprioriController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.PATCH("/apriori/:code", controller.UpdateStatus)
		authorized.POST("/apriori", controller.Create)
		authorized.DELETE("/apriori/:code", controller.Delete)
		authorized.PATCH("/apriori/:code/update/:id", controller.Update)
		authorized.POST("/apriori/generate", controller.Generate)
	}

	unauthorized := router.Group("/api")
	{
		unauthorized.GET("/apriori", controller.FindAll)
		unauthorized.GET("/apriori/:code", controller.FindAllByCode)
		unauthorized.GET("/apriori/:code/detail/:id", controller.FindByCodeAndId)
		unauthorized.GET("/apriori/actives", controller.FindAllByActive)
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

func (controller *AprioriController) FindAllByActive(c *gin.Context) {
	apriories, err := controller.AprioriService.FindAllByActive(c.Request.Context())
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", apriories)
}

func (controller *AprioriController) FindAllByCode(c *gin.Context) {
	codeParam := c.Param("code")
	apriories, err := controller.AprioriService.FindAllByCode(c.Request.Context(), codeParam)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", apriories)
}

func (controller *AprioriController) FindByCodeAndId(c *gin.Context) {
	codeParam := c.Param("code")
	idParam := util.StrToInt(c.Param("id"))
	apriori, err := controller.AprioriService.FindByCodeAndId(c.Request.Context(), codeParam, idParam)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", apriori)
}

func (controller *AprioriController) Update(c *gin.Context) {
	codeParam := c.Param("code")
	idParam := util.StrToInt(c.Param("id"))
	description := c.PostForm("description")
	file, header, err := c.Request.FormFile("image")
	filePath := ""
	if err == nil {
		pathName, err := controller.StorageS3.UploadFileS3(file, header)
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}
		filePath = pathName
	}

	var requestUpdate request.UpdateAprioriRequest
	requestUpdate.IdApriori = idParam
	requestUpdate.Code = codeParam
	requestUpdate.Description = description
	requestUpdate.Image = filePath
	apriories, err := controller.AprioriService.Update(c.Request.Context(), &requestUpdate)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", apriories)
}

func (controller *AprioriController) UpdateStatus(c *gin.Context) {
	codeParam := c.Param("code")
	err := controller.AprioriService.UpdateStatus(c.Request.Context(), codeParam)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)
}

func (controller *AprioriController) Create(c *gin.Context) {
	var generateRequests []*model.GenerateApriori
	err := c.ShouldBindJSON(&generateRequests)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	var aprioriRequests []*request.CreateAprioriRequest
	for _, generateRequest := range generateRequests {
		ItemSet := strings.Join(generateRequest.ItemSet, ", ")
		aprioriRequests = append(aprioriRequests, &request.CreateAprioriRequest{
			Item:       ItemSet,
			Discount:   generateRequest.Discount,
			Support:    generateRequest.Support,
			Confidence: generateRequest.Confidence,
			RangeDate:  generateRequest.RangeDate,
		})
	}

	err = controller.AprioriService.Create(c.Request.Context(), aprioriRequests)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)
}
func (controller *AprioriController) Delete(c *gin.Context) {
	codeParam := c.Param("code")
	err := controller.AprioriService.Delete(c.Request.Context(), codeParam)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)

}
func (controller *AprioriController) Generate(c *gin.Context) {
	var requestGenerate request.GenerateAprioriRequest
	err := c.ShouldBindJSON(&requestGenerate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	apriori, err := controller.AprioriService.Generate(c.Request.Context(), &requestGenerate)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", apriori)
}
