package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	router.GET("/", func(context *gin.Context) {
		context.String(200, "Asuas")
	})

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
