package rest

import (
	"encoding/json"
	"fmt"
	"github.com/arvians-id/apriori/cmd/library/cache"
	"github.com/arvians-id/apriori/internal/http/middleware"
	"github.com/arvians-id/apriori/internal/http/presenter/request"
	response2 "github.com/arvians-id/apriori/internal/http/presenter/response"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/arvians-id/apriori/internal/service"
	"github.com/arvians-id/apriori/util"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type CategoryController struct {
	CategoryService service.CategoryService
	Redis           cache.Redis
}

func NewCategoryController(
	categoryService *service.CategoryService,
	redis *cache.Redis,
) *CategoryController {
	return &CategoryController{
		CategoryService: *categoryService,
		Redis:           *redis,
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
	categoriesCache, err := controller.Redis.Get(c.Request.Context(), "categories")
	if err == redis.Nil {
		categories, err := controller.CategoryService.FindAll(c.Request.Context())
		if err != nil {
			response2.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		err = controller.Redis.Set(c.Request.Context(), "categories", categories)
		if err != nil {
			response2.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		response2.ReturnSuccessOK(c, "OK", categories)
		return
	} else if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	var categoryCacheResponses []model.Category
	err = json.Unmarshal(categoriesCache, &categoryCacheResponses)
	if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response2.ReturnSuccessOK(c, "OK", categoryCacheResponses)
}

func (controller *CategoryController) FindById(c *gin.Context) {
	idParam := util.StrToInt(c.Param("id"))
	category, err := controller.CategoryService.FindById(c.Request.Context(), idParam)
	if err != nil {
		if err.Error() == response2.ErrorNotFound {
			response2.ReturnErrorNotFound(c, err, nil)
			return
		}

		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response2.ReturnSuccessOK(c, "OK", category)
}

func (controller *CategoryController) Create(c *gin.Context) {
	var requestCreate request.CreateCategoryRequest
	if err := c.ShouldBindJSON(&requestCreate); err != nil {
		response2.ReturnErrorBadRequest(c, err, nil)
		return
	}

	category, err := controller.CategoryService.Create(c.Request.Context(), &requestCreate)
	if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	_ = controller.Redis.Del(c.Request.Context(), fmt.Sprintf("categories"))

	response2.ReturnSuccessOK(c, "OK", category)
}

func (controller *CategoryController) Update(c *gin.Context) {
	var requestUpdate request.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&requestUpdate); err != nil {
		response2.ReturnErrorBadRequest(c, err, nil)
		return
	}

	requestUpdate.IdCategory = util.StrToInt(c.Param("id"))
	category, err := controller.CategoryService.Update(c.Request.Context(), &requestUpdate)
	if err != nil {
		if err.Error() == response2.ErrorNotFound {
			response2.ReturnErrorNotFound(c, err, nil)
			return
		}

		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	_ = controller.Redis.Del(c.Request.Context(), fmt.Sprintf("categories"))

	response2.ReturnSuccessOK(c, "OK", category)
}

func (controller *CategoryController) Delete(c *gin.Context) {
	idParam := util.StrToInt(c.Param("id"))
	err := controller.CategoryService.Delete(c.Request.Context(), idParam)
	if err != nil {
		if err.Error() == response2.ErrorNotFound {
			response2.ReturnErrorNotFound(c, err, nil)
			return
		}

		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	_ = controller.Redis.Del(c.Request.Context(), fmt.Sprintf("categories"))

	response2.ReturnSuccessOK(c, "OK", nil)
}
