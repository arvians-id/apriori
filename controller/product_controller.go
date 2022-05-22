package controller

import (
	"apriori/helper"
	"apriori/middleware"
	"apriori/model"
	"apriori/service"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	ProductService service.ProductService
}

func NewProductController(productService *service.ProductService) *ProductController {
	return &ProductController{
		ProductService: *productService,
	}
}

func (controller *ProductController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.GET("/products", controller.FindAll)
		authorized.GET("/products/:code", controller.FindById)
		authorized.POST("/products", controller.Create)
		authorized.PATCH("/products/:code", controller.Update)
		authorized.DELETE("/products/:code", controller.Delete)
	}

	return router
}

func (controller *ProductController) FindAll(c *gin.Context) {
	products, err := controller.ProductService.FindAll(c.Request.Context())
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "OK", products)
}

func (controller *ProductController) FindById(c *gin.Context) {
	params := c.Param("code")

	product, err := controller.ProductService.FindByCode(c.Request.Context(), params)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "OK", product)
}

func (controller *ProductController) Create(c *gin.Context) {
	var request model.CreateProductRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		helper.ReturnErrorBadRequest(c, err, nil)
		return
	}

	err = controller.ProductService.Create(c.Request.Context(), request)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "created", nil)
}

func (controller *ProductController) Update(c *gin.Context) {
	var request model.UpdateProductRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		helper.ReturnErrorBadRequest(c, err, nil)
		return
	}

	params := c.Param("code")

	request.Code = params
	err = controller.ProductService.Update(c.Request.Context(), request)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "updated", nil)
}

func (controller *ProductController) Delete(c *gin.Context) {
	params := c.Param("code")

	err := controller.ProductService.Delete(c.Request.Context(), params)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "deleted", nil)
}
