package controller

import (
	"apriori/app/middleware"
	"apriori/app/response"
	"apriori/entity"
	"apriori/helper"
	"apriori/model"
	"apriori/service"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type CategoryController struct {
	CategoryService service.CategoryService
	CacheService    service.CacheService
}

func NewCategoryController(
	categoryService *service.CategoryService,
	cacheService *service.CacheService,
) *CategoryController {
	return &CategoryController{
		CategoryService: *categoryService,
		CacheService:    *cacheService,
	}
}

func (controller *CategoryController) Route(router *gin.Engine) *gin.Engine {
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

func (controller *CategoryController) FindAll(c *gin.Context) {
	categoriesCache, err := controller.CacheService.Get(c.Request.Context(), "categories")
	if err == redis.Nil {
		categories, err := controller.CategoryService.FindAll(c.Request.Context())
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

	var categoryCacheResponses []entity.Category
	err = json.Unmarshal(categoriesCache, &categoryCacheResponses)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", categoryCacheResponses)
}

func (controller *CategoryController) FindById(c *gin.Context) {
	idParam := helper.StrToInt(c.Param("id"))
	category, err := controller.CategoryService.FindById(c.Request.Context(), idParam)
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

func (controller *CategoryController) Create(c *gin.Context) {
	var request model.CreateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	category, err := controller.CategoryService.Create(c.Request.Context(), &request)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	_ = controller.CacheService.Del(c.Request.Context(), fmt.Sprintf("categories"))

	response.ReturnSuccessOK(c, "OK", category)
}

func (controller *CategoryController) Update(c *gin.Context) {
	var request model.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	request.IdCategory = helper.StrToInt(c.Param("id"))
	category, err := controller.CategoryService.Update(c.Request.Context(), &request)
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

func (controller *CategoryController) Delete(c *gin.Context) {
	idParam := helper.StrToInt(c.Param("id"))
	err := controller.CategoryService.Delete(c.Request.Context(), idParam)
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
