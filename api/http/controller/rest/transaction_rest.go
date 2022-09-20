package rest

import (
	"encoding/json"
	"github.com/arvians-id/apriori/helper"
	"github.com/arvians-id/apriori/http/controller/rest/request"
	response2 "github.com/arvians-id/apriori/http/controller/rest/response"
	"github.com/arvians-id/apriori/http/middleware"
	"github.com/arvians-id/apriori/model"
	"github.com/arvians-id/apriori/service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"sync"
)

type TransactionController struct {
	TransactionService service.TransactionService
	StorageService     service.StorageService
	CacheService       service.CacheService
}

func NewTransactionController(
	transactionService *service.TransactionService,
	storageService *service.StorageService,
	cacheService *service.CacheService,
) *TransactionController {
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
			response2.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		err = controller.CacheService.Set(c.Request.Context(), "all-transaction", transaction)
		if err != nil {
			response2.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		response2.ReturnSuccessOK(c, "OK", transaction)
		return
	} else if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	var transactionCacheResponses []model.Transaction
	err = json.Unmarshal(transactionCache, &transactionCacheResponses)
	if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response2.ReturnSuccessOK(c, "OK", transactionCacheResponses)
}

func (controller *TransactionController) FindByNoTransaction(c *gin.Context) {
	noTransactionParam := c.Param("number_transaction")
	transactions, err := controller.TransactionService.FindByNoTransaction(c.Request.Context(), noTransactionParam)
	if err != nil {
		if err.Error() == response2.ErrorNotFound {
			response2.ReturnErrorNotFound(c, err, nil)
			return
		}

		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response2.ReturnSuccessOK(c, "OK", transactions)
}

func (controller *TransactionController) Create(c *gin.Context) {
	var requestCreate request.CreateTransactionRequest
	err := c.ShouldBindJSON(&requestCreate)
	if err != nil {
		response2.ReturnErrorBadRequest(c, err, nil)
		return
	}

	transaction, err := controller.TransactionService.Create(c.Request.Context(), &requestCreate)
	if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	_ = controller.CacheService.Del(c.Request.Context(), "all-transaction")

	response2.ReturnSuccessOK(c, "created", transaction)
}

func (controller *TransactionController) CreateByCsv(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response2.ReturnErrorBadRequest(c, err, nil)
		return
	}

	var wg sync.WaitGroup
	pathName, err := controller.StorageService.WaitUploadFileS3(file, header, &wg)
	if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}
	wg.Wait()

	data, err := helper.OpenCsvFile(pathName)
	if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	err = controller.TransactionService.CreateByCsv(c.Request.Context(), data)
	if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	_ = controller.CacheService.Del(c.Request.Context(), "all-transaction")

	response2.ReturnSuccessOK(c, "created", nil)
}

func (controller *TransactionController) Update(c *gin.Context) {
	var requestUpdate request.UpdateTransactionRequest
	err := c.ShouldBindJSON(&requestUpdate)
	if err != nil {
		response2.ReturnErrorBadRequest(c, err, nil)
		return
	}

	noTransactionParam := c.Param("number_transaction")
	requestUpdate.NoTransaction = noTransactionParam
	transaction, err := controller.TransactionService.Update(c.Request.Context(), &requestUpdate)
	if err != nil {
		if err.Error() == response2.ErrorNotFound {
			response2.ReturnErrorNotFound(c, err, nil)
			return
		}

		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	_ = controller.CacheService.Del(c.Request.Context(), "all-transaction")

	response2.ReturnSuccessOK(c, "updated", transaction)
}

func (controller *TransactionController) Delete(c *gin.Context) {
	noTransactionParam := c.Param("number_transaction")
	err := controller.TransactionService.Delete(c.Request.Context(), noTransactionParam)
	if err != nil {
		if err.Error() == response2.ErrorNotFound {
			response2.ReturnErrorNotFound(c, err, nil)
			return
		}

		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	_ = controller.CacheService.Del(c.Request.Context(), "all-transaction")

	response2.ReturnSuccessOK(c, "deleted", nil)
}

func (controller *TransactionController) Truncate(c *gin.Context) {
	err := controller.TransactionService.Truncate(c.Request.Context())
	if err != nil {
		response2.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// delete previous cache
	_ = controller.CacheService.Del(c.Request.Context(), "all-transaction")

	response2.ReturnSuccessOK(c, "deleted", nil)
}
