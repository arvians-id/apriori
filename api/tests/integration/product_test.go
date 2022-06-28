package integration

import (
	"apriori/config"
	"apriori/entity"
	"apriori/repository"
	"apriori/tests/setup"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
)

var _ = Describe("User API", func() {

	var server *gin.Engine
	var database *sql.DB
	var tokenJWT string
	var cookie *http.Cookie

	BeforeEach(func() {
		// Setup Configuration
		configuration := config.New("../../.env.test")

		db, err := setup.SuiteSetupMySQL(configuration)
		if err != nil {
			panic(err)
		}
		router := setup.ModuleSetup(configuration)

		database = db
		server = router

		// User Authentication
		// Create user
		tx, _ := database.Begin()
		userRepository := repository.NewUserRepository()
		password, _ := bcrypt.GenerateFromPassword([]byte("Rahasia123"), bcrypt.DefaultCost)
		_, _ = userRepository.Create(context.Background(), tx, entity.User{
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

		writer := httptest.NewRecorder()
		server.ServeHTTP(writer, request)

		response := writer.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		_ = json.Unmarshal(body, &responseBody)

		tokenJWT = responseBody["data"].(map[string]interface{})["access_token"].(string)
		for _, c := range writer.Result().Cookies() {
			if c.Name == "token" {
				cookie = c
			}
		}
	})

	AfterEach(func() {
		// Setup Configuration
		configuration := config.New("../../.env.test")
		db, err := setup.SuiteSetupMySQL(configuration)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		err = setup.TearDownTest(db)
		if err != nil {
			panic(err)
		}
	})

	Describe("Create Product /products", func() {
		When("the fields are incorrect", func() {
			When("the code field is incorrect", func() {
				It("should return error required", func() {
					// Create Product
					requestBody := strings.NewReader(`{"name": "Bantal Biasa","description": "Test"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateProductRequest.Code' Error:Field validation for 'Code' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error duplicate values", func() {
					// Create Product
					requestBody := strings.NewReader(`{"code": "SK6","name": "Bantal Biasa","description": "Test"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					requestBody = strings.NewReader(`{"code": "SK6","name": "Bantal Biasa","description": "Test"}`)
					request = httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer = httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusInternalServerError))
					Expect(responseBody["status"]).To(Equal("Error 1062: Duplicate entry 'SK6' for key 'code'"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error less character of length", func() {
					// Create Product
					requestBody := strings.NewReader(`{"code": "S","name": "Bantal Biasa","description": "Test"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateProductRequest.Code' Error:Field validation for 'Code' failed on the 'min' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
			When("the name field is incorrect", func() {
				It("should return error required", func() {
					// Create Product
					requestBody := strings.NewReader(`{"code": "SK6","description": "Test"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateProductRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error less character of length", func() {
					// Create Product
					requestBody := strings.NewReader(`{"code": "SK6","name": "saa","description": "Test"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateProductRequest.Name' Error:Field validation for 'Name' failed on the 'min' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
		})

		When("the fields are correct", func() {
			When("the description field is empty", func() {
				It("should return successful create product response", func() {
					// Create Product
					requestBody := strings.NewReader(`{"code": "SK6","name": "Bantal Biasa"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
					Expect(responseBody["status"]).To(Equal("created"))
					Expect(responseBody["data"].(map[string]interface{})["code"]).To(Equal("SK6"))
					Expect(responseBody["data"].(map[string]interface{})["name"]).To(Equal("Bantal Biasa"))
					Expect(responseBody["data"].(map[string]interface{})["description"]).To(Equal(""))
				})
			})

			When("the fields are filled", func() {
				It("should return successful create product response", func() {
					// Create Product
					requestBody := strings.NewReader(`{"code": "SK6","name": "Bantal Biasa","description": "Test"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
					Expect(responseBody["status"]).To(Equal("created"))
					Expect(responseBody["data"].(map[string]interface{})["code"]).To(Equal("SK6"))
					Expect(responseBody["data"].(map[string]interface{})["name"]).To(Equal("Bantal Biasa"))
					Expect(responseBody["data"].(map[string]interface{})["description"]).To(Equal("Test"))
				})
			})
		})
	})

	Describe("Update Product /products/:code", func() {
		When("the product is not found", func() {
			It("should return error not found", func() {
				// Update Product
				requestBody := strings.NewReader(`{"code": "SK1","name": "Bantal Biasa","description": "Test"}`)
				request := httptest.NewRequest(http.MethodPatch, "/api/products/SK1", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusInternalServerError))
				Expect(responseBody["status"]).To(Equal("product not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})
		When("the fields are incorrect", func() {
			When("the name field is incorrect", func() {
				It("should return error required", func() {
					// Create Product
					tx, _ := database.Begin()
					productRepository := repository.NewProductRepository()
					row, _ := productRepository.Create(context.Background(), tx, entity.Product{
						Code:        "SK6",
						Name:        "Widdy",
						Description: "Test",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					})
					_ = tx.Commit()

					// Update Product
					requestBody := strings.NewReader(`{"code": "SK1","description": "Test"}`)
					request := httptest.NewRequest(http.MethodPatch, "/api/products/"+row.Code, requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'UpdateProductRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error less character of length", func() {
					// Create Product
					tx, _ := database.Begin()
					productRepository := repository.NewProductRepository()
					row, _ := productRepository.Create(context.Background(), tx, entity.Product{
						Code:        "SK6",
						Name:        "Widdy",
						Description: "Test",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					})
					_ = tx.Commit()

					// Update Product
					requestBody := strings.NewReader(`{"code": "SK1","name": "Kas","description": "Test"}`)
					request := httptest.NewRequest(http.MethodPatch, "/api/products/"+row.Code, requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'UpdateProductRequest.Name' Error:Field validation for 'Name' failed on the 'min' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
		})

		When("the fields are correct", func() {
			When("the description field is empty", func() {
				It("should return successful create product response", func() {
					// Create Product
					tx, _ := database.Begin()
					productRepository := repository.NewProductRepository()
					row, _ := productRepository.Create(context.Background(), tx, entity.Product{
						Code:        "SK6",
						Name:        "Widdy",
						Description: "Test",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					})
					_ = tx.Commit()

					// Update Product
					requestBody := strings.NewReader(`{"code": "SK1","name": "Guling Doti"}`)
					request := httptest.NewRequest(http.MethodPatch, "/api/products/"+row.Code, requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
					Expect(responseBody["status"]).To(Equal("updated"))
					Expect(responseBody["data"].(map[string]interface{})["name"]).ShouldNot(Equal("Widdy"))
					Expect(responseBody["data"].(map[string]interface{})["description"]).To(Equal(""))
				})
			})

			When("the fields are filled", func() {
				It("should return successful create product response", func() {
					// Create Product
					tx, _ := database.Begin()
					productRepository := repository.NewProductRepository()
					row, _ := productRepository.Create(context.Background(), tx, entity.Product{
						Code:        "SK6",
						Name:        "Widdy",
						Description: "Test",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					})
					_ = tx.Commit()

					// Update Product
					requestBody := strings.NewReader(`{"code": "SK1","name": "Guling Doti","description": "Test Bang"}`)
					request := httptest.NewRequest(http.MethodPatch, "/api/products/"+row.Code, requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

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
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusInternalServerError))
				Expect(responseBody["status"]).To(Equal("product not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("product is found", func() {
			It("should return a successful delete product response", func() {
				// Create Product
				tx, _ := database.Begin()
				productRepository := repository.NewProductRepository()
				row, _ := productRepository.Create(context.Background(), tx, entity.Product{
					Code:        "SK6",
					Name:        "Widdy",
					Description: "Test",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
				_ = tx.Commit()

				// Delete Product
				request := httptest.NewRequest(http.MethodDelete, "/api/products/"+row.Code, nil)
				request.Header.Add("Content-Type", "application/json")
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

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
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

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
				product1, _ := productRepository.Create(context.Background(), tx, entity.Product{
					Code:        "SK6",
					Name:        "Guling",
					Description: "Test",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
				product2, _ := productRepository.Create(context.Background(), tx, entity.Product{
					Code:        "SK1",
					Name:        "Bantal",
					Description: "Test Bang",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
				_ = tx.Commit()

				// Find All Products
				request := httptest.NewRequest(http.MethodGet, "/api/products", nil)
				request.Header.Add("Content-Type", "application/json")
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				products := responseBody["data"].([]interface{})

				// Desc
				productResponse1 := products[1].(map[string]interface{})
				productResponse2 := products[0].(map[string]interface{})

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))

				Expect(product1.Code).To(Equal(productResponse1["code"]))
				Expect(product1.Name).To(Equal(productResponse1["name"]))

				Expect(product2.Code).To(Equal(productResponse2["code"]))
				Expect(product2.Name).To(Equal(productResponse2["name"]))
			})
		})
	})

	Describe("Find By Code Product /products/:code", func() {
		When("product is not found", func() {
			It("should return error not found", func() {
				// Find By Code Product
				request := httptest.NewRequest(http.MethodGet, "/api/products/SK5", nil)
				request.Header.Add("Content-Type", "application/json")
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusInternalServerError))
				Expect(responseBody["status"]).To(Equal("product not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})
		When("product is found", func() {
			It("should return a successful find product by code", func() {
				// Create Product
				tx, _ := database.Begin()
				productRepository := repository.NewProductRepository()
				row, _ := productRepository.Create(context.Background(), tx, entity.Product{
					Code:        "SK6",
					Name:        "Widdy",
					Description: "Test",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
				_ = tx.Commit()

				// Find By Code Product
				request := httptest.NewRequest(http.MethodGet, "/api/products/"+row.Code, nil)
				request.Header.Add("Content-Type", "application/json")
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"].(map[string]interface{})["code"]).To(Equal("SK6"))
				Expect(responseBody["data"].(map[string]interface{})["name"]).To(Equal("Widdy"))
				Expect(responseBody["data"].(map[string]interface{})["description"]).To(Equal("Test"))
			})
		})
	})

	Describe("Access Product Endpoint", func() {
		When("the user is not logged in", func() {
			It("should return error unauthorized response", func() {
				request := httptest.NewRequest(http.MethodGet, "/api/products", nil)
				request.Header.Add("Content-Type", "application/json")

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusUnauthorized))
				Expect(responseBody["status"]).To(Equal("invalid token"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})
	})
})
