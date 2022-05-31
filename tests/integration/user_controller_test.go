package integration_test

import (
	"apriori/api/controller"
	"apriori/repository"
	"apriori/service"
	"database/sql"
	gin "github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	"time"
)

var _ = Describe("User Controller", func() {

	BeforeEach(func() {
		db, err := sql.Open("mysql", "root@tcp(localhost:3306)/apriori_test")
		if err != nil {
			panic(err)
		}

		err = db.Ping()
		if err != nil {
			panic(err)
		}

		router := gin.New()
		db.SetMaxIdleConns(5)
		db.SetMaxOpenConns(20)
		db.SetConnMaxLifetime(60 * time.Minute)
		db.SetConnMaxIdleTime(10 * time.Minute)

		userRepository := repository.NewUserRepository()
		userService := service.NewUserService(&userRepository, db)
		userController := controller.NewUserController(&userService)

		userController.Route(router)

	})

	AfterEach(func() {
		db, err := sql.Open("mysql", "root@tcp(localhost:3306)/apriori_test")
		if err != nil {
			panic(err)
		}

		_, err = db.Exec(`TRUNCATE TABLE users;`)

		if err != nil {
			panic(err)
		}
	})
})
