package controller

import (
	"apriori/api/middleware"
	"apriori/api/response"
	"apriori/model"
	"apriori/service"
	"apriori/utils"
	"github.com/gin-gonic/gin"
)

type commentController struct {
	CommentService service.CommentService
}

func NewCommentController(commentService *service.CommentService) *commentController {
	return &commentController{
		CommentService: *commentService,
	}
}

func (controller *commentController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.POST("/comments", controller.Create)
	}

	unauthorized := router.Group("/api")
	{
		unauthorized.GET("/comments/:comment_id", controller.FindById)
		unauthorized.GET("/comments/product/:product_code", controller.FindAllByProductCode)
		unauthorized.GET("/comments/user-order/:user_order_id", controller.FindByUserOrderId)
	}

	return router
}

func (controller *commentController) FindAllByProductCode(c *gin.Context) {
	productCode := c.Param("product_code")
	comments, err := controller.CommentService.FindAllByProductCode(c.Request.Context(), productCode)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", comments)
}

func (controller *commentController) FindById(c *gin.Context) {
	commentId := utils.StrToInt(c.Param("comment_id"))
	comment, err := controller.CommentService.FindById(c.Request.Context(), commentId)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", comment)
}

func (controller *commentController) FindByUserOrderId(c *gin.Context) {
	userOrderId := utils.StrToInt(c.Param("user_order_id"))
	comment, err := controller.CommentService.FindByUserOrderId(c.Request.Context(), userOrderId)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", comment)
}

func (controller *commentController) Create(c *gin.Context) {
	var request model.CreateCommentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	comment, err := controller.CommentService.Create(c.Request.Context(), request)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", comment)
}
