package server

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/arvians-id/apriori/cmd/library/cache"
	directive2 "github.com/arvians-id/apriori/internal/http/controller/graph/directive"
	"github.com/arvians-id/apriori/internal/http/controller/graph/generated"
	"github.com/arvians-id/apriori/internal/http/controller/graph/resolver"
	"github.com/arvians-id/apriori/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewInitializedMainRoute(
	router *gin.Engine,
	aprioriService service.AprioriService,
	redisLibrary *cache.Redis,
	categoryService service.CategoryService,
	commentService service.CommentService,
	emailService service.EmailService,
	jwtService service.JwtService,
	notificationService service.NotificationService,
	passwordResetService service.PasswordResetService,
	paymentService service.PaymentService,
	productService service.ProductService,
	storageService service.StorageService,
	transactionService service.TransactionService,
	userOrderService service.UserOrderService,
	userService service.UserService,
) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Apriori Algorithm API. Created By https://github.com/arvians-id",
		})
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// GraphQL Route
	router.GET("/playground", func(c *gin.Context) {
		h := playground.Handler("GraphQL", "/query")
		h.ServeHTTP(c.Writer, c.Request)
	})

	router.POST("/query", func(c *gin.Context) {
		generatedConfig := generated.Config{
			Resolvers: &resolver.Resolver{
				AprioriService:       aprioriService,
				Redis:                redisLibrary,
				CategoryService:      categoryService,
				CommentService:       commentService,
				EmailService:         emailService,
				JwtService:           jwtService,
				NotificationService:  notificationService,
				PasswordResetService: passwordResetService,
				PaymentService:       paymentService,
				ProductService:       productService,
				StorageService:       storageService,
				TransactionService:   transactionService,
				UserOrderService:     userOrderService,
				UserService:          userService,
			},
		}
		// Schema directives
		generatedConfig.Directives.Binding = directive2.Binding
		generatedConfig.Directives.ApiKey = directive2.ApiKey
		generatedConfig.Directives.HasRole = directive2.HasRoles
		h := handler.NewDefaultServer(generated.NewExecutableSchema(generatedConfig))
		h.ServeHTTP(c.Writer, c.Request)
	})
}
