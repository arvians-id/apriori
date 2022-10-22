package server

import (
	"database/sql"
	"github.com/arvians-id/apriori/cmd/config"
	"github.com/arvians-id/apriori/internal/http/controller/rest"
	"github.com/arvians-id/apriori/internal/http/middleware"
	"github.com/arvians-id/apriori/internal/repository/postgres"
	"github.com/arvians-id/apriori/internal/service"
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

	// Main Middleware
	middleware.RegisterPrometheusMetrics()
	router.Use(middleware.SetupCorsMiddleware())
	router.Use(middleware.GinContextToContextMiddleware())
	router.Use(middleware.PrometheusMetricsMiddleware())

	// Main Route
	NewInitializedMainRoute(router,
		aprioriService,
		cacheService,
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
