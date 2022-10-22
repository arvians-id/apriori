package rest

import (
	"github.com/arvians-id/apriori/internal/http/middleware"
	"github.com/arvians-id/apriori/internal/http/presenter/request"
	"github.com/arvians-id/apriori/internal/http/presenter/response"
	"github.com/arvians-id/apriori/internal/service"
	"github.com/arvians-id/apriori/util"
	"github.com/gin-gonic/gin"
)

type CommentController struct {
	CommentService service.CommentService
}

func NewCommentController(commentService *service.CommentService) *CommentController {
	return &CommentController{
		CommentService: *commentService,
	}
}

func (controller *CommentController) Route(router *gin.Engine) *gin.Engine {
	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.POST("/comments", controller.Create)
	}

	unauthorized := router.Group("/api")
	{
		unauthorized.GET("/comments/:id", controller.FindById)
		unauthorized.GET("/comments/rating/:product_code", controller.FindAllRatingByProductCode)
		unauthorized.GET("/comments/product/:product_code", controller.FindAllByProductCode)
		unauthorized.GET("/comments/user-order/:user_order_id", controller.FindByUserOrderId)
	}

	return router
}

func (controller *CommentController) FindAllByProductCode(c *gin.Context) {
	productCodeParam := c.Param("product_code")
	tagsQuery := c.Query("tags")
	ratingQuery := c.Query("rating")
	comments, err := controller.CommentService.FindAllByProductCode(c.Request.Context(), productCodeParam, ratingQuery, tagsQuery)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", comments)
}

func (controller *CommentController) FindById(c *gin.Context) {
	idParam := util.StrToInt(c.Param("id"))
	comment, err := controller.CommentService.FindById(c.Request.Context(), idParam)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", comment)
}

func (controller *CommentController) FindAllRatingByProductCode(c *gin.Context) {
	productCodeParam := c.Param("product_code")
	comments, err := controller.CommentService.FindAllRatingByProductCode(c.Request.Context(), productCodeParam)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", comments)
}

func (controller *CommentController) FindByUserOrderId(c *gin.Context) {
	userOrderIdParam := util.StrToInt(c.Param("user_order_id"))
	comment, err := controller.CommentService.FindByUserOrderId(c.Request.Context(), userOrderIdParam)
	if err != nil {
		if err.Error() == response.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", comment)
}

func (controller *CommentController) Create(c *gin.Context) {
	var requestCreate request.CreateCommentRequest
	if err := c.ShouldBindJSON(&requestCreate); err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	comment, err := controller.CommentService.Create(c.Request.Context(), &requestCreate)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", comment)
}
