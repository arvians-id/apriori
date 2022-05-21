package controller

import (
	"apriori/helper"
	"apriori/middleware"
	"apriori/model"
	"apriori/service"
	"github.com/gin-gonic/gin"
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

	transactions, err := controller.TransactionService.Create(c.Request.Context(), request)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "OK", transactions)
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
	transactions, err := controller.TransactionService.Update(c.Request.Context(), request)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "OK", transactions)
}

func (controller *TransactionController) Delete(c *gin.Context) {
	noTransaction := c.Param("code")
	err := controller.TransactionService.Delete(c.Request.Context(), noTransaction)
	if err != nil {
		helper.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	helper.ReturnSuccessOK(c, "OK", nil)
}
