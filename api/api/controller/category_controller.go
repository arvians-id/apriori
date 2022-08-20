package controller

import (
	"apriori/api/middleware"
	"apriori/api/response"
	"apriori/model"
	"apriori/service"
	"apriori/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type categoryController struct {
	categoryService service.CategoryService
	CacheService    service.CacheService
}

func NewCategoryController(categoryService *service.CategoryService, cacheService *service.CacheService) *categoryController {
	return &categoryController{
		categoryService: *categoryService,
		CacheService:    *cacheService,
	}
}

func (controller *categoryController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.GET("/categories/:id", controller.FindById)
		authorized.POST("/categories", controller.Create)
		authorized.PATCH("/categories/:id", controller.Update)
		authorized.DELETE("/categories/:id", controller.Delete)
	}

	unauthorized := router.Group("/api")
	{
		unauthorized.GET("/categories", controller.FindAll)
	}

	return router
}

func (controller *categoryController) FindAll(c *gin.Context) {
	categoriesCache, err := controller.CacheService.Get(c.Request.Context(), "categories")
	if err == redis.Nil {
		categories, err := controller.categoryService.FindAll(c.Request.Context())
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		err = controller.CacheService.Set(c.Request.Context(), "categories", categories)
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		response.ReturnSuccessOK(c, "OK", categories)
		return
	} else if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	var categoryCacheResponses []model.GetCategoryResponse
	err = json.Unmarshal(bytes.NewBufferString(categoriesCache).Bytes(), &categoryCacheResponses)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", categoryCacheResponses)
}

func (controller *categoryController) FindById(c *gin.Context) {
	idParam := utils.StrToInt(c.Param("id"))
	category, err := controller.categoryService.FindById(c.Request.Context(), idParam)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", category)
}

func (controller *categoryController) Create(c *gin.Context) {
	var request model.CreateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	category, err := controller.categoryService.Create(c.Request.Context(), request)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	_ = controller.CacheService.Del(c.Request.Context(), fmt.Sprintf("categories"))

	response.ReturnSuccessOK(c, "OK", category)
}

func (controller *categoryController) Update(c *gin.Context) {
	var request model.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	request.IdCategory = utils.StrToInt(c.Param("id"))

	category, err := controller.categoryService.Update(c.Request.Context(), request)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	_ = controller.CacheService.Del(c.Request.Context(), fmt.Sprintf("categories"))

	response.ReturnSuccessOK(c, "OK", category)
}

func (controller *categoryController) Delete(c *gin.Context) {
	idParam := utils.StrToInt(c.Param("id"))
	err := controller.categoryService.Delete(c.Request.Context(), idParam)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	_ = controller.CacheService.Del(c.Request.Context(), fmt.Sprintf("categories"))

	response.ReturnSuccessOK(c, "OK", nil)
}
