package main

import (
	"apriori/app"
	"apriori/controller"
	"apriori/helper"
	"apriori/repository"
	"apriori/service"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	db := app.NewDB()
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db)
	userController := controller.NewUserController(userService)

	authRepository := repository.NewAuthRepository()
	authService := service.NewAuthService(userRepository, authRepository, db)
	jwtService := service.NewJwtService()
	authController := controller.NewAuthController(authService, userService, jwtService)

	err := router.SetTrustedProxies([]string{"127.0.0.1"})
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

	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
