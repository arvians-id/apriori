package middleware

import (
	"apriori/model"
	"apriori/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

func AuthJwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			c.JSON(http.StatusBadRequest, model.WebResponse{
				Code:   http.StatusBadRequest,
				Status: "invalid token",
				Data:   nil,
			})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		jwtService := service.NewJwtService()
		token, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.WebResponse{
				Code:   http.StatusUnauthorized,
				Status: err.Error(),
				Data:   nil,
			})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, model.WebResponse{
				Code:   http.StatusUnauthorized,
				Status: err.Error(),
				Data:   nil,
			})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("id_user", claims["id_user"])

		c.Next()
	}
}
