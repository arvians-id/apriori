package controller

import (
	"apriori/app/middleware"
	"apriori/app/response"
	"apriori/model"
	"apriori/service"
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

func NewProductController(
	productService *service.ProductService,
	storageService *service.StorageService,
	cacheService *service.CacheService,
) *ProductController {
	return &ProductController{
		ProductService: *productService,
		StorageService: *storageService,
		CacheService:   *cacheService,
	}
}

func (controller *ProductController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.GET("/products-admin", controller.FindAllByAdmin)
		authorized.POST("/products", controller.Create)
		authorized.PATCH("/products/:code", controller.Update)
		authorized.DELETE("/products/:code", controller.Delete)
	}

	unauthorized := router.Group("/api")
	{
		unauthorized.GET("/products", controller.FindAllByUser)
		unauthorized.GET("/products/:code/category", controller.FindAllSimilarCategory)
		unauthorized.GET("/products/:code/recommendation", controller.FindAllRecommendation)
		unauthorized.GET("/products/:code", controller.FindByCode)
	}

	return router
}

func (controller *ProductController) FindAllByAdmin(c *gin.Context) {
	products, err := controller.ProductService.FindAllByAdmin(c.Request.Context())
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", products)
}

func (controller *ProductController) FindAllSimilarCategory(c *gin.Context) {
	codeParam := c.Param("code")
	products, err := controller.ProductService.FindAllBySimilarCategory(c.Request.Context(), codeParam)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", products)
}

func (controller *ProductController) FindAllByUser(c *gin.Context) {
	searchQuery := strings.ToLower(c.Query("search"))
	categoryQuery := strings.ToLower(c.Query("category"))
	products, err := controller.ProductService.FindAll(c.Request.Context(), searchQuery, categoryQuery)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", products)
}

func (controller *ProductController) FindAllRecommendation(c *gin.Context) {
	codeParam := c.Param("code")
	products, err := controller.ProductService.FindAllRecommendation(c.Request.Context(), codeParam)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", products)
}

func (controller *ProductController) FindByCode(c *gin.Context) {
	codeParam := c.Param("code")
	key := fmt.Sprintf("product-%s", codeParam)
	productCache, err := controller.CacheService.Get(c, key)
	if err == redis.Nil {
		product, err := controller.ProductService.FindByCode(c.Request.Context(), codeParam)
		if err != nil {
			if err.Error() == response.ErrorNotFound {
				response.ReturnErrorNotFound(c, err, nil)
				return
			}

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

	var productCacheResponse model.GetProductResponse
	err = json.Unmarshal(bytes.NewBufferString(productCache).Bytes(), &productCacheResponse)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", productCacheResponse)
}

func (controller *ProductController) Create(c *gin.Context) {
	var request model.CreateProductRequest
	err := c.ShouldBind(&request)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

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
	product, err := controller.ProductService.Create(c.Request.Context(), &request)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "created", product)
}

func (controller *ProductController) Update(c *gin.Context) {
	var request model.UpdateProductRequest
	err := c.ShouldBind(&request)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err == nil {
		pathName, err := controller.StorageService.UploadFileS3(file, header)
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		request.Image = pathName
	}

	request.Code = c.Param("code")
	product, err := controller.ProductService.Update(c.Request.Context(), &request)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	_ = controller.CacheService.Del(c.Request.Context(), fmt.Sprintf("product-%s", product.Code))

	response.ReturnSuccessOK(c, "updated", product)
}

func (controller *ProductController) Delete(c *gin.Context) {
	codeParam := c.Param("code")
	err := controller.ProductService.Delete(c.Request.Context(), codeParam)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	_ = controller.CacheService.Del(c.Request.Context(), fmt.Sprintf("product-%s", codeParam))

	response.ReturnSuccessOK(c, "deleted", nil)
}
