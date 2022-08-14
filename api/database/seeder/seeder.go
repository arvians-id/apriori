package seeder

import (
	"apriori/database/faker"
	"apriori/service"
	"fmt"
	"time"
)

func RegisterSeeder(service service.ProductService) {
	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond * 1)
		str := fmt.Sprintf("https://source.unsplash.com/random/640x400?sig=%v", i)

		product := faker.NewProductFaker()
		product.SetDefault()
		product.SetDescription()
		product.SetCategory()
		product.SetImage(str)
		product.SetMass()
		product.Seed(service)
	}
}
