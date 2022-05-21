package main

import (
	"apriori/app"
	"apriori/controller"
	"apriori/helper"
	"apriori/repository"
	"apriori/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	// Setup Configuration
	err := godotenv.Load()
	if err != nil {
		helper.PanicIfError(err)
	}
	router := gin.Default()
	db := app.NewDB()

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

	// Setup Controller
	userController := controller.NewUserController(&userService)
	authController := controller.NewAuthController(&authService, &userService, jwtService, emailService, &passwordResetService)
	productController := controller.NewProductController(&productService)
	transactionController := controller.NewTransactionController(&transactionService)

	// Setup Proxies
	err = router.SetTrustedProxies([]string{os.Getenv("APP_URL")})
	if err != nil {
		helper.PanicIfError(err)
	}

	// Setup Router
	authController.Route(router)
	userController.Route(router)
	productController.Route(router)
	transactionController.Route(router)

	// Start App
	addr := fmt.Sprintf(":%v", os.Getenv("APP_PORT"))
	err = router.Run(addr)
	if err != nil {
		panic(err)
	}
}
