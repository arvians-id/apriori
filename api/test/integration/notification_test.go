package integration

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/arvians-id/apriori/cmd/config"
	"github.com/arvians-id/apriori/cmd/library/redis"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/arvians-id/apriori/internal/repository/postgres"
	"github.com/arvians-id/apriori/test/setup"
	"github.com/arvians-id/apriori/util"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
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
	var notification1 *model.Notification
	configuration := config.New("../../.env.test")

	BeforeEach(func() {
		// Setup Configuration
		router, db := setup.ModuleSetup(configuration)

		database = db
		server = router

		// User Authentication
		// Create user
		tx, _ := database.Begin()
		userRepository := postgres.NewUserRepository()
		password, _ := bcrypt.GenerateFromPassword([]byte("Rahasia123"), bcrypt.DefaultCost)
		user, _ := userRepository.Create(context.Background(), tx, &model.User{
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

		tx, _ = database.Begin()
		notificationRepository := postgres.NewNotificationRepository()
		description := "This is first notification"
		url := "https://google.com"
		notificationOne, _ := notificationRepository.Create(context.Background(), tx, &model.Notification{
			UserId:      user.IdUser,
			Title:       "First notification",
			Description: &description,
			URL:         &url,
			IsRead:      false,
			CreatedAt:   time.Now(),
		})

		description = "This is second notification"
		url = "https://facebook.com"
		_, _ = notificationRepository.Create(context.Background(), tx, &model.Notification{
			UserId:      user.IdUser,
			Title:       "Second notification",
			Description: &description,
			URL:         &url,
			IsRead:      false,
			CreatedAt:   time.Now(),
		})
		_ = tx.Commit()

		notification1 = notificationOne
	})

	AfterEach(func() {
		// Setup Configuration
		_, db := setup.ModuleSetup(configuration)
		defer db.Close()

		cacheService := redis.NewCacheService(configuration)
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
				Expect(userOrderResponse[1].(map[string]interface{})["title"]).To(Equal("First notification"))
				Expect(userOrderResponse[1].(map[string]interface{})["description"]).To(Equal("This is first notification"))
				Expect(userOrderResponse[1].(map[string]interface{})["url"]).To(Equal("https://google.com"))
				Expect(userOrderResponse[1].(map[string]interface{})["is_read"].(bool)).To(BeFalse())

				Expect(userOrderResponse[0].(map[string]interface{})["title"]).To(Equal("Second notification"))
				Expect(userOrderResponse[0].(map[string]interface{})["description"]).To(Equal("This is second notification"))
				Expect(userOrderResponse[0].(map[string]interface{})["url"]).To(Equal("https://facebook.com"))
				Expect(userOrderResponse[0].(map[string]interface{})["is_read"].(bool)).To(BeFalse())
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
				Expect(userOrderResponse[1].(map[string]interface{})["title"]).To(Equal("First notification"))
				Expect(userOrderResponse[1].(map[string]interface{})["description"]).To(Equal("This is first notification"))
				Expect(userOrderResponse[1].(map[string]interface{})["url"]).To(Equal("https://google.com"))
				Expect(userOrderResponse[1].(map[string]interface{})["is_read"].(bool)).To(BeFalse())

				Expect(userOrderResponse[0].(map[string]interface{})["title"]).To(Equal("Second notification"))
				Expect(userOrderResponse[0].(map[string]interface{})["description"]).To(Equal("This is second notification"))
				Expect(userOrderResponse[0].(map[string]interface{})["url"]).To(Equal("https://facebook.com"))
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
				Expect(userOrderResponse[1].(map[string]interface{})["title"]).To(Equal("First notification"))
				Expect(userOrderResponse[1].(map[string]interface{})["description"]).To(Equal("This is first notification"))
				Expect(userOrderResponse[1].(map[string]interface{})["url"]).To(Equal("https://google.com"))
				Expect(userOrderResponse[1].(map[string]interface{})["is_read"].(bool)).To(BeTrue())

				Expect(userOrderResponse[0].(map[string]interface{})["title"]).To(Equal("Second notification"))
				Expect(userOrderResponse[0].(map[string]interface{})["description"]).To(Equal("This is second notification"))
				Expect(userOrderResponse[0].(map[string]interface{})["url"]).To(Equal("https://facebook.com"))
				Expect(userOrderResponse[0].(map[string]interface{})["is_read"].(bool)).To(BeTrue())
			})
		})
	})

	Describe("Mark One Notification Id /notifications/mark/:id", func() {
		When("the notification is exists", func() {
			It("should return successful find all notifications response with different is read status", func() {
				// Mark One Notification By Id
				request := httptest.NewRequest(http.MethodPatch, "/api/notifications/mark/"+util.IntToStr(notification1.IdNotification), nil)
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
				Expect(userOrderResponse[1].(map[string]interface{})["title"]).To(Equal("First notification"))
				Expect(userOrderResponse[1].(map[string]interface{})["description"]).To(Equal("This is first notification"))
				Expect(userOrderResponse[1].(map[string]interface{})["url"]).To(Equal("https://google.com"))
				Expect(userOrderResponse[1].(map[string]interface{})["is_read"].(bool)).To(BeTrue())

				Expect(userOrderResponse[0].(map[string]interface{})["title"]).To(Equal("Second notification"))
				Expect(userOrderResponse[0].(map[string]interface{})["description"]).To(Equal("This is second notification"))
				Expect(userOrderResponse[0].(map[string]interface{})["url"]).To(Equal("https://facebook.com"))
				Expect(userOrderResponse[0].(map[string]interface{})["is_read"].(bool)).To(BeFalse())
			})
		})
	})
})
