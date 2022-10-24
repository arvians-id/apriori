package server

import (
	"database/sql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/arvians-id/apriori/cmd/config"
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
	storageService := service.NewStorageService(configuration)
	userService := cache.NewUserCacheService(&userRepository, redisLibrary, db)
	jwtService := service.NewJwtService()
	emailService := service.NewEmailService()
	notificationService := service.NewNotificationService(&notificationRepository, &userRepository, &emailService, db)
	passwordResetService := service.NewPasswordResetService(&passwordRepository, &userRepository, db)
	productService := cache.NewProductCacheService(&productRepository, &storageService, &aprioriRepository, redisLibrary, db)
	transactionService := service.NewTransactionService(&transactionRepository, &productRepository, db)
	aprioriService := service.NewAprioriService(&transactionRepository, storageService, &productRepository, &aprioriRepository, db)
	paymentService := service.NewPaymentService(configuration, &paymentRepository, &userOrderRepository, &transactionRepository, &notificationService, db)
	userOrderService := service.NewUserOrderService(&paymentRepository, &userOrderRepository, &userRepository, db)
	categoryService := cache.NewCategoryCacheService(&categoryRepository, redisLibrary, db)
	commentService := service.NewCommentService(&commentRepository, &productRepository, db)

	// Setup Controller
	userController := rest.NewUserController(&userService)
	authController := rest.NewAuthController(&userService, &jwtService, &emailService, &passwordResetService)
	productController := rest.NewProductController(&productService, &storageService)
	transactionController := rest.NewTransactionController(&transactionService, &storageService)
	aprioriController := rest.NewAprioriController(aprioriService, &storageService)
	paymentController := rest.NewPaymentController(&paymentService, &userOrderService, &emailService, &notificationService)
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
		emailService,
		jwtService,
		notificationService,
		passwordResetService,
		paymentService,
		productService,
		storageService,
		transactionService,
		userOrderService,
		userService,
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
