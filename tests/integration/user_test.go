package integration_test

import (
	"apriori/entity"
	"apriori/repository"
	"apriori/tests/setup"
	"context"
	"database/sql"
	"encoding/json"
	gin "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"
)

var _ = Describe("User API", func() {

	var server *gin.Engine
	var database *sql.DB

	BeforeEach(func() {
		err := setup.TestEnv()
		if err != nil {
			panic(err)
		}

		db, err := setup.SuiteSetup()
		if err != nil {
			panic(err)
		}

		router := setup.ModuleSetup(db)

		database = db
		server = router
	})

	AfterEach(func() {
		db, err := setup.SuiteSetup()
		if err != nil {
			panic(err)
		}

		err = setup.TearDownTest(db)
		if err != nil {
			panic(err)
		}
	})

	Describe("Create User /users", func() {
		When("the fields are incorrect", func() {
			When("the name field is incorrect", func() {
				It("should return error required", func() {
					requestBody := strings.NewReader(`{"email": "widdy@gmail.com","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/users", requestBody)
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
					requestBody := strings.NewReader(`{"name":"asdasdsdsasdsfsdsassssssssssd", "email": "widdy@gmail.com","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/users", requestBody)
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
					requestBody := strings.NewReader(`{"name": "Widdy","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/users", requestBody)
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
					requestBody := strings.NewReader(`{"name": "Widdy","email": "Widdys","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/users", requestBody)
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
					requestBody := strings.NewReader(`{"name": "Widdy","email": "widdy@gmail.com","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/users", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					// Second register with the same email
					requestBody = strings.NewReader(`{"name": "Widdy","email": "widdy@gmail.com","password": "Rahasia123"}`)
					request = httptest.NewRequest(http.MethodPost, "/api/users", requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer = httptest.NewRecorder()
					server.ServeHTTP(writer, request)
					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusInternalServerError))
					Expect(responseBody["status"]).To(Equal("Error 1062: Duplicate entry 'widdy@gmail.com' for key 'email'"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error exceeds the limit character", func() {
					requestBody := strings.NewReader(`{"name":"wids","email": "Widdysdsdasdddddddsadasdasdss@gmail.com","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/users", requestBody)
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
					requestBody := strings.NewReader(`{"name": "Widdy","email": "widdy@gmail.com"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/users", requestBody)
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
					requestBody := strings.NewReader(`{"name": "Widdy","email": "widdy@gmail.com","password": "as"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/users", requestBody)
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
			It("should return a successful create new user response", func() {
				requestBody := strings.NewReader(`{"name": "Widdy","email": "widdy@gmail.com","password": "Rahasia123"}`)
				request := httptest.NewRequest(http.MethodPost, "/api/auth/register", requestBody)
				request.Header.Add("Content-Type", "application/json")

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("created"))
				Expect(responseBody["data"].(map[string]interface{})["name"]).To(Equal("Widdy"))
				Expect(responseBody["data"].(map[string]interface{})["email"]).To(Equal("widdy@gmail.com"))
			})
		})
	})

	Describe("Update User /users/:id", func() {
		When("user not found", func() {
			It("should return error not found", func() {
				requestBody := strings.NewReader(`{"name": "SiGanteng","email": "ganteng@gmail.com","password":"Widdy123"}`)
				request := httptest.NewRequest(http.MethodPatch, "/api/users/23", requestBody)
				request.Header.Add("Content-Type", "application/json")

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusInternalServerError))
				Expect(responseBody["status"]).To(Equal("user not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("the fields are incorrect", func() {
			When("the name field is incorrect", func() {
				It("should return error required", func() {
					// Create User
					tx, _ := database.Begin()
					userRepository := repository.NewUserRepository()
					password, _ := bcrypt.GenerateFromPassword([]byte("Rahasia123"), bcrypt.DefaultCost)
					row, _ := userRepository.Create(context.Background(), tx, entity.User{
						Name:      "Widdy",
						Email:     "widdy@gmail.com",
						Password:  string(password),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					})
					_ = tx.Commit()

					// Update User
					requestBody := strings.NewReader(`{"email": "widdy@gmail.com","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPatch, "/api/users/"+strconv.Itoa(int(row.IdUser)), requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'UpdateUserRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error exceeds the limit character", func() {
					// Create User
					tx, _ := database.Begin()
					userRepository := repository.NewUserRepository()
					password, _ := bcrypt.GenerateFromPassword([]byte("Rahasia123"), bcrypt.DefaultCost)
					row, _ := userRepository.Create(context.Background(), tx, entity.User{
						Name:      "Widdy",
						Email:     "widdy@gmail.com",
						Password:  string(password),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					})
					_ = tx.Commit()

					requestBody := strings.NewReader(`{"name":"asdasdsdsasdsfsdsassssssssssd", "email": "widdy@gmail.com","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPatch, "/api/users/"+strconv.Itoa(int(row.IdUser)), requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'UpdateUserRequest.Name' Error:Field validation for 'Name' failed on the 'max' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("the email field is incorrect", func() {
				It("should return error required", func() {
					// Create User
					tx, _ := database.Begin()
					userRepository := repository.NewUserRepository()
					password, _ := bcrypt.GenerateFromPassword([]byte("Rahasia123"), bcrypt.DefaultCost)
					row, _ := userRepository.Create(context.Background(), tx, entity.User{
						Name:      "Widdy",
						Email:     "widdy@gmail.com",
						Password:  string(password),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					})
					_ = tx.Commit()

					// Update User
					requestBody := strings.NewReader(`{"name": "Widdy","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPatch, "/api/users/"+strconv.Itoa(int(row.IdUser)), requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'UpdateUserRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error the email must be valid email", func() {
					// Create User
					tx, _ := database.Begin()
					userRepository := repository.NewUserRepository()
					password, _ := bcrypt.GenerateFromPassword([]byte("Rahasia123"), bcrypt.DefaultCost)
					row, _ := userRepository.Create(context.Background(), tx, entity.User{
						Name:      "Widdy",
						Email:     "widdy@gmail.com",
						Password:  string(password),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					})
					_ = tx.Commit()

					// Update User
					requestBody := strings.NewReader(`{"name": "Widdy","email": "Widdys","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPatch, "/api/users/"+strconv.Itoa(int(row.IdUser)), requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'UpdateUserRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})

				It("should return error exceeds the limit character", func() {
					// Create User
					tx, _ := database.Begin()
					userRepository := repository.NewUserRepository()
					password, _ := bcrypt.GenerateFromPassword([]byte("Rahasia123"), bcrypt.DefaultCost)
					row, _ := userRepository.Create(context.Background(), tx, entity.User{
						Name:      "Widdy",
						Email:     "widdy@gmail.com",
						Password:  string(password),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					})
					_ = tx.Commit()

					// Update User
					requestBody := strings.NewReader(`{"name":"wids","email": "Widdysdsdasdddddddsadasdasdss@gmail.com","password": "Rahasia123"}`)
					request := httptest.NewRequest(http.MethodPatch, "/api/users/"+strconv.Itoa(int(row.IdUser)), requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'UpdateUserRequest.Email' Error:Field validation for 'Email' failed on the 'max' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("the password field is incorrect", func() {
				It("should return error less character of length", func() {
					// Create User
					tx, _ := database.Begin()
					userRepository := repository.NewUserRepository()
					password, _ := bcrypt.GenerateFromPassword([]byte("Rahasia123"), bcrypt.DefaultCost)
					row, _ := userRepository.Create(context.Background(), tx, entity.User{
						Name:      "Widdy",
						Email:     "widdy@gmail.com",
						Password:  string(password),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					})
					_ = tx.Commit()

					// Update User
					requestBody := strings.NewReader(`{"name": "Widdy","email": "widdy@gmail.com","password": "as"}`)
					request := httptest.NewRequest(http.MethodPatch, "/api/users/"+strconv.Itoa(int(row.IdUser)), requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'UpdateUserRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
		})

		When("the fields are correct", func() {
			When("password is empty", func() {
				It("should return a successful update user response", func() {
					// Create User
					tx, _ := database.Begin()
					userRepository := repository.NewUserRepository()
					password, _ := bcrypt.GenerateFromPassword([]byte("Rahasia123"), bcrypt.DefaultCost)
					row, _ := userRepository.Create(context.Background(), tx, entity.User{
						Name:      "Widdy",
						Email:     "widdy@gmail.com",
						Password:  string(password),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					})
					_ = tx.Commit()

					// Update User
					requestBody := strings.NewReader(`{"name": "SiGanteng","email": "ganteng@gmail.com"}`)
					request := httptest.NewRequest(http.MethodPatch, "/api/users/"+strconv.Itoa(int(row.IdUser)), requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
					Expect(responseBody["status"]).To(Equal("updated"))
					Expect(responseBody["data"].(map[string]interface{})["name"]).ShouldNot(Equal("Widdy"))
					Expect(responseBody["data"].(map[string]interface{})["email"]).ShouldNot(Equal("widdy@gmail.com"))
				})
			})

			When("the fields are filled", func() {
				It("should return a successful update user response", func() {
					// Create User
					tx, _ := database.Begin()
					userRepository := repository.NewUserRepository()
					password, _ := bcrypt.GenerateFromPassword([]byte("Rahasia123"), bcrypt.DefaultCost)
					row, _ := userRepository.Create(context.Background(), tx, entity.User{
						Name:      "Widdy",
						Email:     "widdy@gmail.com",
						Password:  string(password),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					})
					_ = tx.Commit()

					// Update User
					requestBody := strings.NewReader(`{"name": "SiGanteng","email": "ganteng@gmail.com","password":"Widdy123"}`)
					request := httptest.NewRequest(http.MethodPatch, "/api/users/"+strconv.Itoa(int(row.IdUser)), requestBody)
					request.Header.Add("Content-Type", "application/json")

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					response := writer.Result()

					body, _ := io.ReadAll(response.Body)
					var responseBody map[string]interface{}
					_ = json.Unmarshal(body, &responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
					Expect(responseBody["status"]).To(Equal("updated"))
					Expect(responseBody["data"].(map[string]interface{})["name"]).ShouldNot(Equal("Widdy"))
					Expect(responseBody["data"].(map[string]interface{})["email"]).ShouldNot(Equal("widdy@gmail.com"))
				})
			})
		})
	})

	Describe("Delete User /users/:id", func() {
		When("user not found", func() {
			It("should return error not found", func() {
				request := httptest.NewRequest(http.MethodDelete, "/api/users/23", nil)
				request.Header.Add("Content-Type", "application/json")

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusInternalServerError))
				Expect(responseBody["status"]).To(Equal("user not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("user is found", func() {
			It("should return a successful delete user response", func() {
				// Create User
				tx, _ := database.Begin()
				userRepository := repository.NewUserRepository()
				password, _ := bcrypt.GenerateFromPassword([]byte("Rahasia123"), bcrypt.DefaultCost)
				row, _ := userRepository.Create(context.Background(), tx, entity.User{
					Name:      "Widdy",
					Email:     "widdy@gmail.com",
					Password:  string(password),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
				_ = tx.Commit()

				// Update User
				request := httptest.NewRequest(http.MethodDelete, "/api/users/"+strconv.Itoa(int(row.IdUser)), nil)
				request.Header.Add("Content-Type", "application/json")

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

	Describe("Find ALl User /users", func() {
		When("the user is not present", func() {
			It("should return a successful but the data is null", func() {
				request := httptest.NewRequest(http.MethodGet, "/api/users", nil)
				request.Header.Add("Content-Type", "application/json")

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
		When("the user is present", func() {
			It("should return a successful and show all users", func() {
				// Create User
				tx, _ := database.Begin()
				userRepository := repository.NewUserRepository()
				password, _ := bcrypt.GenerateFromPassword([]byte("Rahasia123"), bcrypt.DefaultCost)
				user1, _ := userRepository.Create(context.Background(), tx, entity.User{
					Name:      "Widdy",
					Email:     "widdy@gmail.com",
					Password:  string(password),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
				user2, _ := userRepository.Create(context.Background(), tx, entity.User{
					Name:      "Agung",
					Email:     "agung@gmail.com",
					Password:  string(password),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
				_ = tx.Commit()

				// Find All User
				request := httptest.NewRequest(http.MethodGet, "/api/users", nil)
				request.Header.Add("Content-Type", "application/json")

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				users := responseBody["data"].([]interface{})

				usersResponse1 := users[0].(map[string]interface{})
				usersResponse2 := users[1].(map[string]interface{})

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))

				Expect(user1.IdUser).To(Equal(uint64(usersResponse1["id_user"].(float64))))
				Expect(user1.Name).To(Equal(usersResponse1["name"]))

				Expect(user2.IdUser).To(Equal(uint64(usersResponse2["id_user"].(float64))))
				Expect(user2.Name).To(Equal(usersResponse2["name"]))
			})
		})
	})

	Describe("Find By Id User /users/:id", func() {
		When("user is not found", func() {
			It("should return error not found", func() {
				// Find By Id User
				request := httptest.NewRequest(http.MethodGet, "/api/users/1", nil)
				request.Header.Add("Content-Type", "application/json")

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusInternalServerError))
				Expect(responseBody["status"]).To(Equal("user not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})
		When("user is found", func() {
			It("should return a successful find user by id", func() {
				// Create User
				tx, _ := database.Begin()
				userRepository := repository.NewUserRepository()
				password, _ := bcrypt.GenerateFromPassword([]byte("Rahasia123"), bcrypt.DefaultCost)
				row, _ := userRepository.Create(context.Background(), tx, entity.User{
					Name:      "Widdy",
					Email:     "widdy@gmail.com",
					Password:  string(password),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
				_ = tx.Commit()

				// Find By Id User
				request := httptest.NewRequest(http.MethodGet, "/api/users/"+strconv.Itoa(int(row.IdUser)), nil)
				request.Header.Add("Content-Type", "application/json")

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"].(map[string]interface{})["name"]).To(Equal("Widdy"))
				Expect(responseBody["data"].(map[string]interface{})["email"]).To(Equal("widdy@gmail.com"))
			})
		})
	})
})
