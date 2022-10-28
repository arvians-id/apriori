package rest

import (
	"github.com/arvians-id/apriori/cmd/library/aws"
	"github.com/arvians-id/apriori/internal/http/middleware"
	"github.com/arvians-id/apriori/internal/http/presenter/request"
	"github.com/arvians-id/apriori/internal/http/presenter/response"
	"github.com/arvians-id/apriori/internal/service"
	"github.com/arvians-id/apriori/util"
	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	TransactionService service.TransactionService
	StorageS3          aws.StorageS3
}

func NewTransactionController(transactionService *service.TransactionService, storageS3 *aws.StorageS3) *TransactionController {
	return &TransactionController{
		TransactionService: *transactionService,
		StorageS3:          *storageS3,
	}
}

func (controller *TransactionController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.GET("/transactions", controller.FindAll)
		authorized.GET("/transactions/:number_transaction", controller.FindByNoTransaction)
		authorized.POST("/transactions", controller.Create)
		authorized.POST("/transactions/csv", controller.CreateByCSV)
		authorized.PATCH("/transactions/:number_transaction", controller.Update)
		authorized.DELETE("/transactions/:number_transaction", controller.Delete)
		authorized.DELETE("/transactions/truncate", controller.Truncate)
	}

	return router
}

func (controller *TransactionController) FindAll(c *gin.Context) {
	transaction, err := controller.TransactionService.FindAll(c.Request.Context())
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", transaction)
}

func (controller *TransactionController) FindByNoTransaction(c *gin.Context) {
	noTransactionParam := c.Param("number_transaction")
	transactions, err := controller.TransactionService.FindByNoTransaction(c.Request.Context(), noTransactionParam)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", transactions)
}

func (controller *TransactionController) Create(c *gin.Context) {
	var requestCreate request.CreateTransactionRequest
	err := c.ShouldBindJSON(&requestCreate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	transaction, err := controller.TransactionService.Create(c.Request.Context(), &requestCreate)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "created", transaction)
}

func (controller *TransactionController) CreateByCSV(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	filePath, fileName := util.GetPathAWS(header.Filename)
	err = controller.StorageS3.UploadToAWS(file, fileName, header.Header.Get("Content-Type"))
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	data, err := util.OpenCsvFile(filePath)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	err = controller.TransactionService.CreateByCsv(c.Request.Context(), data)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	defer func(StorageS3 *aws.StorageS3, filePath string) {
		_ = StorageS3.DeleteFromAWS(filePath)
	}(&controller.StorageS3, filePath)

	response.ReturnSuccessOK(c, "created", nil)
}

func (controller *TransactionController) Update(c *gin.Context) {
	var requestUpdate request.UpdateTransactionRequest
	err := c.ShouldBindJSON(&requestUpdate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	noTransactionParam := c.Param("number_transaction")
	requestUpdate.NoTransaction = noTransactionParam
	transaction, err := controller.TransactionService.Update(c.Request.Context(), &requestUpdate)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "updated", transaction)
}

func (controller *TransactionController) Delete(c *gin.Context) {
	noTransactionParam := c.Param("number_transaction")
	err := controller.TransactionService.Delete(c.Request.Context(), noTransactionParam)
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

func (controller *TransactionController) Truncate(c *gin.Context) {
	err := controller.TransactionService.Truncate(c.Request.Context())
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "deleted", nil)
}
