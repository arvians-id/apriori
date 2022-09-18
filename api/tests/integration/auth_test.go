package integration_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/arvians-id/apriori/config"
	"github.com/arvians-id/apriori/entity"
	repository "github.com/arvians-id/apriori/repository/postgres"
	"github.com/arvians-id/apriori/tests/setup"
	"strings"
	"time"

	"log"

	gin "github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"net/http"
	"net/http/httptest"
)

var _ = Describe("Auth API", func() {
	var server *gin.Engine
	var database *sql.DB
	configuration := config.New("../../.env.test")

	BeforeEach(func() {
		router, db := setup.ModuleSetup(configuration)
		database = db
		server = router
	})

	AfterEach(func() {
		_, db := setup.ModuleSetup(configuration)
		defer db.Close()
		err := setup.TearDownTest(db)
		if err != nil {
			log.Fatal(err)
		}
	})

	Describe("Login /auth/login", func() {
		When("the fields is incorrect", func() {
			When("the email field is incorrect", func() {
				It("should return error required", func() {
					requestBody := strings.NewReader(`{"password":"Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/auth/login", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'GetUserCredentialRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error the email must be valid email", func() {
					requestBody := strings.NewReader(`{"email":"Widdys","password":"Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/auth/login", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'GetUserCredentialRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("the email is not found", func() {
				It("should return error not found", func() {
					requestBody := strings.NewReader(`{"email": "widdy@gmail.com","password":"Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/auth/login", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("the password field is incorrect", func() {
				It("should return error required", func() {
					requestBody := strings.NewReader(`{"email": "widdy@gmail.com"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/auth/login", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'GetUserCredentialRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("the password is wrong", func() {
				It("should return error wrong password", func() {
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
					requestBody := strings.NewReader(`{"email": "widdy@gmail.com","password":"Raha123"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/auth/login", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("wrong password"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
		})

		When("the fields is correct", func() {
			It("should return successful login response", func() {
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

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"].(map[string]interface{})["access_token"]).ShouldNot(BeNil())
				Expect(responseBody["data"].(map[string]interface{})["refresh_token"]).ShouldNot(BeNil())
			})
		})
	})

	Describe("Refresh Token /auth/refresh", func() {
		When("the refresh token is incorrect", func() {
			It("should return error invalid token", func() {
				requestBody := strings.NewReader(`{"refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZF91c2VyIjoxLCJleHAiOjE2NTM5MjI1MTJ9.6xJ4ZQdem4ZoWPBuZctJTMKNOkqE93Ea0ncKovpqN44"}`)
				request := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusInternalServerError))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("the refresh token is correct", func() {
			It("should regenerate a new access token and refresh token", func() {
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

				// Refresh token
				signature := responseBody["data"].(map[string]interface{})["refresh_token"].(string)
				sign := fmt.Sprintf(`{"refresh_token": "%s"}`, signature)
				requestBody = strings.NewReader(sign)
				request = httptest.NewRequest(http.MethodPost, "/api/auth/refresh", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"].(map[string]interface{})["access_token"]).ShouldNot(BeNil())
				Expect(responseBody["data"].(map[string]interface{})["refresh_token"]).ShouldNot(BeNil())
			})
		})
	})

	Describe("Logout /auth/logout", func() {
		When("the token is correct", func() {
			It("should delete cookies and cannot log in", func() {
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
				cookies := writer.Result().Cookies()

				Expect(cookies[0].Value).ShouldNot(BeNil())
				// Logout
				request = httptest.NewRequest(http.MethodDelete, "/api/auth/logout", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)
				cookies = writer.Result().Cookies()

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(cookies[0].Value).To(Equal(""))
				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})
	})

	Describe("Forgot Password /auth/forgot-password", func() {
		When("the email field is incorrect", func() {
			It("should return error required", func() {
				requestBody := strings.NewReader(`{}`)
				request := httptest.NewRequest(http.MethodPost, "/api/auth/forgot-password", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
				Expect(responseBody["status"]).To(Equal("Key: 'CreatePasswordResetRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"))
				Expect(responseBody["data"]).To(BeNil())
			})

			It("should return error the email must be valid email", func() {
				requestBody := strings.NewReader(`{"email": "Widdys"}`)
				request := httptest.NewRequest(http.MethodPost, "/api/auth/forgot-password", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
				Expect(responseBody["status"]).To(Equal("Key: 'CreatePasswordResetRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("the email field is correct", func() {
			It("should return success response and send email to user", func() {
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

				// Send forgot password
				requestBody := strings.NewReader(`{"email": "widdy@gmail.com"}`)
				request := httptest.NewRequest(http.MethodPost, "/api/auth/forgot-password", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("mail sent successfully"))
				Expect(responseBody["data"].(map[string]interface{})["signature"]).ShouldNot(BeNil())
			})
		})
	})

	Describe("Verify Password /auth/verify", func() {
		When("the fields is incorrect", func() {
			When("the email field is incorrect", func() {
				It("should return error required", func() {
					requestBody := strings.NewReader(`{"password": "Widdy123"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/auth/verify", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'UpdateResetPasswordUserRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error the email must be valid email", func() {
					requestBody := strings.NewReader(`{"email": "widdyarfiansyah","password": "Widdy123"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/auth/verify", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'UpdateResetPasswordUserRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("the password field is incorrect", func() {
				It("should return error required", func() {
					requestBody := strings.NewReader(`{"email": "widdyarfiansyah@ummi.ac.id"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/auth/verify", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'UpdateResetPasswordUserRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error less character of length", func() {
					requestBody := strings.NewReader(`{"email": "widdyarfiansyah@ummi.ac.id","password": "Wi"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/auth/verify", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'UpdateResetPasswordUserRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("giving the incorrect signature", func() {
				It("should return error invalid credentials", func() {
					// Create user first
					requestBody := strings.NewReader(`{"name": "Widdy","email": "widdy@gmail.com","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/auth/verify", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					// Send forgot password
					requestBody = strings.NewReader(`{"email": "widdy@gmail.com"}`)
					request = httptest.NewRequest(http.MethodPost, "/api/auth/forgot-password", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

					writer = httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					// Verify password
					requestBody = strings.NewReader(`{"email": "widdy@gmail.com","password": "Widdy123"}`)
					request = httptest.NewRequest(http.MethodPost, "/api/auth/verify?signature=asdsa23sda", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

					writer = httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
		})

		When("the field is correct", func() {
			It("should return success response and update password's user", func() {
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

				// Send forgot password
				requestBody := strings.NewReader(`{"email": "widdy@gmail.com"}`)
				request := httptest.NewRequest(http.MethodPost, "/api/auth/forgot-password", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				// Verify password
				signature := responseBody["data"].(map[string]interface{})["signature"].(string)
				requestBody = strings.NewReader(`{"email": "widdy@gmail.com","password": "Widdy123"}`)
				request = httptest.NewRequest(http.MethodPost, "/api/auth/verify?signature="+signature, requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("updated"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})
	})

	Describe("User Registration /auth/register", func() {
		When("the fields is incorrect", func() {
			When("the name field is incorrect", func() {
				It("should return error required", func() {
					requestBody := strings.NewReader(`{"email": "widdy@gmail.com","address":"nganjok","phone":"082299","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/auth/register", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateUserRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error exceeds the limit character", func() {
					requestBody := strings.NewReader(`{"name":"asdasdsdsasdsfsdsassssssssssd","email": "widdy@gmail.com","address":"nganjok","phone":"082299","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/auth/register", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateUserRequest.Name' Error:Field validation for 'Name' failed on the 'max' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("the email field is incorrect", func() {
				It("should return error required", func() {
					requestBody := strings.NewReader(`{"name": "Widdy","address":"nganjok","phone":"082299","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/auth/register", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error the email must be valid email", func() {
					requestBody := strings.NewReader(`{"name":"Widdy","email":"Widdys","address":"nganjok","phone":"082299","password":"Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/auth/register", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error duplicate email", func() {
					// First register
					requestBody := strings.NewReader(`{"name": "Widdy","email": "widdy@gmail.com","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/auth/register", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					// Second register with the same email
					requestBody = strings.NewReader(`{"name": "Widdy","email": "widdy@gmail.com","password": "Rahasia123"}`)
					request = httptest.NewRequest(http.MethodPost, "/api/auth/register", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

					writer = httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error exceeds the limit character", func() {
					requestBody := strings.NewReader(`{"name":"wids","email": "sssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssddddddddddddddddddddddddddddddddddddddd@gmail.com","address":"nganjok","phone":"082299","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/auth/register", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'max' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("the password field is incorrect", func() {
				It("should return error required", func() {
					requestBody := strings.NewReader(`{"name":"Widdy","email":"widdy@gmail.com","address":"nganjok","phone":"082299"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/auth/register", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error less character of length", func() {
					requestBody := strings.NewReader(`{"name":"Widdy","email":"widdy@gmail.com","address":"nganjok","phone":"082299","password":"as"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/auth/register", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
		})

		When("the fields are correct", func() {
			It("should return a successful register response", func() {
				requestBody := strings.NewReader(`{"name": "Widdy","email": "widdy@gmail.com","address":"nganjok","phone":"082299","password": "Rahasia123"}`)
				request := httptest.NewRequest(http.MethodPost, "/api/auth/register", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("created"))
				Expect(responseBody["data"].(map[string]interface{})["name"]).To(Equal("Widdy"))
				Expect(responseBody["data"].(map[string]interface{})["email"]).To(Equal("widdy@gmail.com"))
			})
		})
	})

	Describe("Check Token /auth/token", func() {
		When("the token is empty or invalid", func() {
			It("Should return error unauthorized", func() {
				request := httptest.NewRequest(http.MethodGet, "/api/auth/token", nil)
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

		When("the token is exists", func() {
			It("Should return successful get token", func() {
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

				requestBody := strings.NewReader(`{"email": "widdy@gmail.com", "password": "Rahasia123"}`)
				request := httptest.NewRequest(http.MethodPost, "/api/auth/login", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				var cookie *http.Cookie
				for _, c := range writer.Result().Cookies() {
					if c.Name == "token" {
						cookie = c
					}
				}
				bearerToken := fmt.Sprintf("Bearer %v", responseBody["data"].(map[string]interface{})["access_token"].(string))

				request = httptest.NewRequest(http.MethodGet, "/api/auth/token", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
				request.Header.Set("Authorization", bearerToken)
				request.AddCookie(cookie)

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})
	})
})
