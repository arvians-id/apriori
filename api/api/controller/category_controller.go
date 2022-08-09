package controller

import (
	"apriori/api/middleware"
	"apriori/api/response"
	"apriori/model"
	"apriori/service"
	"apriori/utils"
	"github.com/gin-gonic/gin"
)

type categoryController struct {
	categoryService service.CategoryService
}

func NewCategoryController(categoryService *service.CategoryService) *categoryController {
	return &categoryController{
		categoryService: *categoryService,
	}
}

func (controller *categoryController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.GET("/categories", controller.FindAll)
		authorized.GET("/categories/:id", controller.FindById)
		authorized.POST("/categories", controller.Create)
		authorized.PATCH("/categories/:id", controller.Update)
		authorized.DELETE("/categories/:id", controller.Delete)
	}

	return router
}

func (controller *categoryController) FindAll(c *gin.Context) {
	categories, err := controller.categoryService.FindAll(c.Request.Context())
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", categories)
}

func (controller *categoryController) FindById(c *gin.Context) {
	id := utils.StrToInt(c.Param("id"))
	category, err := controller.categoryService.FindById(c.Request.Context(), id)
	if err != nil {
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
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", category)
}

func (controller *categoryController) Delete(c *gin.Context) {
	id := utils.StrToInt(c.Param("id"))
	err := controller.categoryService.Delete(c.Request.Context(), id)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)
}
