package route

import (
	"apriori/api/controller"
	"apriori/api/middleware"
	"apriori/cache"
	"apriori/config"
	"apriori/repository"
	"apriori/service"
	"database/sql"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func NewInitializedServer(configuration config.Config) (*gin.Engine, *sql.DB) {
	// Setup Configuration
	router := gin.Default()
	db, err := config.NewPostgreSQL(configuration)
	if err != nil {
		log.Fatal(err)
	}

	// Setup Repository
	userRepository := repository.NewUserRepository()
	passwordRepository := repository.NewPasswordResetRepository()
	authRepository := repository.NewAuthRepository()
	productRepository := repository.NewProductRepository()
	transactionRepository := repository.NewTransactionRepository()
	aprioriRepository := repository.NewAprioriRepository()
	paymentRepository := repository.NewPaymentRepository()
	userOrderRepository := repository.NewUserOrderRepository()

	// Setup Service
	storageService := service.NewStorageService(configuration)
	userService := service.NewUserService(&userRepository, db)
	authService := service.NewAuthService(&userRepository, &authRepository, db)
	jwtService := service.NewJwtService()
	emailService := service.NewEmailService()
	passwordResetService := service.NewPasswordResetService(&passwordRepository, &userRepository, db)
	productService := service.NewProductService(&productRepository, storageService, &aprioriRepository, db)
	transactionService := service.NewTransactionService(&transactionRepository, &productRepository, db)
	aprioriService := service.NewAprioriService(&transactionRepository, storageService, &productRepository, &aprioriRepository, db)
	paymentService := service.NewPaymentService(configuration, &paymentRepository, &userOrderRepository, &transactionRepository, db)
	userOrderService := service.NewUserOrderService(&paymentRepository, &userOrderRepository, db)

	// Setup Cache
	transactionCache := cache.NewTransactionCache(configuration)
	aprioriCache := cache.NewAprioriCache(configuration)
	productCache := cache.NewProductCache(configuration)
	paymentCache := cache.NewPaymentCache(configuration)
	userOrderCache := cache.NewUserOrderCache(configuration)

	// Setup Controller
	userController := controller.NewUserController(&userService)
	authController := controller.NewAuthController(&authService, &userService, jwtService, emailService, &passwordResetService)
	productController := controller.NewProductController(&productService, &storageService, &productCache)
	transactionController := controller.NewTransactionController(&transactionService, &storageService, &transactionCache)
	aprioriController := controller.NewAprioriController(aprioriService, &storageService, &aprioriCache)
	paymentController := controller.NewPaymentController(&paymentService, &userOrderService, emailService, &paymentCache, &userOrderCache)
	userOrderController := controller.NewUserOrderController(&paymentService, &userOrderService, &userOrderCache, &paymentCache)

	// CORS Middleware
	router.Use(middleware.SetupCorsMiddleware())

	// Main Router
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Apriori Algorithm API. Created By https://github.com/arvians-id",
		})
	})

	paymentController.Route(router)

	// X API KEY Middleware
	router.Use(middleware.SetupXApiKeyMiddleware())

	// Setup Router
	authController.Route(router)
	userController.Route(router)
	productController.Route(router)
	transactionController.Route(router)
	aprioriController.Route(router)
	userOrderController.Route(router)

	return router, db
}
