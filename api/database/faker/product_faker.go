package faker

import (
	"apriori/model"
	"apriori/service"
	"apriori/utils"
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"os"
)

type Product struct {
	IdProduct   uint64
	Code        string
	Name        string
	Description string
	Price       int
	Image       string
}

type ProductFaker struct {
	Product model.CreateProductRequest
}

func NewProductFaker() *ProductFaker {
	return &ProductFaker{}
}

func (p *ProductFaker) SetDefault() *ProductFaker {
	p.Product.Code = utils.RandomString(5)
	p.Product.Name = gofakeit.CarModel()
	p.Product.Price = int(gofakeit.Price(10000, 100000))
	p.Product.Image = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/assets/%s", os.Getenv("AWS_BUCKET"), os.Getenv("AWS_REGION"), "no-image.png")
	return p
}

func (p *ProductFaker) SetDescription() *ProductFaker {
	p.Product.Description = gofakeit.Sentence(30)
	return p
}

func (p *ProductFaker) Seed(service service.ProductService) *model.CreateProductRequest {
	_, _ = service.Create(context.Background(), model.CreateProductRequest{
		Code:        p.Product.Code,
		Name:        p.Product.Name,
		Price:       p.Product.Price,
		Description: p.Product.Description,
		Image:       p.Product.Image,
	})

	return nil
}
