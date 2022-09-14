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
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"
)

/*
	Error :
		- Find All User Order /user-order/user
*/
var _ = Describe("User Order API", func() {
	var server *gin.Engine
	var database *sql.DB
	var tokenJWT string
	var cookie *http.Cookie
	var order1 *entity.UserOrder
	var order2 *entity.UserOrder
	var payment *entity.Payment
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
			OrderId: sql.NullString{
				String: "QESXmTNzqowsqTNZYmAD",
				Valid:  true,
			},
		})

		// Create User Order
		userOrderRepository := repository.NewUserOrderRepository()
		userOrder1, _ := userOrderRepository.Create(context.Background(), tx, &entity.UserOrder{
			PayloadId:      payload.IdPayload,
			Code:           "aXksCj2",
			Name:           "Bantal Biasa",
			Price:          20000,
			Image:          fmt.Sprintf("https://%s.s3.%s.amazonaws.com/assets/%s", os.Getenv("AWS_BUCKET"), os.Getenv("AWS_REGION"), "no-image.png"),
			Quantity:       1,
			TotalPriceItem: 20000,
		})

		userOrder2, _ := userOrderRepository.Create(context.Background(), tx, &entity.UserOrder{
			PayloadId:      payload.IdPayload,
			Code:           "aXksCj2",
			Name:           "Guling",
			Price:          10000,
			Image:          fmt.Sprintf("https://%s.s3.%s.amazonaws.com/assets/%s", os.Getenv("AWS_BUCKET"), os.Getenv("AWS_REGION"), "no-image.png"),
			Quantity:       2,
			TotalPriceItem: 20000,
		})
		_ = tx.Commit()

		order1 = userOrder1
		order2 = userOrder2
		payment = payload
	})

	AfterEach(func() {
		// Setup Configuration
		_, db := setup.ModuleSetup(configuration)
		defer db.Close()

		cacheService := service.NewCacheService(configuration)
		_ = cacheService.FlushDB(context.Background())

		err := setup.TearDownTest(db)
		if err != nil {
			log.Fatal(err)
		}
	})

	Describe("Find All Payment On User Order /user-order", func() {
		When("user not logged in yet", func() {
			It("should return error unauthorized/invalid token", func() {
				// Find All Payment on User Order
				request := httptest.NewRequest(http.MethodGet, "/api/user-order", nil)
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

		When("the user order is exists", func() {
			It("should return successful find all payment on user order response", func() {
				// Find All Payment on User Order
				request := httptest.NewRequest(http.MethodGet, "/api/user-order", nil)
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

				userOrderResponse := responseBody["data"].([]interface{})
				Expect(userOrderResponse[0].(map[string]interface{})["order_id"]).To(Equal("QESXmTNzqowsqTNZYmAD"))
			})
		})
	})

	//Describe("Find All User Order /user-order/user", func() {
	//	When("the user order is exists", func() {
	//		It("should return successful find all user order response", func() {
	//			// Find All User Order
	//			request := httptest.NewRequest(http.MethodGet, "/api/user-order/user", nil)
	//			request.Header.Add("Content-Type", "application/json")
	//			request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
	//			request.AddCookie(cookie)
	//			request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))
	//
	//			writer := httptest.NewRecorder()
	//			server.ServeHTTP(writer, request)
	//
	//			var responseBody map[string]interface{}
	//			_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)
	//
	//			Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
	//			Expect(responseBody["status"]).To(Equal("OK"))
	//
	//			userOrderResponse := responseBody["data"].([]interface{})
	//
	//			Expect(userOrderResponse[0].(map[string]interface{})["code"]).To(Equal(order1.Code))
	//			Expect(userOrderResponse[0].(map[string]interface{})["name"]).To(Equal(order1.Name))
	//			Expect(int(userOrderResponse[0].(map[string]interface{})["price"].(float64))).To(Equal(order1.Price))
	//			Expect(userOrderResponse[0].(map[string]interface{})["image"]).To(Equal(order1.Image))
	//			Expect(int(userOrderResponse[0].(map[string]interface{})["quantity"].(float64))).To(Equal(order1.Quantity))
	//			Expect(int(userOrderResponse[0].(map[string]interface{})["total_price_item"].(float64))).To(Equal(order1.TotalPriceItem))
	//
	//			Expect(userOrderResponse[1].(map[string]interface{})["code"]).To(Equal(order2.Code))
	//			Expect(userOrderResponse[1].(map[string]interface{})["name"]).To(Equal(order2.Name))
	//			Expect(int(userOrderResponse[1].(map[string]interface{})["price"].(float64))).To(Equal(order2.Price))
	//			Expect(userOrderResponse[1].(map[string]interface{})["image"]).To(Equal(order2.Image))
	//			Expect(int(userOrderResponse[1].(map[string]interface{})["quantity"].(float64))).To(Equal(order2.Quantity))
	//			Expect(int(userOrderResponse[1].(map[string]interface{})["total_price_item"].(float64))).To(Equal(order2.TotalPriceItem))
	//		})
	//	})
	//})

	Describe("Find All User Order By Order Id /user-order/:order_id", func() {
		When("the user order is not found", func() {
			It("should return error not found response", func() {
				// Find All User Order
				request := httptest.NewRequest(http.MethodGet, "/api/user-order/asasdw", nil)
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

		When("the user order is exists", func() {
			It("should return successful find all user order by order id response", func() {
				// Find All User Order
				request := httptest.NewRequest(http.MethodGet, "/api/user-order/"+payment.OrderId.String, nil)
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

				userOrderResponse := responseBody["data"].([]interface{})

				Expect(userOrderResponse[0].(map[string]interface{})["code"]).To(Equal(order1.Code))
				Expect(userOrderResponse[0].(map[string]interface{})["name"]).To(Equal(order1.Name))
				Expect(int64(userOrderResponse[0].(map[string]interface{})["price"].(float64))).To(Equal(order1.Price))
				Expect(userOrderResponse[0].(map[string]interface{})["image"]).To(Equal(order1.Image))
				Expect(int(userOrderResponse[0].(map[string]interface{})["quantity"].(float64))).To(Equal(order1.Quantity))
				Expect(int64(userOrderResponse[0].(map[string]interface{})["total_price_item"].(float64))).To(Equal(order1.TotalPriceItem))

				Expect(userOrderResponse[1].(map[string]interface{})["code"]).To(Equal(order2.Code))
				Expect(userOrderResponse[1].(map[string]interface{})["name"]).To(Equal(order2.Name))
				Expect(int64(userOrderResponse[1].(map[string]interface{})["price"].(float64))).To(Equal(order2.Price))
				Expect(userOrderResponse[1].(map[string]interface{})["image"]).To(Equal(order2.Image))
				Expect(int(userOrderResponse[1].(map[string]interface{})["quantity"].(float64))).To(Equal(order2.Quantity))
				Expect(int64(userOrderResponse[1].(map[string]interface{})["total_price_item"].(float64))).To(Equal(order2.TotalPriceItem))
			})
		})
	})

	Describe("Find User Order By Id /user-order/single/:id", func() {
		When("the user order is not found", func() {
			It("should return error not found response", func() {
				// Find All User Order
				request := httptest.NewRequest(http.MethodGet, "/api/user-order/single/12121", nil)
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

		When("the user order is exists", func() {
			It("should return successful find user order by id response", func() {
				// Find All User Order
				request := httptest.NewRequest(http.MethodGet, "/api/user-order/single/"+helper.IntToStr(order1.IdOrder), nil)
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

				Expect(responseBody["data"].(map[string]interface{})["code"]).To(Equal(order1.Code))
				Expect(responseBody["data"].(map[string]interface{})["name"]).To(Equal(order1.Name))
				Expect(int64(responseBody["data"].(map[string]interface{})["price"].(float64))).To(Equal(order1.Price))
				Expect(responseBody["data"].(map[string]interface{})["image"]).To(Equal(order1.Image))
				Expect(int(responseBody["data"].(map[string]interface{})["quantity"].(float64))).To(Equal(order1.Quantity))
				Expect(int64(responseBody["data"].(map[string]interface{})["total_price_item"].(float64))).To(Equal(order1.TotalPriceItem))
			})
		})
	})
})
