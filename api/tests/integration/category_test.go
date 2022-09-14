package integration

import (
	"apriori/config"
	"apriori/entity"
	"apriori/helper"
	repository "apriori/repository/postgres"
	"apriori/service"
	"apriori/tests/setup"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

var _ = Describe("Category API", func() {
	var server *gin.Engine
	var database *sql.DB
	var tokenJWT string
	var cookie *http.Cookie
	configuration := config.New("../../.env.test")

	BeforeEach(func() {
		// Setup Configuration
		router, db := setup.ModuleSetup(configuration)

		database = db
		server = router

		// User Authentication
		// Create user
		tx, _ := database.Begin()
		userRepository := repository.NewUserRepository()
		password, _ := bcrypt.GenerateFromPassword([]byte("Rahasia123"), bcrypt.DefaultCost)
		_, _ = userRepository.Create(context.Background(), tx, &entity.User{
			Role:      1,
			Name:      "Widdy",
			Email:     "widdy@gmail.com",
			Password:  string(password),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		_ = tx.Commit()

		// Login
		requestBody := strings.NewReader(`{"email": "widdy@gmail.com","password":"Rahasia123"}`)
		request := httptest.NewRequest(http.MethodPost, "/api/auth/login", requestBody)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

		writer := httptest.NewRecorder()
		server.ServeHTTP(writer, request)

		var responseBody map[string]interface{}
		_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

		tokenJWT = responseBody["data"].(map[string]interface{})["access_token"].(string)
		for _, c := range writer.Result().Cookies() {
			if c.Name == "token" {
				cookie = c
			}
		}
	})

	AfterEach(func() {
		// Setup Configuration
		_, db := setup.ModuleSetup(configuration)
		defer db.Close()

		productCache := service.NewCacheService(configuration)
		_ = productCache.FlushDB(context.Background())

		err := setup.TearDownTest(db)
		if err != nil {
			log.Fatal(err)
		}
	})

	Describe("Create Category /categories", func() {
		When("the fields are filled", func() {
			It("should return successful create category response", func() {
				// Create Category
				requestBody := strings.NewReader(`{"name": "Produk Kasur"}`)
				request := httptest.NewRequest(http.MethodPost, "/api/categories", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))

				Expect(responseBody["data"].(map[string]interface{})["name"]).To(Equal("Produk Kasur"))
			})
		})
	})

	Describe("Find All Category /categories", func() {
		When("the categories is not found", func() {
			It("should return successful find all category response with empty data", func() {
				// Find All Category
				request := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("the categories is exist", func() {
			It("should return successful find all category response", func() {
				// Create Category
				tx, _ := database.Begin()
				categoryRepository := repository.NewCategoryRepository()
				_, _ = categoryRepository.Create(context.Background(), tx, &entity.Category{
					Name:      "Produk Kasur",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
				_ = tx.Commit()

				// Find All Category
				request := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))

				categoryResponse := responseBody["data"].([]interface{})
				Expect(categoryResponse[0].(map[string]interface{})["name"]).To(Equal("Produk Kasur"))
			})
		})
	})

	Describe("Find Category By Id /categories/:id", func() {
		When("the categories is not found", func() {
			It("should return error not found", func() {
				// Find Category By Id
				request := httptest.NewRequest(http.MethodGet, "/api/categories/1", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["status"]).To(Equal("data not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("the categories is exist", func() {
			It("should return successful find category by id response", func() {
				// Create Category
				tx, _ := database.Begin()
				categoryRepository := repository.NewCategoryRepository()
				category, _ := categoryRepository.Create(context.Background(), tx, &entity.Category{
					Name:      "Produk Kasur",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
				_ = tx.Commit()

				// Find Category By Id
				request := httptest.NewRequest(http.MethodGet, "/api/categories/"+helper.IntToStr(category.IdCategory), nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"].(map[string]interface{})["name"]).To(Equal("Produk Kasur"))
			})
		})
	})

	Describe("Update Category By Id /categories/:id", func() {
		When("the categories is not found", func() {
			It("should return error not found", func() {
				// Update Category By Id
				requestBody := strings.NewReader(`{"name": "Produk Bantal"}`)
				request := httptest.NewRequest(http.MethodPatch, "/api/categories/1", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["status"]).To(Equal("data not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("the categories is exist", func() {
			It("should return successful update category response", func() {
				// Create Category
				tx, _ := database.Begin()
				categoryRepository := repository.NewCategoryRepository()
				category, _ := categoryRepository.Create(context.Background(), tx, &entity.Category{
					Name:      "Produk Kasur",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
				_ = tx.Commit()

				// Update Category By Id
				requestBody := strings.NewReader(`{"name": "Produk Bantal"}`)
				request := httptest.NewRequest(http.MethodPatch, "/api/categories/"+helper.IntToStr(category.IdCategory), requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				// Find Category By Id
				request = httptest.NewRequest(http.MethodGet, "/api/categories/"+helper.IntToStr(category.IdCategory), nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"].(map[string]interface{})["name"]).To(Equal("Produk Bantal"))
			})
		})
	})

	Describe("Delete Category By Id /categories/:id", func() {
		When("the categories is not found", func() {
			It("should return error not found", func() {
				// Delete Category By Id
				request := httptest.NewRequest(http.MethodDelete, "/api/categories/1", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["status"]).To(Equal("data not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("the categories is exist", func() {
			It("should return successful delete category by id and set not found after deleted", func() {
				// Create Category
				tx, _ := database.Begin()
				categoryRepository := repository.NewCategoryRepository()
				category, _ := categoryRepository.Create(context.Background(), tx, &entity.Category{
					Name:      "Produk Kasur",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
				_ = tx.Commit()

				// Delete Category By Id
				request := httptest.NewRequest(http.MethodDelete, "/api/categories/"+helper.IntToStr(category.IdCategory), nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				// Find Category By Id
				request = httptest.NewRequest(http.MethodGet, "/api/categories/"+helper.IntToStr(category.IdCategory), nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["status"]).To(Equal("data not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})
	})
})
