package controller

import (
	"apriori/api/middleware"
	"apriori/api/response"
	"apriori/cache"
	"apriori/model"
	"apriori/service"
	"apriori/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"os"
)

type ProductController struct {
	ProductService service.ProductService
	StorageService service.StorageService
	ProductCache   cache.ProductCache
}

func NewProductController(productService *service.ProductService, storageService *service.StorageService, productCache *cache.ProductCache) *ProductController {
	return &ProductController{
		ProductService: *productService,
		StorageService: *storageService,
		ProductCache:   *productCache,
	}
}

func (controller *ProductController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api")
	{
		authorized.GET("/products", controller.FindAll)
		authorized.GET("/products/:code", controller.FindById)
		authorized.POST("/products", controller.Create, middleware.AuthJwtMiddleware())
		authorized.PATCH("/products/:code", controller.Update, middleware.AuthJwtMiddleware())
		authorized.DELETE("/products/:code", controller.Delete, middleware.AuthJwtMiddleware())
		authorized.GET("/products/:code/recommendation", controller.FindAllRecommendation)
	}

	return router
}

func (controller *ProductController) FindAll(c *gin.Context) {
	productsCache, err := controller.ProductCache.Get(c, "all-product")
	if err == redis.Nil {
		products, err := controller.ProductService.FindAll(c.Request.Context())
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		err = controller.ProductCache.Set(c.Request.Context(), "all-product", products)
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

	response.ReturnSuccessOK(c, "OK", productsCache)
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
	productCache, err := controller.ProductCache.SingleGet(c, key)
	if err == redis.Nil {
		product, err := controller.ProductService.FindByCode(c.Request.Context(), params)
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		err = controller.ProductCache.SingleSet(c.Request.Context(), key, product)
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

	response.ReturnSuccessOK(c, "OK", productCache)
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

	// Recover cache
	dataProduct, _ := controller.ProductService.FindAll(c.Request.Context())
	_ = controller.ProductCache.Set(c.Request.Context(), "all-product", dataProduct)

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

	// Recover cache
	dataProduct, _ := controller.ProductService.FindAll(c.Request.Context())
	singleProduct := fmt.Sprintf("product-%s", product.Code)
	_ = controller.ProductCache.Set(c.Request.Context(), "all-product", dataProduct)
	_ = controller.ProductCache.SingleSet(c.Request.Context(), singleProduct, product)

	response.ReturnSuccessOK(c, "updated", product)
}

func (controller *ProductController) Delete(c *gin.Context) {
	params := c.Param("code")

	err := controller.ProductService.Delete(c.Request.Context(), params)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// Recover cache
	dataProduct, _ := controller.ProductService.FindAll(c.Request.Context())
	_ = controller.ProductCache.Set(c.Request.Context(), "all-product", dataProduct)

	response.ReturnSuccessOK(c, "deleted", nil)
}
