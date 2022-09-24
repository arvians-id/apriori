package route

import (
	"database/sql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/arvians-id/apriori/config"
	"github.com/arvians-id/apriori/http/controller/graph/directive"
	"github.com/arvians-id/apriori/http/controller/graph/generated"
	"github.com/arvians-id/apriori/http/controller/graph/resolver"
	"github.com/arvians-id/apriori/http/controller/rest"
	"github.com/arvians-id/apriori/http/middleware"
	repository "github.com/arvians-id/apriori/repository/postgres"
	"github.com/arvians-id/apriori/service"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
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
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// Setup Configuration
	router := gin.Default()

	// Setup Database
	db, err := NewInitializedDatabase(configuration)
	if err != nil {
		log.Fatal(err)
	}

	// Setup Repository
	userRepository := repository.NewUserRepository()
	passwordRepository := repository.NewPasswordResetRepository()
	productRepository := repository.NewProductRepository()
	transactionRepository := repository.NewTransactionRepository()
	aprioriRepository := repository.NewAprioriRepository()
	paymentRepository := repository.NewPaymentRepository()
	userOrderRepository := repository.NewUserOrderRepository()
	categoryRepository := repository.NewCategoryRepository()
	commentRepository := repository.NewCommentRepository()
	notificationRepository := repository.NewNotificationRepository()

	// Setup Service
	storageService := service.NewStorageService(configuration)
	userService := service.NewUserService(&userRepository, db)
	jwtService := service.NewJwtService()
	emailService := service.NewEmailService()
	notificationService := service.NewNotificationService(&notificationRepository, &userRepository, &emailService, db)
	passwordResetService := service.NewPasswordResetService(&passwordRepository, &userRepository, db)
	productService := service.NewProductService(&productRepository, &storageService, &aprioriRepository, db)
	transactionService := service.NewTransactionService(&transactionRepository, &productRepository, db)
	aprioriService := service.NewAprioriService(&transactionRepository, storageService, &productRepository, &aprioriRepository, db)
	paymentService := service.NewPaymentService(configuration, &paymentRepository, &userOrderRepository, &transactionRepository, &notificationService, db)
	userOrderService := service.NewUserOrderService(&paymentRepository, &userOrderRepository, &userRepository, db)
	cacheService := service.NewCacheService(configuration)
	categoryService := service.NewCategoryService(&categoryRepository, db)
	commentService := service.NewCommentService(&commentRepository, &productRepository, db)

	// Setup Controller
	userController := rest.NewUserController(&userService)
	authController := rest.NewAuthController(&userService, &jwtService, &emailService, &passwordResetService)
	productController := rest.NewProductController(&productService, &storageService, &cacheService)
	transactionController := rest.NewTransactionController(&transactionService, &storageService, &cacheService)
	aprioriController := rest.NewAprioriController(aprioriService, &storageService, &cacheService)
	paymentController := rest.NewPaymentController(&paymentService, &userOrderService, &emailService, &cacheService, &notificationService)
	userOrderController := rest.NewUserOrderController(&paymentService, &userOrderService, &cacheService)
	categoryController := rest.NewCategoryController(&categoryService, &cacheService)
	commentController := rest.NewCommentController(&commentService)
	rajaOngkirController := rest.NewRajaOngkirController()
	notificationController := rest.NewNotificationController(&notificationService)

	// CORS Middleware
	router.Use(middleware.SetupCorsMiddleware())
	router.Use(middleware.GinContextToContextMiddleware())

	// Main Route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Apriori Algorithm API. Created By https://github.com/arvians-id",
		})
	})

	// GraphQL Route
	router.GET("/playground", func(c *gin.Context) {
		h := playground.Handler("GraphQL", "/query")
		h.ServeHTTP(c.Writer, c.Request)
	})

	router.POST("/query", func(c *gin.Context) {
		generatedConfig := generated.Config{
			Resolvers: &resolver.Resolver{
				AprioriService:       aprioriService,
				CacheService:         cacheService,
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
		generatedConfig.Directives.Binding = directive.Binding
		generatedConfig.Directives.ApiKey = directive.ApiKey
		generatedConfig.Directives.HasRole = directive.HasRoles
		h := handler.NewDefaultServer(generated.NewExecutableSchema(generatedConfig))
		h.ServeHTTP(c.Writer, c.Request)
	})

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
