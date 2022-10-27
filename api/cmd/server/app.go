package server

import (
	"database/sql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/arvians-id/apriori/cmd/config"
	"github.com/arvians-id/apriori/cmd/library/auth"
	"github.com/arvians-id/apriori/cmd/library/aws"
	"github.com/arvians-id/apriori/cmd/library/messaging"
	"github.com/arvians-id/apriori/cmd/library/redis"
	directive2 "github.com/arvians-id/apriori/internal/http/controller/graph/directive"
	"github.com/arvians-id/apriori/internal/http/controller/graph/generated"
	"github.com/arvians-id/apriori/internal/http/controller/graph/resolver"
	"github.com/arvians-id/apriori/internal/http/controller/rest"
	"github.com/arvians-id/apriori/internal/http/middleware"
	"github.com/arvians-id/apriori/internal/repository/postgres"
	"github.com/arvians-id/apriori/internal/service"
	"github.com/arvians-id/apriori/internal/service/cache"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
)

func NewInitializedDatabase(configuration config.Config) (*sql.DB, error) {
	db, err := config.NewPostgreSQL(configuration)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewInitializedServer(configuration config.Config) (*gin.Engine, *sql.DB) {
	// Write log to file
	//f, _ := os.Create("gin.log")
	//gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// Setup Configuration
	router := gin.Default()

	// Setup Database
	db, err := NewInitializedDatabase(configuration)
	if err != nil {
		log.Fatal(err)
	}

	// Setup Library
	redisLibrary := redis.NewCacheService(configuration)
	storageLibrary := aws.NewStorageS3(configuration)
	jwtLibrary := auth.NewJsonWebToken()
	messagingProducerLibrary := messaging.NewProducer(messaging.ProducerConfig{
		NsqdAddress: "nsqd:4150",
	})

	// Setup Repository
	userRepository := postgres.NewUserRepository()
	passwordRepository := postgres.NewPasswordResetRepository()
	productRepository := postgres.NewProductRepository()
	transactionRepository := postgres.NewTransactionRepository()
	aprioriRepository := postgres.NewAprioriRepository()
	paymentRepository := postgres.NewPaymentRepository()
	userOrderRepository := postgres.NewUserOrderRepository()
	categoryRepository := postgres.NewCategoryRepository()
	commentRepository := postgres.NewCommentRepository()
	notificationRepository := postgres.NewNotificationRepository()

	// Setup Service
	userService := cache.NewUserCacheService(&userRepository, redisLibrary, db)
	notificationService := service.NewNotificationService(&notificationRepository, &userRepository, db)
	passwordResetService := service.NewPasswordResetService(&passwordRepository, &userRepository, db)
	productService := cache.NewProductCacheService(&productRepository, &aprioriRepository, redisLibrary, db)
	transactionService := service.NewTransactionService(&transactionRepository, &productRepository, db)
	aprioriService := service.NewAprioriService(&transactionRepository, &productRepository, &aprioriRepository, db)
	paymentService := service.NewPaymentService(configuration, &paymentRepository, &userOrderRepository, &transactionRepository, db)
	userOrderService := service.NewUserOrderService(&paymentRepository, &userOrderRepository, &userRepository, db)
	categoryService := cache.NewCategoryCacheService(&categoryRepository, redisLibrary, db)
	commentService := service.NewCommentService(&commentRepository, &productRepository, db)

	// Setup Controller
	userController := rest.NewUserController(&userService)
	authController := rest.NewAuthController(&userService, &passwordResetService, jwtLibrary, messagingProducerLibrary)
	productController := rest.NewProductController(&productService, storageLibrary, messagingProducerLibrary)
	transactionController := rest.NewTransactionController(&transactionService, storageLibrary)
	aprioriController := rest.NewAprioriController(aprioriService, storageLibrary)
	paymentController := rest.NewPaymentController(&paymentService, &userOrderService, &notificationService, &userService, messagingProducerLibrary)
	userOrderController := rest.NewUserOrderController(&paymentService, &userOrderService)
	categoryController := rest.NewCategoryController(&categoryService)
	commentController := rest.NewCommentController(&commentService)
	rajaOngkirController := rest.NewRajaOngkirController()
	notificationController := rest.NewNotificationController(&notificationService)

	// Main Middleware
	middleware.RegisterPrometheusMetrics()
	router.Use(middleware.SetupCorsMiddleware())
	router.Use(middleware.GinContextToContextMiddleware())
	router.Use(middleware.PrometheusMetricsMiddleware())

	// Main Route
	NewInitializedMainRoute(router,
		aprioriService,
		categoryService,
		commentService,
		notificationService,
		passwordResetService,
		paymentService,
		productService,
		transactionService,
		userOrderService,
		userService,
		storageLibrary,
		jwtLibrary,
		messagingProducerLibrary,
	)

	// REST API Route
	paymentController.Route(router)
	// X API KEY Middleware
	router.Use(middleware.SetupXApiKeyMiddleware())
	authController.Route(router)
	userController.Route(router)
	productController.Route(router)
	transactionController.Route(router)
	aprioriController.Route(router)
	userOrderController.Route(router)
	categoryController.Route(router)
	commentController.Route(router)
	rajaOngkirController.Route(router)
	notificationController.Route(router)

	return router, db
}

func NewInitializedMainRoute(
	router *gin.Engine,
	aprioriService service.AprioriService,
	categoryService service.CategoryService,
	commentService service.CommentService,
	notificationService service.NotificationService,
	passwordResetService service.PasswordResetService,
	paymentService service.PaymentService,
	productService service.ProductService,
	transactionService service.TransactionService,
	userOrderService service.UserOrderService,
	userService service.UserService,
	storageS3Library *aws.StorageS3,
	jwtLibrary *auth.JsonWebToken,
	messagingProducerLibrary *messaging.Producer,
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
				CategoryService:      categoryService,
				CommentService:       commentService,
				NotificationService:  notificationService,
				PasswordResetService: passwordResetService,
				PaymentService:       paymentService,
				ProductService:       productService,
				TransactionService:   transactionService,
				UserOrderService:     userOrderService,
				UserService:          userService,
				StorageS3:            storageS3Library,
				Jwt:                  jwtLibrary,
				Producer:             messagingProducerLibrary,
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
