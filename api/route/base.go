package route

import (
	"apriori/api/controller"
	"apriori/api/middleware"
	"apriori/config"
	"apriori/database/seeder"
	repository "apriori/repository/postgres"
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
	cacheService := service.NewCacheService(configuration)

	// Setup Controller
	userController := controller.NewUserController(&userService)
	authController := controller.NewAuthController(&authService, &userService, jwtService, emailService, &passwordResetService)
	productController := controller.NewProductController(&productService, &storageService, &cacheService)
	transactionController := controller.NewTransactionController(&transactionService, &storageService, &cacheService)
	aprioriController := controller.NewAprioriController(aprioriService, &storageService, &cacheService)
	paymentController := controller.NewPaymentController(&paymentService, &userOrderService, emailService, &cacheService)
	userOrderController := controller.NewUserOrderController(&paymentService, &userOrderService, &cacheService)

	// CORS Middleware
	router.Use(middleware.SetupCorsMiddleware())

	// Main Route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Apriori Algorithm API. Created By https://github.com/arvians-id",
		})
	})
	router.GET("/api/seeds", middleware.SetupXApiKeyMiddleware(), func(c *gin.Context) {
		// Seeder
		seeder.RegisterSeeder(productService)
		c.JSON(200, gin.H{
			"message": "Seeding Success",
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
