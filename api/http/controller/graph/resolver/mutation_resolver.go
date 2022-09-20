package resolver

import (
	"context"
	"fmt"
	"github.com/arvians-id/apriori/http/controller/graph/generated"
	"github.com/arvians-id/apriori/http/controller/rest/request"
	"github.com/arvians-id/apriori/model"
	"os"
)

func (r *mutationResolver) CreateProduct(ctx context.Context, input model.CreateProductRequest) (*model.Product, error) {
	fileName := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/assets/%s", os.Getenv("AWS_BUCKET"), os.Getenv("AWS_REGION"), "no-image.png")
	if input.Image.Filename != "" {
		initFileName := fmt.Sprintf("%v-%s", input.Image.Size, input.Image.Filename)
		fileNameGenerated, err := r.StorageService.UploadFileS3GraphQL(input.Image, initFileName)
		if err != nil {
			return nil, err
		}
		defer os.Remove(initFileName)
		fileName = fileNameGenerated
	}

	product, err := r.ProductService.Create(ctx, &request.CreateProductRequest{
		Code:        input.Code,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Category:    input.Category,
		Mass:        input.Mass,
		Image:       fileName,
	})
	if err != nil {
		return nil, err
	}

	return product, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
