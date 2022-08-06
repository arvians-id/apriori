package seeder

import (
	"apriori/database/faker"
	"apriori/service"
	"time"
)

type Seeder struct {
	Seeder interface{}
}

func RegisterSeeder(service service.ProductService) {
	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond * 1)
		product := faker.NewProductFaker()
		product.SetDefault()
		product.SetDescription()
		product.Seed(service)
	}
}
