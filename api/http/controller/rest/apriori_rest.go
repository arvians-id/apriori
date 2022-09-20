package rest

import (
	"encoding/json"
	"fmt"
	"github.com/arvians-id/apriori/helper"
	"github.com/arvians-id/apriori/http/controller/rest/request"
	response2 "github.com/arvians-id/apriori/http/controller/rest/response"
	"github.com/arvians-id/apriori/http/middleware"
	"github.com/arvians-id/apriori/model"
	"github.com/arvians-id/apriori/service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"strings"
)

type AprioriController struct {
	AprioriService service.AprioriService
	StorageService service.StorageService
	CacheService   service.CacheService
}

func NewAprioriController(
	aprioriService service.AprioriService,
	storageService *service.StorageService,
	cacheService *service.CacheService,
) *AprioriController {
	return &AprioriController{
		AprioriService: aprioriService,
		StorageService: *storageService,
		CacheService:   *cacheService,
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
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response2.ReturnSuccessOK(c, "OK", apriories)
}

func (controller *AprioriController) FindAllByActive(c *gin.Context) {
	apriories, err := controller.AprioriService.FindAllByActive(c.Request.Context())
	if err != nil {
		if err.Error() == response2.ErrorNotFound {
			response2.ReturnErrorNotFound(c, err, nil)
			return
		}

		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response2.ReturnSuccessOK(c, "OK", apriories)
}

func (controller *AprioriController) FindAllByCode(c *gin.Context) {
	codeParam := c.Param("code")
	apriories, err := controller.AprioriService.FindAllByCode(c.Request.Context(), codeParam)
	if err != nil {
		if err.Error() == response2.ErrorNotFound {
			response2.ReturnErrorNotFound(c, err, nil)
			return
		}

		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response2.ReturnSuccessOK(c, "OK", apriories)
}

func (controller *AprioriController) FindByCodeAndId(c *gin.Context) {
	codeParam := c.Param("code")
	idParam := helper.StrToInt(c.Param("id"))
	apriori, err := controller.AprioriService.FindByCodeAndId(c.Request.Context(), codeParam, idParam)
	if err != nil {
		if err.Error() == response2.ErrorNotFound {
			response2.ReturnErrorNotFound(c, err, nil)
			return
		}

		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response2.ReturnSuccessOK(c, "OK", apriori)
}

func (controller *AprioriController) Update(c *gin.Context) {
	codeParam := c.Param("code")
	idParam := helper.StrToInt(c.Param("id"))
	description := c.PostForm("description")
	file, header, err := c.Request.FormFile("image")
	filePath := ""
	if err == nil {
		pathName, err := controller.StorageService.UploadFileS3(file, header)
		if err != nil {
			response2.ReturnErrorInternalServerError(c, err, nil)
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
		if err.Error() == response2.ErrorNotFound {
			response2.ReturnErrorNotFound(c, err, nil)
			return
		}

		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response2.ReturnSuccessOK(c, "OK", apriories)
}

func (controller *AprioriController) UpdateStatus(c *gin.Context) {
	codeParam := c.Param("code")
	err := controller.AprioriService.UpdateStatus(c.Request.Context(), codeParam)
	if err != nil {
		if err.Error() == response2.ErrorNotFound {
			response2.ReturnErrorNotFound(c, err, nil)
			return
		}

		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response2.ReturnSuccessOK(c, "OK", nil)
}

func (controller *AprioriController) Create(c *gin.Context) {
	var generateRequests []*model.GenerateApriori
	err := c.ShouldBindJSON(&generateRequests)
	if err != nil {
		response2.ReturnErrorBadRequest(c, err, nil)
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
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response2.ReturnSuccessOK(c, "OK", nil)
}
func (controller *AprioriController) Delete(c *gin.Context) {
	codeParam := c.Param("code")
	err := controller.AprioriService.Delete(c.Request.Context(), codeParam)
	if err != nil {
		if err.Error() == response2.ErrorNotFound {
			response2.ReturnErrorNotFound(c, err, nil)
			return
		}

		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response2.ReturnSuccessOK(c, "OK", nil)

}
func (controller *AprioriController) Generate(c *gin.Context) {
	var requestGenerate request.GenerateAprioriRequest
	err := c.ShouldBindJSON(&requestGenerate)
	if err != nil {
		response2.ReturnErrorBadRequest(c, err, nil)
		return
	}

	key := fmt.Sprintf(
		"%v%v%v%v%s%s",
		requestGenerate.MinimumDiscount,
		requestGenerate.MaximumDiscount,
		requestGenerate.MinimumSupport,
		requestGenerate.MinimumConfidence,
		requestGenerate.StartDate,
		requestGenerate.EndDate,
	)
	aprioriCache, err := controller.CacheService.Get(c, key)
	if err == redis.Nil {
		apriori, err := controller.AprioriService.Generate(c.Request.Context(), &requestGenerate)
		if err != nil {
			response2.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		err = controller.CacheService.Set(c.Request.Context(), key, apriori)
		if err != nil {
			response2.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		response2.ReturnSuccessOK(c, "OK", apriori)
		return
	} else if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	var aprioriCacheResponses []model.GenerateApriori
	err = json.Unmarshal(aprioriCache, &aprioriCacheResponses)
	if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response2.ReturnSuccessOK(c, "OK", aprioriCacheResponses)
}
