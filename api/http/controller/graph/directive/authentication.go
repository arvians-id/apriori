package directive

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/arvians-id/apriori/http/middleware"
	"github.com/arvians-id/apriori/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"os"
	"strings"
)

func IsAuthenticated(c *gin.Context) (jwt.Claims, error) {
	authorizationHeader := c.GetHeader("Authorization")
	if !strings.Contains(authorizationHeader, "Bearer") {
		return nil, errors.New("invalid authorization header")
	}

	tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

	jwtService := service.NewJwtService()
	token, err := jwtService.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims := token.Claims

	return claims, nil
}

func HasRoles(ctx context.Context, obj interface{}, next graphql.Resolver, roles string, useAPIKey *bool) (interface{}, error) {
	c, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if useAPIKey == nil || *useAPIKey == true {
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey != os.Getenv("X_API_KEY") {
			return nil, errors.New("you are not authorized to access this resource")
		}
	}

	token, err := IsAuthenticated(c)
	if err != nil {
		return nil, err
	}

	claims := token.(jwt.MapClaims)

	var rolesTemp = make(map[string]bool)
	admin := 1
	user := 2
	role := int(claims["role"].(float64))
	if role == admin {
		rolesTemp["admin"] = true
		rolesTemp["user"] = true
	} else if role == user {
		rolesTemp["user"] = true
		rolesTemp["admin"] = false
	}

	if !rolesTemp[roles] {
		return nil, errors.New("you don't have permission to access this resource")
	}

	c.Set("id_user", claims["id_user"])
	c.Set("role", claims["role"])

	return next(ctx)
}
