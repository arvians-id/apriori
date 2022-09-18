package integration

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/arvians-id/apriori/config"
	"github.com/arvians-id/apriori/entity"
	repository "github.com/arvians-id/apriori/repository/postgres"
	"github.com/arvians-id/apriori/service"
	"github.com/arvians-id/apriori/tests/setup"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
)

var _ = Describe("Product API", func() {
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

	Describe("Create Product /products", func() {
		When("the fields are correct", func() {
			When("the fields are filled", func() {
				It("should return successful create product response", func() {
					// Create Product
					requestBody := strings.NewReader(`{"code": "SK6","name": "Bantal Biasa","description": "Test","category": "Bantal, Kasur","mass": 1000,"price": 7000}`)
					request := httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
					Expect(responseBody["status"]).To(Equal("created"))
					Expect(responseBody["data"].(map[string]interface{})["code"]).To(Equal("SK6"))
					Expect(responseBody["data"].(map[string]interface{})["name"]).To(Equal("Bantal Biasa"))
					Expect(responseBody["data"].(map[string]interface{})["description"]).To(Equal("Test"))
					Expect(int(responseBody["data"].(map[string]interface{})["price"].(float64))).To(Equal(7000))
				})
			})
		})
	})

	Describe("Update Product /products/:code", func() {
		When("the product is not found", func() {
			It("should return error not found", func() {
				// Update Product
				requestBody := strings.NewReader(`{"code": "SK1","name": "Bantal Biasa","description": "Test","category": "Bantal, Kasur","mass":1000}`)
				request := httptest.NewRequest(http.MethodPatch, "/api/products/SK1", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("the fields are correct", func() {
			When("the fields are filled", func() {
				It("should return successful update product response", func() {
					// Create Product
					tx, _ := database.Begin()
					productRepository := repository.NewProductRepository()
					description := "Test"
					row, _ := productRepository.Create(context.Background(), tx, &entity.Product{
						Code:        "SK6",
						Name:        "Widdy",
						Description: &description,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					})
					_ = tx.Commit()

					// Update Product
					requestBody := strings.NewReader(`{"code": "SK1","name": "Guling Doti","description": "Test Bang","category": "Bantal, Kasur","mass":1000}`)
					request := httptest.NewRequest(http.MethodPatch, "/api/products/"+row.Code, requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
					Expect(responseBody["status"]).To(Equal("updated"))
					Expect(responseBody["data"].(map[string]interface{})["name"]).ShouldNot(Equal("Widdy"))
					Expect(responseBody["data"].(map[string]interface{})["description"]).ShouldNot(Equal("Test"))
				})
			})
		})
	})

	Describe("Delete Product /products/:code", func() {
		When("product is not found", func() {
			It("should return error not found", func() {
				// Delete Product
				request := httptest.NewRequest(http.MethodDelete, "/api/products/SK9", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("product is found", func() {
			It("should return a successful delete product response", func() {
				// Create Product
				tx, _ := database.Begin()
				productRepository := repository.NewProductRepository()
				description := "Test"
				row, _ := productRepository.Create(context.Background(), tx, &entity.Product{
					Code:        "SK6",
					Name:        "Widdy",
					Description: &description,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
				_ = tx.Commit()

				// Delete Product
				request := httptest.NewRequest(http.MethodDelete, "/api/products/"+row.Code, nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("deleted"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})
	})

	Describe("Find All Product /products", func() {
		When("the product is not present", func() {
			It("should return a successful but the data is null", func() {
				// Find All Product
				request := httptest.NewRequest(http.MethodGet, "/api/products", nil)
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
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("the product is present", func() {
			It("should return a successful and show all products", func() {
				// Create Product
				tx, _ := database.Begin()
				productRepository := repository.NewProductRepository()
				description := "Test"
				product1, _ := productRepository.Create(context.Background(), tx, &entity.Product{
					Code:        "SK6",
					Name:        "Guling",
					Description: &description,
					Category:    "Bantal, Kasur",
					Mass:        1000,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
				description = "Test Bang"
				product2, _ := productRepository.Create(context.Background(), tx, &entity.Product{
					Code:        "SK1",
					Name:        "Bantal",
					Description: &description,
					Category:    "Bantal, Kasur",
					Mass:        1000,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
				_ = tx.Commit()

				// Find All Products
				request := httptest.NewRequest(http.MethodGet, "/api/products", nil)
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

				products := responseBody["data"].([]interface{})
				Expect(product1.Code).To(Equal(products[1].(map[string]interface{})["code"]))
				Expect(product1.Name).To(Equal(products[1].(map[string]interface{})["name"]))

				Expect(product2.Code).To(Equal(products[0].(map[string]interface{})["code"]))
				Expect(product2.Name).To(Equal(products[0].(map[string]interface{})["name"]))
			})
		})
	})

	Describe("Find All Product On Admin /products-admin", func() {
		When("the product is not present", func() {
			It("should return a successful but the data is null", func() {
				// Find All Product
				request := httptest.NewRequest(http.MethodGet, "/api/products-admin", nil)
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
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("the product is present", func() {
			It("should return a successful and show all products with different status empty", func() {
				// Create Product
				tx, _ := database.Begin()
				productRepository := repository.NewProductRepository()
				description := "Test"
				product1, _ := productRepository.Create(context.Background(), tx, &entity.Product{
					Code:        "SK6",
					Name:        "Guling",
					Description: &description,
					Category:    "Bantal, Kasur",
					Mass:        1000,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
				description = "Test Bang"
				_, _ = productRepository.Create(context.Background(), tx, &entity.Product{
					Code:        "SK1",
					Name:        "Bantal",
					Description: &description,
					Category:    "Bantal, Kasur",
					Mass:        1000,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
				_ = tx.Commit()

				// Update Product
				requestBody := strings.NewReader(`{"name": "Guling Doti Bang","description": "Test Bang","category": "Bantal, Kasur","mass":1000}`)
				request := httptest.NewRequest(http.MethodPatch, "/api/products/"+product1.Code, requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				// Find All Products
				request = httptest.NewRequest(http.MethodGet, "/api/products-admin", nil)
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

				products := responseBody["data"].([]interface{})
				Expect(products[1].(map[string]interface{})["code"]).To(Equal("SK6"))
				Expect(products[1].(map[string]interface{})["name"]).To(Equal("Guling Doti Bang"))

				Expect(products[0].(map[string]interface{})["code"]).To(Equal("SK1"))
				Expect(products[0].(map[string]interface{})["name"]).To(Equal("Bantal"))
			})
		})
	})

	Describe("Find All Product By Similar Category /products/:code/category", func() {
		When("the product similar is not present", func() {
			It("should return a successful but the data is null", func() {
				// Create Product
				tx, _ := database.Begin()
				productRepository := repository.NewProductRepository()
				description := "Test"
				product1, _ := productRepository.Create(context.Background(), tx, &entity.Product{
					Code:        "SK6",
					Name:        "Guling",
					Description: &description,
					Category:    "Bantal, Kasur",
					Mass:        1000,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
				description = "Test Bang"
				_, _ = productRepository.Create(context.Background(), tx, &entity.Product{
					Code:        "SK1",
					Name:        "Bantal",
					Description: &description,
					Category:    "Elektronik, Guling",
					Mass:        1000,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
				_ = tx.Commit()

				// Find All Products
				request := httptest.NewRequest(http.MethodGet, "/api/products/"+product1.Code+"/category", nil)
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
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("the product is present", func() {
			It("should return a successful and show all products by similar category", func() {
				// Create Product
				tx, _ := database.Begin()
				productRepository := repository.NewProductRepository()
				description := "Test"
				product1, _ := productRepository.Create(context.Background(), tx, &entity.Product{
					Code:        "SK6",
					Name:        "Guling",
					Description: &description,
					Category:    "Bantal, Kasur",
					Mass:        1000,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
				description = "Test Bang"
				_, _ = productRepository.Create(context.Background(), tx, &entity.Product{
					Code:        "SK1",
					Name:        "Bantal",
					Description: &description,
					Category:    "Bantal, Kasur",
					Mass:        1000,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
				_ = tx.Commit()

				// Find All Products
				request := httptest.NewRequest(http.MethodGet, "/api/products/"+product1.Code+"/category", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				products := responseBody["data"].([]interface{})
				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(products[0].(map[string]interface{})["code"]).To(Equal("SK1"))
				Expect(products[0].(map[string]interface{})["name"]).To(Equal("Bantal"))
				Expect(products[0].(map[string]interface{})["description"]).To(Equal("Test Bang"))
				Expect(products[0].(map[string]interface{})["category"]).To(Equal("Bantal, Kasur"))
			})
		})
	})

	Describe("Find By Code Product /products/:code", func() {
		When("product is not found", func() {
			It("should return error not found", func() {
				// Find By Code Product
				request := httptest.NewRequest(http.MethodGet, "/api/products/SK5", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("product is found", func() {
			It("should return a successful find product by code", func() {
				// Create Product
				tx, _ := database.Begin()
				productRepository := repository.NewProductRepository()
				description := "Test"
				row, _ := productRepository.Create(context.Background(), tx, &entity.Product{
					Code:        "SK6",
					Name:        "Widdy",
					Description: &description,
					Category:    "Bantal, Kasur",
					Mass:        1000,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
				_ = tx.Commit()

				// Find By Code Product
				request := httptest.NewRequest(http.MethodGet, "/api/products/"+row.Code, nil)
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
				Expect(responseBody["data"].(map[string]interface{})["code"]).To(Equal("SK6"))
				Expect(responseBody["data"].(map[string]interface{})["name"]).To(Equal("Widdy"))
				Expect(responseBody["data"].(map[string]interface{})["description"]).To(Equal("Test"))
			})
		})
	})

	Describe("Find All Recommendation Product By Code /products/:code/recommendation", func() {
		When("product recommendation is not found", func() {
			It("should return error not found", func() {
				// Create Product
				tx, _ := database.Begin()
				productRepository := repository.NewProductRepository()
				description := "Test Bang"
				product, _ := productRepository.Create(context.Background(), tx, &entity.Product{
					Code:        "SK1",
					Name:        "Bantal Biasa",
					Description: &description,
					Category:    "Bantal, Kasur",
					Mass:        1000,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})

				// Create Apriori
				aprioriRepository := repository.NewAprioriRepository()
				var aprioriRequests []*entity.Apriori
				image := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/assets/%s", os.Getenv("AWS_BUCKET"), os.Getenv("AWS_REGION"), "no-image.png")
				aprioriRequests = append(aprioriRequests, &entity.Apriori{
					Code:       "uRwCmCplpF",
					Item:       "guling biasa",
					Discount:   25.00,
					Support:    50.00,
					Confidence: 71.43,
					RangeDate:  "2021-05-21 - 2022-05-21",
					IsActive:   true,
					Image:      &image,
					CreatedAt:  time.Now(),
				})
				_ = aprioriRepository.Create(context.Background(), tx, aprioriRequests)
				_ = tx.Commit()

				// Find All Recommendation
				request := httptest.NewRequest(http.MethodGet, "/api/products/"+product.Code+"/recommendation", nil)
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
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("product recommendation is found", func() {
			It("should return a successful find recommendation product by code", func() {
				// Create Product
				tx, _ := database.Begin()
				productRepository := repository.NewProductRepository()
				description := "Test Bang"
				product, _ := productRepository.Create(context.Background(), tx, &entity.Product{
					Code:        "SK1",
					Name:        "Bantal Biasa",
					Description: &description,
					Category:    "Bantal, Kasur",
					Mass:        1000,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})

				// Create Apriori
				aprioriRepository := repository.NewAprioriRepository()
				var aprioriRequests []*entity.Apriori
				image := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/assets/%s", os.Getenv("AWS_BUCKET"), os.Getenv("AWS_REGION"), "no-image.png")
				aprioriRequests = append(aprioriRequests, &entity.Apriori{
					Code:       "uRwCmCplpF",
					Item:       "bantal biasa",
					Discount:   25.00,
					Support:    50.00,
					Confidence: 71.43,
					RangeDate:  "2021-05-21 - 2022-05-21",
					IsActive:   true,
					Image:      &image,
					CreatedAt:  time.Now(),
				})
				_ = aprioriRepository.Create(context.Background(), tx, aprioriRequests)
				_ = tx.Commit()

				// Find All Recommendation
				request := httptest.NewRequest(http.MethodGet, "/api/products/"+product.Code+"/recommendation", nil)
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

				products := responseBody["data"].([]interface{})
				Expect(products[0].(map[string]interface{})["apriori_code"]).To(Equal("uRwCmCplpF"))
				Expect(products[0].(map[string]interface{})["apriori_item"]).To(Equal("bantal biasa"))
			})
		})
	})
})
