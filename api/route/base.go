package route

import (
	"apriori/api/controller"
	"apriori/api/middleware"
	"apriori/config"
	"apriori/repository"
	"apriori/service"
	"database/sql"
	"github.com/gin-gonic/gin"
)

func NewInitializedServer(configuration config.Config) (*gin.Engine, *sql.DB) {
	// Setup Configuration
	router := gin.Default()
	db, err := config.NewPostgreSQL(configuration)
	if err != nil {
		panic(err)
	}

	// Setup Repository
	userRepository := repository.NewUserRepository()
	passwordRepository := repository.NewPasswordResetRepository()
	authRepository := repository.NewAuthRepository()
	productRepository := repository.NewProductRepository()
	transactionRepository := repository.NewTransactionRepository()
	aprioriRepository := repository.NewAprioriRepository()

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

	// Setup Controller
	userController := controller.NewUserController(&userService)
	authController := controller.NewAuthController(&authService, &userService, jwtService, emailService, &passwordResetService)
	productController := controller.NewProductController(&productService, &storageService)
	transactionController := controller.NewTransactionController(&transactionService, &storageService)
	aprioriController := controller.NewAprioriController(aprioriService, &storageService)

	// CORS Middleware
	router.Use(middleware.SetupCorsMiddleware())

	// Setup Router
	authController.Route(router)

	userController.Route(router)
	productController.Route(router)
	transactionController.Route(router)
	aprioriController.Route(router)

	return router, db
}
