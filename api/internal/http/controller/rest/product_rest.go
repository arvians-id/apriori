package rest

import (
	"fmt"
	"github.com/arvians-id/apriori/cmd/library/aws"
	"github.com/arvians-id/apriori/internal/http/middleware"
	"github.com/arvians-id/apriori/internal/http/presenter/request"
	"github.com/arvians-id/apriori/internal/http/presenter/response"
	"github.com/arvians-id/apriori/internal/service"
	"github.com/arvians-id/apriori/util"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strings"
)

type ProductController struct {
	ProductService service.ProductService
	StorageS3      aws.StorageS3
}

func NewProductController(productService *service.ProductService, storageS3 *aws.StorageS3) *ProductController {
	return &ProductController{
		ProductService: *productService,
		StorageS3:      *storageS3,
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
	product, err := controller.ProductService.FindByCode(c.Request.Context(), codeParam)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", product)
}

func (controller *ProductController) Create(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	filePath := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/assets/%s", os.Getenv("AWS_BUCKET"), os.Getenv("AWS_REGION"), "no-image.png")
	if err == nil {
		path, fileName := util.GetPathAWS(header.Filename)
		go func() {
			err = controller.StorageS3.UploadToAWS(file, fileName, header.Header.Get("Content-Type"))
			if err != nil {
				log.Println("[Product][Create][UploadFileS3Test] error upload file S3, err: ", err.Error())
			}
		}()
		filePath = path
	}

	var requestCreate request.CreateProductRequest
	err = c.ShouldBind(&requestCreate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	requestCreate.Image = filePath
	_, err = controller.ProductService.Create(c.Request.Context(), &requestCreate)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "created", nil)
}

func (controller *ProductController) Update(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	var filePath string
	if err == nil {
		path, fileName := util.GetPathAWS(header.Filename)
		go func() {
			err = controller.StorageS3.UploadToAWS(file, fileName, header.Header.Get("Content-Type"))
			if err != nil {
				log.Println("[Product][Create][UploadFileS3Test] error upload file S3, err: ", err.Error())
			}
		}()
		filePath = path
	}

	var requestUpdate request.UpdateProductRequest
	err = c.ShouldBind(&requestUpdate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	requestUpdate.Image = filePath
	requestUpdate.Code = c.Param("code")
	product, err := controller.ProductService.Update(c.Request.Context(), &requestUpdate)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

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

	response.ReturnSuccessOK(c, "deleted", nil)
}
