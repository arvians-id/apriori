package faker

import (
	"apriori/helper"
	"apriori/model"
	"apriori/service"
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
	Category    string
	isEmpty     int
	mass        int
	Image       string
}

type ProductFaker struct {
	Product model.CreateProductRequest
}

func NewProductFaker() *ProductFaker {
	return &ProductFaker{}
}

func (p *ProductFaker) SetDefault() *ProductFaker {
	p.Product.Code = helper.RandomString(5)
	p.Product.Name = gofakeit.CarModel()
	p.Product.Price = int(gofakeit.Price(10000, 100000))
	return p
}

func (p *ProductFaker) SetDescription() *ProductFaker {
	p.Product.Description = gofakeit.Sentence(30)
	return p
}

func (p *ProductFaker) SetCategory() *ProductFaker {
	p.Product.Category = gofakeit.RandomString([]string{"Produk Bantal", "Produk Guling", "Produk Kasur", "Produk Elektronik", "Produk Kamar Mandi"})
	return p
}

func (p *ProductFaker) SetMass() *ProductFaker {
	p.Product.Mass = gofakeit.Number(100, 2000)
	return p
}

func (p *ProductFaker) SetImage(image string) *ProductFaker {
	if image == "" {
		p.Product.Image = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/assets/%s", os.Getenv("AWS_BUCKET"), os.Getenv("AWS_REGION"), "no-image.png")
	} else {
		p.Product.Image = image
	}
	return p
}

func (p *ProductFaker) Seed(service service.ProductService) *model.CreateProductRequest {
	_, _ = service.Create(context.Background(), &model.CreateProductRequest{
		Code:        p.Product.Code,
		Name:        p.Product.Name,
		Price:       p.Product.Price,
		Category:    p.Product.Category,
		Mass:        p.Product.Mass,
		Description: p.Product.Description,
		Image:       p.Product.Image,
	})

	return nil
}
