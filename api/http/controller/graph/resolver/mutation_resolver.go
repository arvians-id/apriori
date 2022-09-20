package resolver

import (
	"context"
	"github.com/arvians-id/apriori/http/controller/graph/generated"
	"github.com/arvians-id/apriori/http/request"
	"github.com/arvians-id/apriori/model"
)

func (r *mutationResolver) CreateProduct(ctx context.Context, input model.NewProduct) (*model.Product, error) {
	product, err := r.ProductService.Create(ctx, &request.CreateProductRequest{
		Code:        input.Code,
		Name:        input.Name,
		Description: "",
		Price:       input.Price,
		Category:    input.Category,
		Mass:        input.Mass,
		Image:       "",
	})
	if err != nil {
		return nil, err
	}

	return &model.Product{
		IdProduct:   product.IdProduct,
		Code:        product.Code,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category,
		IsEmpty:     product.IsEmpty,
		Mass:        product.Mass,
		Image:       product.Image,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
