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
	"strings"
	"time"
)

var _ = Describe("Notification API", func() {
	var server *gin.Engine
	var database *sql.DB
	var tokenJWT string
	var cookie *http.Cookie
	var notification1 *entity.Notification
	var notification2 *entity.Notification
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

		tx, _ = database.Begin()
		notificationRepository := repository.NewNotificationRepository()
		notificationOne, _ := notificationRepository.Create(context.Background(), tx, &entity.Notification{
			UserId: user.IdUser,
			Title:  "First notification",
			Description: sql.NullString{
				String: "This is first notification",
				Valid:  true,
			},
			URL: sql.NullString{
				String: "https://google.com",
				Valid:  true,
			},
			IsRead:    false,
			CreatedAt: time.Now(),
		})

		notificationTwo, _ := notificationRepository.Create(context.Background(), tx, &entity.Notification{
			UserId: user.IdUser,
			Title:  "Second notification",
			Description: sql.NullString{
				String: "This is second notification",
				Valid:  true,
			},
			URL: sql.NullString{
				String: "https://facebook.com",
				Valid:  true,
			},
			IsRead:    false,
			CreatedAt: time.Now(),
		})
		_ = tx.Commit()

		notification1 = notificationOne
		notification2 = notificationTwo
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

	Describe("Find All Notification /notifications", func() {
		When("the notification is exists", func() {
			It("should return successful find all notifications response", func() {
				// Find All Notification
				request := httptest.NewRequest(http.MethodGet, "/api/notifications", nil)
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
				Expect(userOrderResponse[1].(map[string]interface{})["title"]).To(Equal(notification1.Title))
				Expect(userOrderResponse[1].(map[string]interface{})["description"]).To(Equal(notification1.Description.String))
				Expect(userOrderResponse[1].(map[string]interface{})["url"]).To(Equal(notification1.URL.String))
				Expect(userOrderResponse[1].(map[string]interface{})["is_read"].(bool)).To(Equal(notification1.IsRead))

				Expect(userOrderResponse[0].(map[string]interface{})["title"]).To(Equal(notification2.Title))
				Expect(userOrderResponse[0].(map[string]interface{})["description"]).To(Equal(notification2.Description.String))
				Expect(userOrderResponse[0].(map[string]interface{})["url"]).To(Equal(notification2.URL.String))
				Expect(userOrderResponse[0].(map[string]interface{})["is_read"].(bool)).To(Equal(notification2.IsRead))
			})
		})
	})

	Describe("Find All Notification By User Id /notifications/user", func() {
		When("the notification is exists", func() {
			It("should return successful find all notifications response", func() {
				// Find All Notification By User Id
				request := httptest.NewRequest(http.MethodGet, "/api/notifications/user", nil)
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
				Expect(userOrderResponse[1].(map[string]interface{})["title"]).To(Equal(notification1.Title))
				Expect(userOrderResponse[1].(map[string]interface{})["description"]).To(Equal(notification1.Description.String))
				Expect(userOrderResponse[1].(map[string]interface{})["url"]).To(Equal(notification1.URL.String))
				Expect(userOrderResponse[1].(map[string]interface{})["is_read"].(bool)).To(BeFalse())

				Expect(userOrderResponse[0].(map[string]interface{})["title"]).To(Equal(notification2.Title))
				Expect(userOrderResponse[0].(map[string]interface{})["description"]).To(Equal(notification2.Description.String))
				Expect(userOrderResponse[0].(map[string]interface{})["url"]).To(Equal(notification2.URL.String))
				Expect(userOrderResponse[0].(map[string]interface{})["is_read"].(bool)).To(BeFalse())
			})
		})
	})

	Describe("Mark All Notification By User Id /notifications/mark", func() {
		When("the notification is exists", func() {
			It("should return successful find all notifications response", func() {
				// Mark All Notification By User Id
				request := httptest.NewRequest(http.MethodPatch, "/api/notifications/mark", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				// Find All Notification By User Id
				request = httptest.NewRequest(http.MethodGet, "/api/notifications/user", nil)
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

				userOrderResponse := responseBody["data"].([]interface{})
				Expect(userOrderResponse[1].(map[string]interface{})["title"]).To(Equal(notification1.Title))
				Expect(userOrderResponse[1].(map[string]interface{})["description"]).To(Equal(notification1.Description.String))
				Expect(userOrderResponse[1].(map[string]interface{})["url"]).To(Equal(notification1.URL.String))
				Expect(userOrderResponse[1].(map[string]interface{})["is_read"].(bool)).To(BeTrue())

				Expect(userOrderResponse[0].(map[string]interface{})["title"]).To(Equal(notification2.Title))
				Expect(userOrderResponse[0].(map[string]interface{})["description"]).To(Equal(notification2.Description.String))
				Expect(userOrderResponse[0].(map[string]interface{})["url"]).To(Equal(notification2.URL.String))
				Expect(userOrderResponse[0].(map[string]interface{})["is_read"].(bool)).To(BeTrue())
			})
		})
	})

	Describe("Mark One Notification Id /notifications/mark/:id", func() {
		When("the notification is exists", func() {
			It("should return successful find all notifications response", func() {
				// Mark One Notification By Id
				request := httptest.NewRequest(http.MethodPatch, "/api/notifications/mark/"+helper.IntToStr(notification1.IdNotification), nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.Get("X_API_KEY"))
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				// Find All Notification By User Id
				request = httptest.NewRequest(http.MethodGet, "/api/notifications/user", nil)
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

				userOrderResponse := responseBody["data"].([]interface{})
				Expect(userOrderResponse[1].(map[string]interface{})["title"]).To(Equal(notification1.Title))
				Expect(userOrderResponse[1].(map[string]interface{})["description"]).To(Equal(notification1.Description.String))
				Expect(userOrderResponse[1].(map[string]interface{})["url"]).To(Equal(notification1.URL.String))
				Expect(userOrderResponse[1].(map[string]interface{})["is_read"].(bool)).To(BeTrue())

				Expect(userOrderResponse[0].(map[string]interface{})["title"]).To(Equal(notification2.Title))
				Expect(userOrderResponse[0].(map[string]interface{})["description"]).To(Equal(notification2.Description.String))
				Expect(userOrderResponse[0].(map[string]interface{})["url"]).To(Equal(notification2.URL.String))
				Expect(userOrderResponse[0].(map[string]interface{})["is_read"].(bool)).To(BeFalse())
			})
		})
	})
})
