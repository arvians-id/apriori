package middleware

import (
	"apriori/model"
	"apriori/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
)

func AuthJwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		value, err := c.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				c.JSON(http.StatusUnauthorized, model.WebResponse{
					Code:   http.StatusUnauthorized,
					Status: err.Error(),
					Data:   nil,
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusBadRequest, model.WebResponse{
				Code:   http.StatusBadRequest,
				Status: err.Error(),
				Data:   nil,
			})
			c.Abort()
			return
		}

		jwtService := service.NewJwtService()
		token, err := jwtService.ValidateToken(value)
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
