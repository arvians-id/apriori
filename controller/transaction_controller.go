package controller

import (
	"apriori/helper"
	"apriori/middleware"
	"apriori/model"
	"apriori/service"
	"errors"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
)

type TransactionController struct {
	TransactionService service.TransactionService
}

func NewTransactionController(transactionService *service.TransactionService) *TransactionController {
	return &TransactionController{
		TransactionService: *transactionService,
	}
}

func (controller *TransactionController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.GET("/transactions", controller.FindAll)
		authorized.GET("/transactions/:code", controller.FindByTransaction)
		authorized.POST("/transactions", controller.Create)
		authorized.POST("/transactions/csv", controller.CreateFromCsv)
		authorized.PATCH("/transactions/:code", controller.Update)
		authorized.DELETE("/transactions/:code", controller.Delete)
	}

	return router
}

func (controller *TransactionController) FindAll(c *gin.Context) {
	transactions, err := controller.TransactionService.FindAll(c.Request.Context())
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "OK", transactions)
}

func (controller *TransactionController) FindByTransaction(c *gin.Context) {
	noTransaction := c.Param("code")
	transactions, err := controller.TransactionService.FindByTransaction(c.Request.Context(), noTransaction)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "OK", transactions)
}

func (controller *TransactionController) Create(c *gin.Context) {
	var request model.CreateTransactionRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		helper.ReturnErrorBadRequest(c, err, nil)
		return
	}

	err = controller.TransactionService.Create(c.Request.Context(), request)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "created", nil)
}

func (controller *TransactionController) CreateFromCsv(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		helper.ReturnErrorBadRequest(c, err, nil)
		return
	}
	if file.Header.Get("Content-Type") != "text/csv" {
		helper.ReturnErrorInternalServerError(c, errors.New("file not allowed"), nil)
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	path := "/assets/" + file.Filename
	path = filepath.Join(dir, path)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	data, err := helper.OpenCsvFile(path)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	err = controller.TransactionService.CreateFromCsv(c.Request.Context(), data)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "created", nil)
}

func (controller *TransactionController) Update(c *gin.Context) {
	var request model.UpdateTransactionRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		helper.ReturnErrorBadRequest(c, err, nil)
		return
	}

	noTransaction := c.Param("code")

	request.NoTransaction = noTransaction
	err = controller.TransactionService.Update(c.Request.Context(), request)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "updated", nil)
}

func (controller *TransactionController) Delete(c *gin.Context) {
	noTransaction := c.Param("code")
	err := controller.TransactionService.Delete(c.Request.Context(), noTransaction)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "deleted", nil)
}
