package integration_test

import (
	"apriori/api/controller"
	"apriori/repository"
	"apriori/service"
	"database/sql"
	"encoding/json"
	"fmt"
	gin "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"io"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"net/http"
	"net/http/httptest"
)

func TestUserController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth controller test")
}

var _ = Describe("Auth Controller", func() {

	var server *gin.Engine

	BeforeEach(func() {
		db, err := sql.Open("mysql", "root@tcp(localhost:3306)/apriori_test?parseTime=true")
		if err != nil {
			panic(err)
		}

		db.SetMaxIdleConns(5)
		db.SetMaxOpenConns(20)
		db.SetConnMaxLifetime(60 * time.Minute)
		db.SetConnMaxIdleTime(10 * time.Minute)

		err = godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		router := gin.New()

		userRepository := repository.NewUserRepository()
		authRepository := repository.NewAuthRepository()
		passwordRepository := repository.NewPasswordResetRepository()

		userService := service.NewUserService(&userRepository, db)
		authService := service.NewAuthService(&userRepository, &authRepository, db)
		jwtService := service.NewJwtService()
		emailService := service.NewEmailService()
		passwordResetService := service.NewPasswordResetService(&passwordRepository, &userRepository, db)

		authController := controller.NewAuthController(&authService, &userService, jwtService, emailService, &passwordResetService)

		authController.Route(router)

		server = router
	})

	AfterEach(func() {
		db, err := sql.Open("mysql", "root@tcp(localhost:3306)/apriori_test")
		if err != nil {
			panic(err)
		}

		_, err = db.Exec(`TRUNCATE TABLE users;`)
		_, err = db.Exec(`TRUNCATE TABLE password_resets;`)

		if err != nil {
			panic(err)
		}
	})

	Describe("Login /auth/login", func() {
		When("the fields is incorrect", func() {
			When("the email field is incorrect", func() {
				It("should return error required", func() {
					requestBody := strings.NewReader(`{"password":"Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/login", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'GetUserCredentialRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
				It("should return error the email must be valid email", func() {
					requestBody := strings.NewReader(`{"email":"agungs","password":"Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/login", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'GetUserCredentialRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
			When("the email is not found", func() {
				It("should return error not found", func() {
					requestBody := strings.NewReader(`{"email": "agungs@gmail.com","password":"Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/login", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusInternalServerError))
					Expect(responseBody["status"]).To(Equal("email not found"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
			When("the password field is incorrect", func() {
				It("should return error required", func() {
					requestBody := strings.NewReader(`{"email": "agungs@gmail.com"}`)
					request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/login", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'GetUserCredentialRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
			When("the password is wrong", func() {
				It("should return error wrong password", func() {
					// Create user first
					requestBody := strings.NewReader(`{"name": "Agung","email": "agungs@gmail.com","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/register", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					// Login
					requestBody = strings.NewReader(`{"email": "agungs@gmail.com","password":"Raha123"}`)
					request = httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/login", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer = httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusInternalServerError))
					Expect(responseBody["status"]).To(Equal("wrong password"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
		})

		When("the fields is correct", func() {
			It("should return successful login response", func() {
				// Create user first
				requestBody := strings.NewReader(`{"name": "Agung","email": "agungs@gmail.com","password": "Rahasia123"}`)
				request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/register", requestBody)
				request.Header.Add("Content-Type", "application/json")

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				// Login
				requestBody = strings.NewReader(`{"email": "agungs@gmail.com","password":"Rahasia123"}`)
				request = httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/login", requestBody)
				request.Header.Add("Content-Type", "application/json")

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"].(map[string]interface{})["access_token"]).ShouldNot(BeNil())
				Expect(responseBody["data"].(map[string]interface{})["refresh_token"]).ShouldNot(BeNil())

				cookies := response.Cookies()

				Expect(cookies[0].Value).ShouldNot(BeNil())
			})
		})
	})

	Describe("Refresh Token /auth/refresh", func() {
		When("the refresh token is incorrect", func() {
			It("should return error invalid token", func() {
				requestBody := strings.NewReader(`{"refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZF91c2VyIjoxLCJleHAiOjE2NTM5MjI1MTJ9.6xJ4ZQdem4ZoWPBuZctJTMKNOkqE93Ea0ncKovpqN44"}`)
				request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/refresh", requestBody)
				request.Header.Add("Content-Type", "application/json")

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusInternalServerError))
				Expect(responseBody["data"]).To(BeNil())
			})
		})
		When("the refresh token is correct", func() {
			It("should regenerate a new access token and refresh token", func() {
				// Create user first
				requestBody := strings.NewReader(`{"name": "Agung","email": "agungs@gmail.com","password": "Rahasia123"}`)
				request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/register", requestBody)
				request.Header.Add("Content-Type", "application/json")

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				// Login
				requestBody = strings.NewReader(`{"email": "agungs@gmail.com","password":"Rahasia123"}`)
				request = httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/login", requestBody)
				request.Header.Add("Content-Type", "application/json")

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				// Refresh token
				signature := responseBody["data"].(map[string]interface{})["refresh_token"].(string)
				sign := fmt.Sprintf(`{"refresh_token": "%s"}`, signature)
				requestBody = strings.NewReader(sign)
				request = httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/refresh", requestBody)
				request.Header.Add("Content-Type", "application/json")

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response = writer.Result()

				body, _ = io.ReadAll(response.Body)
				_ = json.Unmarshal(body, &responseBody)

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
				// Create user first
				requestBody := strings.NewReader(`{"name": "Agung","email": "agungs@gmail.com","password": "Rahasia123"}`)
				request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/register", requestBody)
				request.Header.Add("Content-Type", "application/json")

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				// Login
				requestBody = strings.NewReader(`{"email": "agungs@gmail.com","password":"Rahasia123"}`)
				request = httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/login", requestBody)
				request.Header.Add("Content-Type", "application/json")

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()
				cookies := response.Cookies()

				Expect(cookies[0].Value).ShouldNot(BeNil())

				// Logout
				request = httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/auth/logout", nil)
				request.Header.Add("Content-Type", "application/json")

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response = writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				cookies = response.Cookies()

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
				request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/forgot-password", requestBody)
				request.Header.Add("Content-Type", "application/json")

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
				Expect(responseBody["status"]).To(Equal("Key: 'CreatePasswordResetRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"))
				Expect(responseBody["data"]).To(BeNil())
			})

			It("should return error the email must be valid email", func() {
				requestBody := strings.NewReader(`{"email": "agungs"}`)
				request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/forgot-password", requestBody)
				request.Header.Add("Content-Type", "application/json")

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
				Expect(responseBody["status"]).To(Equal("Key: 'CreatePasswordResetRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("the email field is correct", func() {
			It("should return success response and send email to user", func() {
				// Create user first
				requestBody := strings.NewReader(`{"name": "Agung","email": "agungs@gmail.com","password": "Rahasia123"}`)
				request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/register", requestBody)
				request.Header.Add("Content-Type", "application/json")

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				// Send forgot password
				requestBody = strings.NewReader(`{"email": "agungs@gmail.com"}`)
				request = httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/forgot-password", requestBody)
				request.Header.Add("Content-Type", "application/json")

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

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
					request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/verify", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'UpdateResetPasswordUserRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error the email must be valid email", func() {
					requestBody := strings.NewReader(`{"email": "widdyarfiansyah","password": "Widdy123"}`)
					request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/verify", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'UpdateResetPasswordUserRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("the password field is incorrect", func() {
				It("should return error required", func() {
					requestBody := strings.NewReader(`{"email": "widdyarfiansyah@ummi.ac.id"}`)
					request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/verify", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'UpdateResetPasswordUserRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error less character of length", func() {
					requestBody := strings.NewReader(`{"email": "widdyarfiansyah@ummi.ac.id","password": "Wi"}`)
					request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/verify", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'UpdateResetPasswordUserRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("giving the incorrect signature", func() {
				It("should return error invalid credentials", func() {
					// Create user first
					requestBody := strings.NewReader(`{"name": "Agung","email": "agungs@gmail.com","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/verify", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					// Send forgot password
					requestBody = strings.NewReader(`{"email": "agungs@gmail.com"}`)
					request = httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/forgot-password", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer = httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					// Verify password
					requestBody = strings.NewReader(`{"email": "agungs@gmail.com","password": "Widdy123"}`)
					request = httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/verify?signature=asdsa23sda", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer = httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusInternalServerError))
					Expect(responseBody["status"]).To(Equal("invalid credentials"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
		})

		When("the field is correct", func() {
			It("should return success response and update password's user", func() {
				// Create user first
				requestBody := strings.NewReader(`{"name": "Agung","email": "agungs@gmail.com","password": "Rahasia123"}`)
				request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/register", requestBody)
				request.Header.Add("Content-Type", "application/json")

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				// Send forgot password
				requestBody = strings.NewReader(`{"email": "agungs@gmail.com"}`)
				request = httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/forgot-password", requestBody)
				request.Header.Add("Content-Type", "application/json")

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()
				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				// Verify password
				signature := responseBody["data"].(map[string]interface{})["signature"].(string)
				requestBody = strings.NewReader(`{"email": "agungs@gmail.com","password": "Widdy123"}`)
				request = httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/verify?signature="+signature, requestBody)
				request.Header.Add("Content-Type", "application/json")

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response = writer.Result()
				body, _ = io.ReadAll(response.Body)
				_ = json.Unmarshal(body, &responseBody)

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
					requestBody := strings.NewReader(`{"email": "agungs@gmail.com","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/register", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateUserRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error exceeds the limit character", func() {
					requestBody := strings.NewReader(`{"name":"asdasdsdsasdsfsdsassssssssssd", "email": "agungs@gmail.com","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/register", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateUserRequest.Name' Error:Field validation for 'Name' failed on the 'max' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("the email field is incorrect", func() {
				It("should return error required", func() {
					requestBody := strings.NewReader(`{"name": "Agung","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/register", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error the email must be valid email", func() {
					requestBody := strings.NewReader(`{"name": "Agung","email": "agungs","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/register", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error duplicate email", func() {
					// First register
					requestBody := strings.NewReader(`{"name": "Agung","email": "agungs@gmail.com","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/register", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					// Second register with the same email
					requestBody = strings.NewReader(`{"name": "Agung","email": "agungs@gmail.com","password": "Rahasia123"}`)
					request = httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/register", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer = httptest.NewRecorder()
					server.ServeHTTP(writer, request)
					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusInternalServerError))
					Expect(responseBody["status"]).To(Equal("Error 1062: Duplicate entry 'agungs@gmail.com' for key 'email'"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error exceeds the limit character", func() {
					requestBody := strings.NewReader(`{"name":"wids","email": "agungsdsdasdddddddsadasdasdss@gmail.com","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/register", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'max' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("the password field is incorrect", func() {
				It("should return error required", func() {
					requestBody := strings.NewReader(`{"name": "Agung","email": "agungs@gmail.com"}`)
					request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/register", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error less character of length", func() {
					requestBody := strings.NewReader(`{"name": "Agung","email": "agungs@gmail.com","password": "as"}`)
					request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/register", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
		})

		When("the fields are correct", func() {
			It("Should return a successful register response", func() {
				requestBody := strings.NewReader(`{"name": "Agung","email": "agungs@gmail.com","password": "Rahasia123"}`)
				request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/auth/register", requestBody)
				request.Header.Add("Content-Type", "application/json")

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("created"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})
	})

})
