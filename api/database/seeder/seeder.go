package seeder

import (
	"fmt"
	"github.com/arvians-id/apriori/database/faker"
	"github.com/arvians-id/apriori/service"
	"time"
)

func RegisterSeeder(service service.ProductService, totalSeeds int) {
	for i := 0; i < totalSeeds; i++ {
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
