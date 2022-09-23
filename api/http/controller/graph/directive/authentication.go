package directive

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/arvians-id/apriori/http/middleware"
	"github.com/arvians-id/apriori/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"strings"
)

func IsAuthenticated(ctx context.Context) (jwt.Claims, *gin.Context, error) {
	c, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, nil, err
	}

	authorizationHeader := c.GetHeader("Authorization")
	if !strings.Contains(authorizationHeader, "Bearer") {
		return nil, nil, errors.New("invalid authorization header")
	}

	tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

	jwtService := service.NewJwtService()
	token, err := jwtService.ValidateToken(tokenString)
	if err != nil {
		return nil, nil, err
	}

	if !token.Valid {
		return nil, nil, errors.New("invalid token")
	}

	claims := token.Claims

	return claims, c, nil
}

func HasRoles(ctx context.Context, obj interface{}, next graphql.Resolver, roles string) (interface{}, error) {
	token, c, err := IsAuthenticated(ctx)
	if err != nil {
		return nil, err
	}

	claims := token.(jwt.MapClaims)

	var rolesTemp string
	admin := 1
	user := 2
	role := int(claims["role"].(float64))
	if role == admin {
		rolesTemp = "admin"
	} else if role == user {
		rolesTemp = "user"
	}

	if rolesTemp != roles {
		return nil, errors.New("you are not authorized")
	}

	c.Set("id_user", claims["id_user"])
	c.Set("role", claims["role"])

	return next(ctx)
}
