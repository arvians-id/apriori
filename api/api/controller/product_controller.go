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
	"os"
	"strings"
)

type ProductController struct {
	ProductService service.ProductService
	StorageService service.StorageService
	CacheService   service.CacheService
}

func NewProductController(productService *service.ProductService, storageService *service.StorageService, cacheService *service.CacheService) *ProductController {
	return &ProductController{
		ProductService: *productService,
		StorageService: *storageService,
		CacheService:   *cacheService,
	}
}

func (controller *ProductController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.GET("/products-admin", controller.FindAllOnAdmin)
		authorized.POST("/products", controller.Create)
		authorized.PATCH("/products/:code", controller.Update)
		authorized.DELETE("/products/:code", controller.Delete)
	}

	unauthorized := router.Group("/api")
	{
		unauthorized.GET("/products", controller.FindAll)
		unauthorized.GET("/products/:code/category", controller.FindAllSimilarCategory)
		unauthorized.GET("/products/:code/recommendation", controller.FindAllRecommendation)
		unauthorized.GET("/products/:code", controller.FindById)
	}

	return router
}

func (controller *ProductController) FindAllOnAdmin(c *gin.Context) {
	products, err := controller.ProductService.FindAllOnAdmin(c.Request.Context())
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", products)
}

func (controller *ProductController) FindAllSimilarCategory(c *gin.Context) {
	params := c.Param("code")
	products, err := controller.ProductService.FindAllSimilarCategory(c.Request.Context(), params)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", products)
}

func (controller *ProductController) FindAll(c *gin.Context) {
	search := strings.ToLower(c.Query("search"))
	category := strings.ToLower(c.Query("category"))
	products, err := controller.ProductService.FindAll(c.Request.Context(), search, category)
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

	key := fmt.Sprintf("product-%s", params)
	productCache, err := controller.CacheService.Get(c, key)
	if err == redis.Nil {
		product, err := controller.ProductService.FindByCode(c.Request.Context(), params)
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		err = controller.CacheService.Set(c.Request.Context(), key, product)
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		response.ReturnSuccessOK(c, "OK", product)
		return
	} else if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	var product model.GetProductResponse
	err = json.Unmarshal(bytes.NewBufferString(productCache).Bytes(), &product)
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
	request.Category = c.PostForm("category")
	request.Mass = utils.StrToInt(c.PostForm("mass"))

	file, header, err := c.Request.FormFile("image")
	filePath := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/assets/%s", os.Getenv("AWS_BUCKET"), os.Getenv("AWS_REGION"), "no-image.png")
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
	request.Code = c.Param("code")
	request.Name = c.PostForm("name")
	request.Description = c.PostForm("description")
	request.Price = utils.StrToInt(c.PostForm("price"))
	request.Category = c.PostForm("category")
	request.IsEmpty = utils.StrToInt(c.PostForm("is_empty"))
	request.Mass = utils.StrToInt(c.PostForm("mass"))
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

	// delete previous cache
	_ = controller.CacheService.Del(c.Request.Context(), fmt.Sprintf("product-%s", product.Code))

	response.ReturnSuccessOK(c, "updated", product)
}

func (controller *ProductController) Delete(c *gin.Context) {
	params := c.Param("code")

	err := controller.ProductService.Delete(c.Request.Context(), params)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	_ = controller.CacheService.Del(c.Request.Context(), fmt.Sprintf("product-%s", params))

	response.ReturnSuccessOK(c, "deleted", nil)
}
