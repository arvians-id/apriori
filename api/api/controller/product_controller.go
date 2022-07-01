package controller

import (
	"apriori/api/response"
	"apriori/model"
	"apriori/service"
	"apriori/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

type ProductController struct {
	ProductService service.ProductService
	StorageService service.StorageService
}

func NewProductController(productService *service.ProductService, storageService *service.StorageService) *ProductController {
	return &ProductController{
		ProductService: *productService,
		StorageService: *storageService,
	}
}

func (controller *ProductController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api")
	{
		authorized.GET("/products", controller.FindAll)
		authorized.GET("/products/:code", controller.FindById)
		authorized.POST("/products", controller.Create)
		authorized.PATCH("/products/:code", controller.Update)
		authorized.DELETE("/products/:code", controller.Delete)
		authorized.GET("/products/:code/recommendation", controller.FindAllRecommendation)
	}

	return router
}

func (controller *ProductController) FindAll(c *gin.Context) {
	products, err := controller.ProductService.FindAll(c.Request.Context())
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", products)
}

func (controller *ProductController) FindAllRecommendation(c *gin.Context) {
	params := c.Param("code")

	products, err := controller.ProductService.FindAllRecommendation(c.Request.Context(), params)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", products)
}

func (controller *ProductController) FindById(c *gin.Context) {
	params := c.Param("code")

	product, err := controller.ProductService.FindByCode(c.Request.Context(), params)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", product)
}

func (controller *ProductController) Create(c *gin.Context) {
	var request model.CreateProductRequest
	request.Code = c.PostForm("code")
	request.Name = c.PostForm("name")
	request.Description = c.PostForm("description")
	request.Price = utils.StrToInt(c.PostForm("price"))

	file, header, err := c.Request.FormFile("image")
	filePath := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", os.Getenv("AWS_BUCKET"), os.Getenv("AWS_REGION"), "no-image.png")
	if err == nil {
		pathName, err := controller.StorageService.UploadFileS3(file, header)
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}
		filePath = pathName
	}

	request.Image = filePath
	product, err := controller.ProductService.Create(c.Request.Context(), request)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "created", product)
}

func (controller *ProductController) Update(c *gin.Context) {
	var request model.UpdateProductRequest
	request.Code = c.PostForm("code")
	request.Name = c.PostForm("name")
	request.Description = c.PostForm("description")
	request.Price = utils.StrToInt(c.PostForm("price"))
	params := c.Param("code")

	file, header, err := c.Request.FormFile("image")
	if err == nil {
		pathName, err := controller.StorageService.UploadFileS3(file, header)
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		request.Image = pathName
	}

	request.Code = params
	product, err := controller.ProductService.Update(c.Request.Context(), request)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "updated", product)
}

func (controller *ProductController) Delete(c *gin.Context) {
	params := c.Param("code")

	err := controller.ProductService.Delete(c.Request.Context(), params)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "deleted", nil)
}
