package controller

import (
	"apriori/api/middleware"
	"apriori/api/response"
	"apriori/model"
	"apriori/service"
	"apriori/utils"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"sync"
)

type TransactionController struct {
	TransactionService service.TransactionService
	StorageService     service.StorageService
	CacheService       service.CacheService
}

func NewTransactionController(transactionService *service.TransactionService, storageService *service.StorageService, cacheService *service.CacheService) *TransactionController {
	return &TransactionController{
		TransactionService: *transactionService,
		StorageService:     *storageService,
		CacheService:       *cacheService,
	}
}

func (controller *TransactionController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.GET("/transactions", controller.FindAll)
		authorized.GET("/transactions/:number_transaction", controller.FindByNoTransaction)
		authorized.POST("/transactions", controller.Create)
		authorized.POST("/transactions/csv", controller.CreateByCsv)
		authorized.PATCH("/transactions/:number_transaction", controller.Update)
		authorized.DELETE("/transactions/:number_transaction", controller.Delete)
		authorized.DELETE("/transactions/truncate", controller.Truncate)
	}

	return router
}

func (controller *TransactionController) FindAll(c *gin.Context) {
	transactionCache, err := controller.CacheService.Get(c, "all-transaction")
	if err == redis.Nil {
		transaction, err := controller.TransactionService.FindAll(c.Request.Context())
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		err = controller.CacheService.Set(c.Request.Context(), "all-transaction", transaction)
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		response.ReturnSuccessOK(c, "OK", transaction)
		return
	} else if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	var transactionCacheResponses []model.GetTransactionResponse
	err = json.Unmarshal(bytes.NewBufferString(transactionCache).Bytes(), &transactionCacheResponses)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", transactionCacheResponses)
}

func (controller *TransactionController) FindByNoTransaction(c *gin.Context) {
	noTransactionParam := c.Param("number_transaction")
	transactions, err := controller.TransactionService.FindByNoTransaction(c.Request.Context(), noTransactionParam)
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

	transaction, err := controller.TransactionService.Create(c.Request.Context(), request)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	_ = controller.CacheService.Del(c.Request.Context(), "all-transaction")

	response.ReturnSuccessOK(c, "created", transaction)
}

func (controller *TransactionController) CreateByCsv(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	var wg sync.WaitGroup
	pathName, err := controller.StorageService.WaitUploadFileS3(file, header, &wg)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}
	wg.Wait()

	data, err := utils.OpenCsvFile(pathName)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	err = controller.TransactionService.CreateByCsv(c.Request.Context(), data)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	_ = controller.CacheService.Del(c.Request.Context(), "all-transaction")

	response.ReturnSuccessOK(c, "created", nil)
}

func (controller *TransactionController) Update(c *gin.Context) {
	var request model.UpdateTransactionRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	noTransactionParam := c.Param("number_transaction")

	request.NoTransaction = noTransactionParam
	transaction, err := controller.TransactionService.Update(c.Request.Context(), request)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	_ = controller.CacheService.Del(c.Request.Context(), "all-transaction")

	response.ReturnSuccessOK(c, "updated", transaction)
}

func (controller *TransactionController) Delete(c *gin.Context) {
	noTransactionParam := c.Param("number_transaction")
	err := controller.TransactionService.Delete(c.Request.Context(), noTransactionParam)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	_ = controller.CacheService.Del(c.Request.Context(), "all-transaction")

	response.ReturnSuccessOK(c, "deleted", nil)
}

func (controller *TransactionController) Truncate(c *gin.Context) {
	err := controller.TransactionService.Truncate(c.Request.Context())
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	_ = controller.CacheService.Del(c.Request.Context(), "all-transaction")

	response.ReturnSuccessOK(c, "deleted", nil)
}
