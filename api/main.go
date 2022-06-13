package main

import (
	"apriori/api/controller"
	"apriori/api/middleware"
	"apriori/config"
	"apriori/repository"
	"apriori/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	config.SetupConfiguration()
}

func main() {
	// Setup Configuration
	err := godotenv.Load()
	if err != nil {
		log.Fatal("cannot load env ", err)
		return
	}

	router := gin.Default()
	db := config.NewDB()

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

	// Setup Proxies
	err = router.SetTrustedProxies([]string{os.Getenv("APP_URL")})
	if err != nil {
		log.Fatal("cannot set proxies ", err)
		return
	}

	// CORS Middleware
	router.Use(middleware.SetupCorsMiddleware())

	// Setup Router
	authController.Route(router)

	// Authentication Middleware
	//router.Use(middleware.AuthJwtMiddleware())

	userController.Route(router)
	productController.Route(router)
	transactionController.Route(router)
	aprioriController.Route(router)

	// Start App
	addr := fmt.Sprintf(":%v", os.Getenv("APP_PORT"))
	err = router.Run(addr)
	if err != nil {
		log.Fatal("cannot run server ", err)
		return
	}
}
