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
	"os"
	"strings"
	"time"
)

var _ = Describe("Comment API", func() {
	var server *gin.Engine
	var database *sql.DB
	var tokenJWT string
	var cookie *http.Cookie
	var order *entity.UserOrder
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
		user, _ := userRepository.Create(context.Background(), tx, &entity.User{
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

		// Create product
		tx, _ = database.Begin()
		productRepository := repository.NewProductRepository()
		_, _ = productRepository.Create(context.Background(), tx, &entity.Product{
			Code:        "Lfanp",
			Name:        "Bantal Biasa",
			Description: "Test Bang",
			Category:    "Bantal, Kasur",
			Mass:        1000,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})

		// Create payload
		payloadRepository := repository.NewPaymentRepository()
		payload, _ := payloadRepository.Create(context.Background(), tx, &entity.Payment{
			UserId: sql.NullString{
				String: helper.IntToStr(user.IdUser),
				Valid:  true,
			},
		})

		// Create User Order
		userOrderRepository := repository.NewUserOrderRepository()
		userOrder, _ := userOrderRepository.Create(context.Background(), tx, &entity.UserOrder{
			PayloadId:      payload.IdPayload,
			Code:           "aXksCj2",
			Name:           "Bantal",
			Price:          20000,
			Image:          fmt.Sprintf("https://%s.s3.%s.amazonaws.com/assets/%s", os.Getenv("AWS_BUCKET"), os.Getenv("AWS_REGION"), "no-image.png"),
			Quantity:       1,
			TotalPriceItem: 20000,
		})
		_ = tx.Commit()

		order = userOrder
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

	Describe("Create Comment /comments", func() {
		When("the fields are filled", func() {
			It("should return successful create comment response", func() {
				// Create Comment
				stringBody := fmt.Sprintf(`{"user_order_id": %v,"product_code": "Lfanp","description": "mantap bang","tag": "keren, mantap","rating": 4}`, order.IdOrder)
				requestBody := strings.NewReader(stringBody)
				request := httptest.NewRequest(http.MethodPost, "/api/comments", requestBody)
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

				Expect(responseBody["data"].(map[string]interface{})["product_code"]).To(Equal("Lfanp"))
				Expect(responseBody["data"].(map[string]interface{})["description"]).To(Equal("mantap bang"))
				Expect(responseBody["data"].(map[string]interface{})["tag"]).To(Equal("keren, mantap"))
				Expect(int(responseBody["data"].(map[string]interface{})["rating"].(float64))).To(Equal(4))
			})
		})
	})

	Describe("Find By Id Comment /comments/:id", func() {
		When("comment is not found", func() {
			It("should return error not found", func() {
				// Find By Id Comment
				request := httptest.NewRequest(http.MethodGet, "/api/comments/1", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["status"]).To(Equal("data not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("comment is found", func() {
			It("should return a successful find comment by id", func() {
				// Create Comment
				tx, _ := database.Begin()
				commentRepository := repository.NewCommentRepository()
				comment, _ := commentRepository.Create(context.Background(), tx, &entity.Comment{
					UserOrderId: order.IdOrder,
					ProductCode: "Lfanp",
					Description: sql.NullString{
						String: "mantap bang",
						Valid:  true,
					},
					Tag: sql.NullString{
						String: "keren, mantap",
						Valid:  true,
					},
					Rating: 4,
				})
				_ = tx.Commit()

				// Find By Id Comment
				request := httptest.NewRequest(http.MethodGet, "/api/comments/"+helper.IntToStr(comment.IdComment), nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"].(map[string]interface{})["product_code"]).To(Equal("Lfanp"))
				Expect(responseBody["data"].(map[string]interface{})["description"]).To(Equal("mantap bang"))
				Expect(responseBody["data"].(map[string]interface{})["tag"]).To(Equal("keren, mantap"))
				Expect(int(responseBody["data"].(map[string]interface{})["rating"].(float64))).To(Equal(4))
			})
		})
	})

	Describe("Find Comment By User Order Id /comments/user-order/:user_order_id", func() {
		When("comment is not found", func() {
			It("should return error not found", func() {
				// Find Comment By User Order Id
				request := httptest.NewRequest(http.MethodGet, "/api/comments/user-order/1", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["status"]).To(Equal("data not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("comment is found", func() {
			It("should return a successful find comment by user order id", func() {
				// Create Comment
				tx, _ := database.Begin()
				commentRepository := repository.NewCommentRepository()
				comment, _ := commentRepository.Create(context.Background(), tx, &entity.Comment{
					UserOrderId: order.IdOrder,
					ProductCode: "Lfanp",
					Description: sql.NullString{
						String: "mantap bang",
						Valid:  true,
					},
					Tag: sql.NullString{
						String: "keren, mantap",
						Valid:  true,
					},
					Rating: 4,
				})
				_ = tx.Commit()

				// Find Comment By User Order Id
				request := httptest.NewRequest(http.MethodGet, "/api/comments/user-order/"+helper.IntToStr(comment.UserOrderId), nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"].(map[string]interface{})["product_code"]).To(Equal("Lfanp"))
				Expect(responseBody["data"].(map[string]interface{})["description"]).To(Equal("mantap bang"))
				Expect(responseBody["data"].(map[string]interface{})["tag"]).To(Equal("keren, mantap"))
				Expect(int(responseBody["data"].(map[string]interface{})["rating"].(float64))).To(Equal(4))
			})
		})
	})

	Describe("Find All Rating By Product Code /comments/rating/:product_code", func() {
		When("rating's comment by product code is not found", func() {
			It("should return error not found", func() {
				// Find All Rating By Product Code
				request := httptest.NewRequest(http.MethodGet, "/api/comments/rating/XX1", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["status"]).To(Equal("data not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("rating's comment is exists", func() {
			It("should return a successful find comment by user order id", func() {
				// Create Comment
				tx, _ := database.Begin()
				commentRepository := repository.NewCommentRepository()
				comment1, _ := commentRepository.Create(context.Background(), tx, &entity.Comment{
					UserOrderId: order.IdOrder,
					ProductCode: "Lfanp",
					Description: sql.NullString{
						String: "mantap bang",
						Valid:  true,
					},
					Tag: sql.NullString{
						String: "keren, mantap",
						Valid:  true,
					},
					Rating: 4,
				})

				_, _ = commentRepository.Create(context.Background(), tx, &entity.Comment{
					UserOrderId: order.IdOrder,
					ProductCode: "Lfanp",
					Rating:      3,
				})

				_, _ = commentRepository.Create(context.Background(), tx, &entity.Comment{
					UserOrderId: order.IdOrder,
					ProductCode: "Lfanp",
					Rating:      4,
				})
				_ = tx.Commit()

				// Find All Rating By Product Code
				request := httptest.NewRequest(http.MethodGet, "/api/comments/rating/"+comment1.ProductCode, nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))

				commentResponse := responseBody["data"].([]interface{})
				Expect(int(commentResponse[0].(map[string]interface{})["rating"].(float64))).To(Equal(4))
				Expect(int(commentResponse[0].(map[string]interface{})["result_comment"].(float64))).To(Equal(1))
				Expect(int(commentResponse[0].(map[string]interface{})["result_rating"].(float64))).To(Equal(8))

				Expect(int(commentResponse[1].(map[string]interface{})["rating"].(float64))).To(Equal(3))
				Expect(int(commentResponse[1].(map[string]interface{})["result_comment"].(float64))).To(Equal(0))
				Expect(int(commentResponse[1].(map[string]interface{})["result_rating"].(float64))).To(Equal(3))
			})
		})
	})

	Describe("Find All Comment By Product Code /comments/product/:product_code", func() {
		When("comment by product code is not found", func() {
			It("should return error not found", func() {
				// Find All Comment By Product Code
				request := httptest.NewRequest(http.MethodGet, "/api/comments/product/XX1", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["status"]).To(Equal("data not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("comment is exists", func() {
			It("should return a successful find all comment by product code", func() {
				// Create Comment
				tx, _ := database.Begin()
				commentRepository := repository.NewCommentRepository()
				comment1, _ := commentRepository.Create(context.Background(), tx, &entity.Comment{
					UserOrderId: order.IdOrder,
					ProductCode: "Lfanp",
					Description: sql.NullString{
						String: "mantap bang",
						Valid:  true,
					},
					Tag: sql.NullString{
						String: "keren, mantap",
						Valid:  true,
					},
					Rating: 4,
				})

				comment2, _ := commentRepository.Create(context.Background(), tx, &entity.Comment{
					UserOrderId: order.IdOrder,
					ProductCode: "Lfanp",
					Tag: sql.NullString{
						String: "jelek, tidak memuaskan",
						Valid:  true,
					},
					Rating: 2,
				})
				tx.Commit()

				// Find All Comment By Product Code
				request := httptest.NewRequest(http.MethodGet, "/api/comments/product/"+comment1.ProductCode, nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))

				commentResponse := responseBody["data"].([]interface{})
				Expect(commentResponse[0].(map[string]interface{})["product_code"]).To(Equal(comment1.ProductCode))
				Expect(commentResponse[0].(map[string]interface{})["description"]).To(Equal(comment1.Description.String))
				Expect(commentResponse[0].(map[string]interface{})["tag"]).To(Equal(comment1.Tag.String))
				Expect(int(commentResponse[0].(map[string]interface{})["rating"].(float64))).To(Equal(comment1.Rating))

				Expect(commentResponse[1].(map[string]interface{})["product_code"]).To(Equal(comment2.ProductCode))
				Expect(commentResponse[1].(map[string]interface{})["description"]).To(BeNil())
				Expect(commentResponse[1].(map[string]interface{})["tag"]).To(Equal(comment2.Tag.String))
				Expect(int(commentResponse[1].(map[string]interface{})["rating"].(float64))).To(Equal(comment2.Rating))
			})
		})
	})
})
