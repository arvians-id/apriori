package resolver

import (
	"bytes"
	"context"
	"fmt"
	"github.com/arvians-id/apriori/http/controller/graph/generated"
	"github.com/arvians-id/apriori/http/controller/rest/request"
	"github.com/arvians-id/apriori/model"
	"io/ioutil"
	"os"
)

func (r *mutationResolver) CreateProduct(ctx context.Context, input model.CreateProductRequest) (*model.Product, error) {
	filePath := fmt.Sprintf("%s-%s", input.Code, input.Image.Filename)
	stream, err := ioutil.ReadAll(input.Image.File)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(filePath, stream, 0644)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer := make([]byte, input.Image.Size)
	_, _ = file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)

	fileName, err := r.StorageService.UploadFileS3GraphQL(fileBytes, input.Image.Filename)
	if err != nil {
		return nil, err
	}

	_ = os.Remove(filePath)

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
