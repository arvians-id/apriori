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
		authorized.POST("/products", controller.Create)
		authorized.PATCH("/products/:code", controller.Update)
		authorized.DELETE("/products/:code", controller.Delete)
	}

	notAuthorized := router.Group("/api")
	{
		notAuthorized.GET("/products", controller.FindAll)
		notAuthorized.GET("/products/:code/recommendation", controller.FindAllRecommendation)
		notAuthorized.GET("/products/:code", controller.FindById)
	}

	return router
}

func (controller *ProductController) FindAll(c *gin.Context) {
	search := strings.ToLower(c.Query("search"))
	searchKey := search
	if search == "" {
		searchKey = "all"
	}
	key := fmt.Sprintf("%s-product", searchKey)
	productsCache, err := controller.CacheService.Get(c, key)
	if err == redis.Nil {
		products, err := controller.ProductService.FindAll(c.Request.Context(), search)
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		err = controller.CacheService.Set(c.Request.Context(), key, products)
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		response.ReturnSuccessOK(c, "OK", products)
		return
	} else if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	var product []model.GetProductResponse
	err = json.Unmarshal(bytes.NewBufferString(productsCache).Bytes(), &product)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", product)
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

	// delete previous cache
	_ = controller.CacheService.Del(c.Request.Context(), "all-product")

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

	// delete previous cache
	_ = controller.CacheService.Del(c.Request.Context(), "all-product", fmt.Sprintf("product-%s", product.Code))

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
	_ = controller.CacheService.Del(c.Request.Context(), "all-product")

	response.ReturnSuccessOK(c, "deleted", nil)
}
