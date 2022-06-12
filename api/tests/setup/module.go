package setup

import (
	"apriori/api/controller"
	"apriori/api/middleware"
	"apriori/repository"
	"apriori/service"
	"database/sql"
	"github.com/gin-gonic/gin"
)

func ModuleSetup(db *sql.DB) *gin.Engine {
	router := gin.New()

	// Setup Repository
	userRepository := repository.NewUserRepository()
	passwordRepository := repository.NewPasswordResetRepository()
	authRepository := repository.NewAuthRepository()
	productRepository := repository.NewProductRepository()
	transactionRepository := repository.NewTransactionRepository()
	aprioriRepository := repository.NewAprioriRepository()

	// Setup Service
	userService := service.NewUserService(&userRepository, db)
	authService := service.NewAuthService(&userRepository, &authRepository, db)
	jwtService := service.NewJwtService()
	emailService := service.NewEmailService()
	passwordResetService := service.NewPasswordResetService(&passwordRepository, &userRepository, db)
	productService := service.NewProductService(&productRepository, db)
	transactionService := service.NewTransactionService(&transactionRepository, &productRepository, db)
	aprioriService := service.NewAprioriService(&transactionRepository, &aprioriRepository, db)

	// Setup Controller
	userController := controller.NewUserController(&userService)
	authController := controller.NewAuthController(&authService, &userService, jwtService, emailService, &passwordResetService)
	productController := controller.NewProductController(&productService)
	transactionController := controller.NewTransactionController(&transactionService)
	aprioriController := controller.NewAprioriController(aprioriService)

	// Setup Router
	authController.Route(router)

	// Auth Middleware
	router.Use(middleware.AuthJwtMiddleware())
	userController.Route(router)
	productController.Route(router)
	transactionController.Route(router)
	aprioriController.Route(router)

	return router
}
