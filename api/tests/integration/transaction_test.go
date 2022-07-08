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

var _ = Describe("Transaction API", func() {

	var server *gin.Engine
	var database *sql.DB
	var tokenJWT string
	var cookie *http.Cookie

	BeforeEach(func() {
		// Setup Configuration
		configuration := config.New("../../.env.test")

		router, db := setup.ModuleSetup(configuration)

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
		_, db := setup.ModuleSetup(configuration)
		defer db.Close()

		err := setup.TearDownTest(db)
		if err != nil {
			panic(err)
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
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

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
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

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
				Expect(responseBody["data"].(map[string]interface{})["product_name"]).To(Equal("Kasur cinta, Bantal memori"))
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
	//			response := rec.Result()
	//
	//			resp, _ := io.ReadAll(response.Body)
	//			var responseBody map[string]interface{}
	//			_ = json.Unmarshal(resp, &responseBody)
	//
	//			log.Println(responseBody["status"])
	//			Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
	//			Expect(responseBody["status"]).To(Equal("created"))
	//			Expect(responseBody["data"]).To(BeNil())
	//		})
	//	})
	//})

	Describe("Update Transaction /transactions/:no_transaction", func() {
		When("the fields are incorrect", func() {
			When("the product name field is incorrect", func() {
				It("should return error required", func() {
					// Create Transaction
					tx, _ := database.Begin()
					transactionRepository := repository.NewTransactionRepository()
					row, _ := transactionRepository.Create(context.Background(), tx, entity.Transaction{
						ProductName:   "Kasur cinta, Bantal memori",
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
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

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
					row, _ := transactionRepository.Create(context.Background(), tx, entity.Transaction{
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
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

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
				row, _ := transactionRepository.Create(context.Background(), tx, entity.Transaction{
					ProductName:   "Kasur cinta, Bantal memori",
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
				Expect(responseBody["data"].(map[string]interface{})["product_name"]).ShouldNot(Equal("Kasur cinta, Bantal memori"))
				Expect(responseBody["data"].(map[string]interface{})["customer_name"]).ShouldNot(Equal("Wids"))
			})
		})
	})

	Describe("Delete Transaction /transactions/:no_transaction", func() {
		When("transaction is not found", func() {
			It("should return error not found", func() {
				// Delete Transaction
				request := httptest.NewRequest(http.MethodDelete, "/api/transactions/32412", nil)
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
				Expect(responseBody["status"]).To(Equal("transaction not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("transaction is found", func() {
			It("should return a successful delete transaction response", func() {
				// Create Transaction
				tx, _ := database.Begin()
				transactionRepository := repository.NewTransactionRepository()
				row, _ := transactionRepository.Create(context.Background(), tx, entity.Transaction{
					ProductName:   "Kasur cinta, Bantal memori",
					CustomerName:  "Wids",
					NoTransaction: "202320",
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				})
				_ = tx.Commit()

				// Delete Transaction
				request := httptest.NewRequest(http.MethodDelete, "/api/transactions/"+row.NoTransaction, nil)
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

	Describe("Find All Transaction /transactions", func() {
		When("the transaction is not present", func() {
			It("should return a successful but the data is null", func() {
				// Find All Transaction
				request := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
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
		When("the transactions is present", func() {
			It("should return a successful and show all transactions", func() {
				// Create Transaction
				tx, _ := database.Begin()
				transactionRepository := repository.NewTransactionRepository()
				transaction1, _ := transactionRepository.Create(context.Background(), tx, entity.Transaction{
					ProductName:  "Kasur cinta, Bantal memori",
					CustomerName: "Wids",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				})
				transaction2, _ := transactionRepository.Create(context.Background(), tx, entity.Transaction{
					ProductName:  "Guling cinta, Guling memori",
					CustomerName: "Goengs",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				})
				_ = tx.Commit()

				// Find All Transaction
				request := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
				request.Header.Add("Content-Type", "application/json")
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				transactions := responseBody["data"].([]interface{})

				// Desc
				transactionResponse1 := transactions[1].(map[string]interface{})
				transactionResponse2 := transactions[0].(map[string]interface{})

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))

				Expect(transaction1.ProductName).To(Equal(transactionResponse1["product_name"]))
				Expect(transaction1.CustomerName).To(Equal(transactionResponse1["customer_name"]))

				Expect(transaction2.ProductName).To(Equal(transactionResponse2["product_name"]))
				Expect(transaction2.CustomerName).To(Equal(transactionResponse2["customer_name"]))
			})
		})
	})

	Describe("Find By No Transaction /transactions/:no_transaction", func() {
		When("transaction is not found", func() {
			It("should return error not found", func() {
				// Find By No Transaction
				request := httptest.NewRequest(http.MethodGet, "/api/transactions/52324", nil)
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
				Expect(responseBody["status"]).To(Equal("transaction not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})
		When("transaction is found", func() {
			It("should return a successful find transaction by no transaction", func() {
				// Create Transaction
				tx, _ := database.Begin()
				transactionRepository := repository.NewTransactionRepository()
				row, _ := transactionRepository.Create(context.Background(), tx, entity.Transaction{
					ProductName:   "Kasur cinta, Bantal memori",
					CustomerName:  "Wids",
					NoTransaction: "202320",
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				})
				_ = tx.Commit()

				// Find By No Transaction
				request := httptest.NewRequest(http.MethodGet, "/api/transactions/"+row.NoTransaction, nil)
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
				Expect(responseBody["data"].(map[string]interface{})["product_name"]).To(Equal("Kasur cinta, Bantal memori"))
				Expect(responseBody["data"].(map[string]interface{})["customer_name"]).To(Equal("Wids"))
			})
		})
	})

	Describe("Access Transaction Endpoint", func() {
		When("the user is not logged in", func() {
			It("should return error unauthorized response", func() {
				request := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
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
