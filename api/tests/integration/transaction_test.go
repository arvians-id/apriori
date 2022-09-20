package integration

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/arvians-id/apriori/config"
	"github.com/arvians-id/apriori/model"
	repository "github.com/arvians-id/apriori/repository/postgres"
	"github.com/arvians-id/apriori/service"
	"github.com/arvians-id/apriori/tests/setup"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
)

/*
	Error :
		- Create Transactions By CSV File /transactions/csv
*/
var _ = Describe("Transaction API", func() {
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
		_, _ = userRepository.Create(context.Background(), tx, &model.User{
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

		transactionCache := service.NewCacheService(configuration)
		_ = transactionCache.FlushDB(context.Background())

		err := setup.TearDownTest(db)
		if err != nil {
			log.Fatal(err)
		}
	})

	Describe("Create Transaction /transactions", func() {
		When("the fields are incorrect", func() {
			When("the product name field is incorrect", func() {
				It("should return error required", func() {
					// Create Transaction
					requestBody := strings.NewReader(`{"customer_name": "Wids"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/transactions", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateTransactionRequest.ProductName' Error:Field validation for 'ProductName' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("the customer name field is incorrect", func() {
				It("should return error required", func() {
					// Create Transaction
					requestBody := strings.NewReader(`{"product_name": "Kasur cinta, Bantal memori"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/transactions", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateTransactionRequest.CustomerName' Error:Field validation for 'CustomerName' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
		})

		When("the fields are correct", func() {
			It("should return successful create transaction response", func() {
				// Create Transaction
				requestBody := strings.NewReader(`{"product_name": "Kasur cinta, Bantal memori","customer_name": "Wids"}`)
				request := httptest.NewRequest(http.MethodPost, "/api/transactions", requestBody)
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
				Expect(responseBody["data"].(map[string]interface{})["product_name"]).To(Equal("kasur cinta, bantal memori"))
				Expect(responseBody["data"].(map[string]interface{})["customer_name"]).To(Equal("Wids"))
			})
		})
	})

	//Describe("Create Transactions By CSV File /transactions/csv", func() {
	//	When("file exist", func() {
	//		It("should return error no such file", func() {
	//			path := "./assets/example1.csv"
	//			body := new(bytes.Buffer)
	//			writer := multipart.NewWriter(body)
	//			part, _ := writer.CreateFormFile("file", path)
	//			sample, _ := os.Open(path)
	//
	//			_, _ = io.Copy(part, sample)
	//			writer.Close()
	//
	//			// Create Transaction
	//			request := httptest.NewRequest(http.MethodPost, "/api/transactions/csv", body)
	//			request.Header.Add("Content-Type", writer.FormDataContentType())
	//			request.AddCookie(cookie)
	//			request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))
	//
	//			rec := httptest.NewRecorder()
	//			server.ServeHTTP(rec, request)
	//
	//			var responseBody map[string]interface{}
	//			_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)
	//
	//			Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
	//			Expect(responseBody["status"]).To(Equal("created"))
	//			Expect(responseBody["data"]).To(BeNil())
	//		})
	//	})
	//})

	Describe("Update Transaction /transactions/:number_transaction", func() {
		When("the fields are incorrect", func() {
			When("the product name field is incorrect", func() {
				It("should return error required", func() {
					// Create Transaction
					tx, _ := database.Begin()
					transactionRepository := repository.NewTransactionRepository()
					row, _ := transactionRepository.Create(context.Background(), tx, &model.Transaction{
						ProductName:   "kasur cinta, bantal memori",
						CustomerName:  "Wids",
						NoTransaction: "202320",
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
					})
					_ = tx.Commit()

					// Update Transaction
					requestBody := strings.NewReader(`{"customer_name": "Wids"}`)
					request := httptest.NewRequest(http.MethodPatch, "/api/transactions/"+row.NoTransaction, requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'UpdateTransactionRequest.ProductName' Error:Field validation for 'ProductName' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("the customer name field is incorrect", func() {
				It("should return error required", func() {
					// Create Transaction
					tx, _ := database.Begin()
					transactionRepository := repository.NewTransactionRepository()
					row, _ := transactionRepository.Create(context.Background(), tx, &model.Transaction{
						ProductName:   "Kasur cinta, Bantal memori",
						CustomerName:  "Wids",
						NoTransaction: "202320",
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
					})
					_ = tx.Commit()

					// Update Transaction
					requestBody := strings.NewReader(`{"product_name": "Kasur cinta, Bantal memori"}`)
					request := httptest.NewRequest(http.MethodPatch, "/api/transactions/"+row.NoTransaction, requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'UpdateTransactionRequest.CustomerName' Error:Field validation for 'CustomerName' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
		})

		When("the fields are correct", func() {
			It("should return successful update transaction response", func() {
				// Create Transaction
				tx, _ := database.Begin()
				transactionRepository := repository.NewTransactionRepository()
				row, _ := transactionRepository.Create(context.Background(), tx, &model.Transaction{
					ProductName:   "kasur cinta, bantal memori",
					CustomerName:  "Wids",
					NoTransaction: "202320",
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				})
				_ = tx.Commit()

				// Update Transaction
				requestBody := strings.NewReader(`{"product_name": "Guling cinta, Guling memori","customer_name": "Goengs"}`)
				request := httptest.NewRequest(http.MethodPatch, "/api/transactions/"+row.NoTransaction, requestBody)
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
				Expect(responseBody["data"].(map[string]interface{})["product_name"]).To(Equal("guling cinta, guling memori"))
				Expect(responseBody["data"].(map[string]interface{})["product_name"]).ShouldNot(Equal("kasur cinta, bantal memori"))
				Expect(responseBody["data"].(map[string]interface{})["customer_name"]).ShouldNot(Equal("Wids"))
			})
		})
	})

	Describe("Delete Transaction /transactions/:number_transaction", func() {
		When("transaction is not found", func() {
			It("should return error not found", func() {
				// Delete Transaction
				request := httptest.NewRequest(http.MethodDelete, "/api/transactions/32412", nil)
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

		When("transaction is found", func() {
			It("should return a successful delete transaction response", func() {
				// Create Transaction
				tx, _ := database.Begin()
				transactionRepository := repository.NewTransactionRepository()
				row, _ := transactionRepository.Create(context.Background(), tx, &model.Transaction{
					ProductName:   "kasur cinta, bantal memori",
					CustomerName:  "Wids",
					NoTransaction: "202320",
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				})
				_ = tx.Commit()

				// Delete Transaction
				request := httptest.NewRequest(http.MethodDelete, "/api/transactions/"+row.NoTransaction, nil)
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

	Describe("Find All Transaction /transactions", func() {
		When("the transaction is not present", func() {
			It("should return a successful but the data is null", func() {
				// Find All Transaction
				request := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
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

		When("the transactions is present", func() {
			It("should return a successful and show all transactions", func() {
				// Create Transaction
				tx, _ := database.Begin()
				transactionRepository := repository.NewTransactionRepository()
				transaction1, _ := transactionRepository.Create(context.Background(), tx, &model.Transaction{
					ProductName:  "kasur cinta, bantal memori",
					CustomerName: "Wids",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				})
				transaction2, _ := transactionRepository.Create(context.Background(), tx, &model.Transaction{
					ProductName:  "guling cinta, guling memori",
					CustomerName: "Goengs",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				})
				_ = tx.Commit()

				// Find All Transaction
				request := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
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

				transactions := responseBody["data"].([]interface{})
				Expect(transaction1.ProductName).To(Equal(transactions[1].(map[string]interface{})["product_name"]))
				Expect(transaction1.CustomerName).To(Equal(transactions[1].(map[string]interface{})["customer_name"]))

				Expect(transaction2.ProductName).To(Equal(transactions[0].(map[string]interface{})["product_name"]))
				Expect(transaction2.CustomerName).To(Equal(transactions[0].(map[string]interface{})["customer_name"]))
			})
		})
	})

	Describe("Find By No Transaction /transactions/:number_transaction", func() {
		When("transaction is not found", func() {
			It("should return error not found", func() {
				// Find By No Transaction
				request := httptest.NewRequest(http.MethodGet, "/api/transactions/52324", nil)
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

		When("transaction is found", func() {
			It("should return a successful find transaction by no transaction", func() {
				// Create Transaction
				tx, _ := database.Begin()
				transactionRepository := repository.NewTransactionRepository()
				row, _ := transactionRepository.Create(context.Background(), tx, &model.Transaction{
					ProductName:   "kasur cinta, bantal memori",
					CustomerName:  "Wids",
					NoTransaction: "202320",
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				})
				_ = tx.Commit()

				// Find By No Transaction
				request := httptest.NewRequest(http.MethodGet, "/api/transactions/"+row.NoTransaction, nil)
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
				Expect(responseBody["data"].(map[string]interface{})["product_name"]).To(Equal("kasur cinta, bantal memori"))
				Expect(responseBody["data"].(map[string]interface{})["customer_name"]).To(Equal("Wids"))
			})
		})
	})

	Describe("Truncate Transaction /transactions/truncate", func() {
		When("transaction is found", func() {
			It("should return successful delete all transactions", func() {
				// Create Transaction
				tx, _ := database.Begin()
				transactionRepository := repository.NewTransactionRepository()
				_, _ = transactionRepository.Create(context.Background(), tx, &model.Transaction{
					ProductName:   "kasur cinta, bantal memori",
					CustomerName:  "Wids",
					NoTransaction: "202320",
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				})
				_ = tx.Commit()

				// Delete Transaction
				request := httptest.NewRequest(http.MethodDelete, "/api/transactions/truncate", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["data"]).To(BeNil())
			})
		})
	})

	Describe("Access Transaction Endpoint", func() {
		When("the user is not logged in", func() {
			It("should return error unauthorized response", func() {
				request := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusUnauthorized))
				Expect(responseBody["status"]).To(Equal("invalid token"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})
	})
})
