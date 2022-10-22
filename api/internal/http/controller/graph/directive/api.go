package directive

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/arvians-id/apriori/internal/http/middleware"
	"os"
)

func ApiKey(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	c, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	apiKey := c.GetHeader("X-API-KEY")
	if apiKey != os.Getenv("X_API_KEY") {
		return nil, errors.New("you are not authorized to access this resource")
	}

	return next(ctx)
}
