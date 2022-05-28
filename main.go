package main

import (
	"apriori/api/controller"
	"apriori/config"
	"apriori/repository"
	"apriori/service"
	"apriori/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	// Setup Configuration
	err := godotenv.Load()
	utils.PanicIfError(err)

	router := gin.Default()
	db := config.NewDB()

	// Setup Repository
	userRepository := repository.NewUserRepository()
	passwordRepository := repository.NewPasswordResetRepository()
	authRepository := repository.NewAuthRepository()
	productRepository := repository.NewProductRepository()
	transactionRepository := repository.NewTransactionRepository()

	// Setup Service
	userService := service.NewUserService(&userRepository, db)
	authService := service.NewAuthService(&userRepository, &authRepository, db)
	jwtService := service.NewJwtService()
	emailService := service.NewEmailService()
	passwordResetService := service.NewPasswordResetService(&passwordRepository, &userRepository, db)
	productService := service.NewProductService(&productRepository, db)
	transactionService := service.NewTransactionService(&transactionRepository, &productRepository, db)
	aprioriService := service.NewAprioriService(&transactionRepository, db)

	// Setup Controller
	userController := controller.NewUserController(&userService)
	authController := controller.NewAuthController(&authService, &userService, jwtService, emailService, &passwordResetService)
	productController := controller.NewProductController(&productService)
	transactionController := controller.NewTransactionController(&transactionService)
	aprioriController := controller.NewAprioriController(aprioriService)

	// Setup Proxies
	err = router.SetTrustedProxies([]string{os.Getenv("APP_URL")})
	utils.PanicIfError(err)

	// Setup Router
	authController.Route(router)
	userController.Route(router)
	productController.Route(router)
	transactionController.Route(router)
	aprioriController.Route(router)

	// Start App
	addr := fmt.Sprintf(":%v", os.Getenv("APP_PORT"))
	err = router.Run(addr)
	utils.PanicIfError(err)
}
