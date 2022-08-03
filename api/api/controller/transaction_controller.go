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
		authorized.GET("/transactions/:code", controller.FindByTransaction)
		authorized.POST("/transactions", controller.Create)
		authorized.POST("/transactions/csv", controller.CreateFromCsv)
		authorized.PATCH("/transactions/:numberTransaction", controller.Update)
		authorized.DELETE("/transactions/:numberTransaction", controller.Delete)
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

	var transaction []model.GetTransactionResponse
	err = json.Unmarshal(bytes.NewBufferString(transactionCache).Bytes(), &transaction)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", transaction)
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

	transaction, err := controller.TransactionService.Create(c.Request.Context(), request)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	_ = controller.CacheService.Del(c.Request.Context(), "all-transaction")

	response.ReturnSuccessOK(c, "created", transaction)
}

func (controller *TransactionController) CreateFromCsv(c *gin.Context) {
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

	err = controller.TransactionService.CreateFromCsv(c.Request.Context(), data)
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

	noTransaction := c.Param("numberTransaction")

	request.NoTransaction = noTransaction
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
	noTransaction := c.Param("numberTransaction")
	err := controller.TransactionService.Delete(c.Request.Context(), noTransaction)
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
