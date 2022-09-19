package resolver

import (
	"context"
	"github.com/arvians-id/apriori/helper"
	"github.com/arvians-id/apriori/http/controller/graph/generated"
	"github.com/arvians-id/apriori/http/controller/graph/model"
	"github.com/arvians-id/apriori/http/middleware"
)

func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct {
	*Resolver
}

func (r *queryResolver) Products(ctx context.Context) ([]*model.Product, error) {
	c, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	products, err := r.ProductService.FindAllByAdmin(c.Request.Context())
	if err != nil {
		return nil, err
	}

	var result []*model.Product
	for _, product := range products {
		result = append(result, &model.Product{
			IDProduct:   helper.IntToStr(product.IdProduct),
			Code:        product.Code,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Category:    product.Category,
			IsEmpty:     product.IsEmpty,
			Mass:        product.Mass,
			Image:       product.Image,
			CreatedAt:   product.CreatedAt.String(),
			UpdatedAt:   product.UpdatedAt.String(),
		})
	}

	return result, nil
}
