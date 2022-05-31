package controller

import (
	"apriori/api/response"
	"apriori/model"
	"apriori/service"
	"apriori/utils"
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
	authorized := router.Group("/api")
	{
		authorized.GET("/transactions", controller.FindAll)
		authorized.GET("/transactions/:code", controller.FindByTransaction)
		authorized.POST("/transactions", controller.Create)
		authorized.POST("/transactions/csv", controller.CreateFromCsv)
		authorized.PATCH("/transactions/:numberTransaction", controller.Update)
		authorized.DELETE("/transactions/:numberTransaction", controller.Delete)
	}

	return router
}

func (controller *TransactionController) FindAll(c *gin.Context) {
	transactions, err := controller.TransactionService.FindAll(c.Request.Context())
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", transactions)
}

func (controller *TransactionController) FindByTransaction(c *gin.Context) {
	noTransaction := c.Param("code")
	transactions, err := controller.TransactionService.FindByTransaction(c.Request.Context(), noTransaction)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", transactions)
}

func (controller *TransactionController) Create(c *gin.Context) {
	var request model.CreateTransactionRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	err = controller.TransactionService.Create(c.Request.Context(), request)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "created", nil)
}

func (controller *TransactionController) CreateFromCsv(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}
	if file.Header.Get("Content-Type") != "text/csv" {
		response.ReturnErrorInternalServerError(c, errors.New("file not allowed"), nil)
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	path := "/assets/" + file.Filename
	path = filepath.Join(dir, path)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	data, err := utils.OpenCsvFile(path)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	err = controller.TransactionService.CreateFromCsv(c.Request.Context(), data)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "created", nil)
}

func (controller *TransactionController) Update(c *gin.Context) {
	var request model.UpdateTransactionRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	noTransaction := c.Param("numberTransaction")

	request.NoTransaction = noTransaction
	err = controller.TransactionService.Update(c.Request.Context(), request)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "updated", nil)
}

func (controller *TransactionController) Delete(c *gin.Context) {
	noTransaction := c.Param("numberTransaction")
	err := controller.TransactionService.Delete(c.Request.Context(), noTransaction)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "deleted", nil)
}
