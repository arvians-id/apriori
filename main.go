package main

import (
	"apriori/app"
	"apriori/controller"
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

	userController.Route(router)

	router.GET("/", func(c *gin.Context) {
		c.String(200, "asus")
	})

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
