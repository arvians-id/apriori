package route

import (
	"apriori/app/controller"
	"apriori/app/middleware"
	"apriori/config"
	repository "apriori/repository/postgres"
	"apriori/service"
	"database/sql"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
)

func NewInitializedServer(configuration config.Config) (*gin.Engine, *sql.DB) {
	// Write log to file
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// Setup Configuration
	router := gin.Default()
	db, err := config.NewPostgreSQL(configuration)
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

	// Setup Service
	storageService := service.NewStorageService(configuration)
	userService := service.NewUserService(&userRepository, db)
	jwtService := service.NewJwtService()
	emailService := service.NewEmailService()
	passwordResetService := service.NewPasswordResetService(&passwordRepository, &userRepository, db)
	productService := service.NewProductService(&productRepository, storageService, &aprioriRepository, db)
	transactionService := service.NewTransactionService(&transactionRepository, &productRepository, db)
	aprioriService := service.NewAprioriService(&transactionRepository, storageService, &productRepository, &aprioriRepository, db)
	paymentService := service.NewPaymentService(configuration, &paymentRepository, &userOrderRepository, &transactionRepository, db)
	userOrderService := service.NewUserOrderService(&paymentRepository, &userOrderRepository, &userRepository, db)
	cacheService := service.NewCacheService(configuration)
	categoryService := service.NewCategoryService(&categoryRepository, db)
	commentService := service.NewCommentService(&commentRepository, &productRepository, db)

	// Setup Controller
	userController := controller.NewUserController(&userService)
	authController := controller.NewAuthController(&userService, jwtService, emailService, &passwordResetService)
	productController := controller.NewProductController(&productService, &storageService, &cacheService)
	transactionController := controller.NewTransactionController(&transactionService, &storageService, &cacheService)
	aprioriController := controller.NewAprioriController(aprioriService, &storageService, &cacheService)
	paymentController := controller.NewPaymentController(&paymentService, &userOrderService, emailService, &cacheService)
	userOrderController := controller.NewUserOrderController(&paymentService, &userOrderService, &cacheService)
	categoryController := controller.NewCategoryController(&categoryService, &cacheService)
	commentController := controller.NewCommentController(&commentService)
	rajaOngkirController := controller.NewRajaOngkirController()

	// CORS Middleware
	router.Use(middleware.SetupCorsMiddleware())

	// Main Route
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
	categoryController.Route(router)
	commentController.Route(router)
	rajaOngkirController.Route(router)

	return router, db
}
