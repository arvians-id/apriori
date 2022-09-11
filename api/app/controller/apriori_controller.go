package controller

import (
	"apriori/app/middleware"
	"apriori/app/response"
	"apriori/helper"
	"apriori/model"
	"apriori/service"
	"bytes"
	"encoding/json"
	"fmt"
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
	idParam := helper.StrToInt(c.Param("id"))
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
	idParam := helper.StrToInt(c.Param("id"))
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
	request.IdApriori = idParam
	request.Code = codeParam
	request.Description = description
	request.Image = filePath
	apriories, err := controller.AprioriService.Update(c.Request.Context(), &request)
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
	var generateRequests []*model.GetGenerateAprioriResponse
	err := c.ShouldBindJSON(&generateRequests)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	var aprioriRequests []*model.CreateAprioriRequest
	for _, generateRequest := range generateRequests {
		ItemSet := strings.Join(generateRequest.ItemSet, ", ")
		aprioriRequests = append(aprioriRequests, &model.CreateAprioriRequest{
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
	var request model.GenerateAprioriRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	key := fmt.Sprintf(
		"%v%v%v%v%s%s",
		request.MinimumDiscount,
		request.MaximumDiscount,
		request.MinimumSupport,
		request.MinimumConfidence,
		request.StartDate,
		request.EndDate,
	)
	aprioriCache, err := controller.CacheService.Get(c, key)
	if err == redis.Nil {
		apriori, err := controller.AprioriService.Generate(c.Request.Context(), &request)
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		err = controller.CacheService.Set(c.Request.Context(), key, apriori)
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		response.ReturnSuccessOK(c, "OK", apriori)
		return
	} else if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	var aprioriCacheResponses []model.GetGenerateAprioriResponse
	err = json.Unmarshal(bytes.NewBufferString(aprioriCache).Bytes(), &aprioriCacheResponses)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", aprioriCacheResponses)
}
