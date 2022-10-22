package rest

import (
	"github.com/arvians-id/apriori/internal/http/middleware"
	"github.com/arvians-id/apriori/internal/http/presenter/request"
	"github.com/arvians-id/apriori/internal/http/presenter/response"
	"github.com/arvians-id/apriori/internal/service"
	"github.com/arvians-id/apriori/util"
	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	CategoryService service.CategoryService
}

func NewCategoryController(categoryService *service.CategoryService) *CategoryController {
	return &CategoryController{
		CategoryService: *categoryService,
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
	categories, err := controller.CategoryService.FindAll(c.Request.Context())
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", categories)
}

func (controller *CategoryController) FindById(c *gin.Context) {
	idParam := util.StrToInt(c.Param("id"))
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
	var requestCreate request.CreateCategoryRequest
	if err := c.ShouldBindJSON(&requestCreate); err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	category, err := controller.CategoryService.Create(c.Request.Context(), &requestCreate)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", category)
}

func (controller *CategoryController) Update(c *gin.Context) {
	var requestUpdate request.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&requestUpdate); err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	requestUpdate.IdCategory = util.StrToInt(c.Param("id"))
	category, err := controller.CategoryService.Update(c.Request.Context(), &requestUpdate)
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

func (controller *CategoryController) Delete(c *gin.Context) {
	idParam := util.StrToInt(c.Param("id"))
	err := controller.CategoryService.Delete(c.Request.Context(), idParam)
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
