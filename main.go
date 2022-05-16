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
	err := godotenv.Load()
	if err != nil {
		helper.PanicIfError(err)
	}

	router := gin.New()
	db := app.NewDB()
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db)
	userController := controller.NewUserController(userService)

	authRepository := repository.NewAuthRepository()
	authService := service.NewAuthService(userRepository, authRepository, db)
	jwtService := service.NewJwtService()
	authController := controller.NewAuthController(authService, userService, jwtService)

	err = router.SetTrustedProxies([]string{os.Getenv("APP_URL")})
	if err != nil {
		helper.PanicIfError(err)
	}
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	authController.Route(router)
	userController.Route(router)

	router.GET("/", func(c *gin.Context) {
		c.String(200, "asus")
	})

	addr := fmt.Sprintf(":%v", os.Getenv("APP_PORT"))
	err = router.Run(addr)
	if err != nil {
		panic(err)
	}
}
