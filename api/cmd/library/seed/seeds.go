package seed

import (
	"fmt"
	"github.com/arvians-id/apriori/cmd/config"
	"github.com/arvians-id/apriori/cmd/library/aws"
	"github.com/arvians-id/apriori/database/seeder"
	"github.com/arvians-id/apriori/internal/repository/postgres"
	"github.com/arvians-id/apriori/internal/service"
	"github.com/spf13/cobra"
)

// seedsCmd represents the seeds command
var seedsCmd = &cobra.Command{
	Use:   "seeds",
	Short: "Insert data into the database",
	Long:  `Use seeds to insert data into the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		getName, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println(err)
			return
		}
		getTotal, err := cmd.Flags().GetInt("count")
		if err != nil {
			fmt.Println(err)
			return
		}

		if getName != "" && getTotal > 0 {
			if getName == "product" {
				err := productSeeder(getTotal)
				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Println("seed successfully executed")
				return
			} else {
				fmt.Println("ERROR: total seed should at least 1 data or seed name is not found")
				fmt.Println("ERROR: see 'apriori seeds list'")
				return
			}
		}

		if len(args) > 0 {
			fmt.Println("list of seeds:", listSeeds)
			return
		}

		fmt.Println("ERROR: command not found")
		fmt.Println("ERROR: see 'apriori seeds -h'")
	},
}

var nameSeeder string
var totalSeeds int
var allSeeds []string
var listSeeds = []string{"product seed - product"}

func init() {
	rootCmd.AddCommand(seedsCmd)

	seedsCmd.PersistentFlags().StringVarP(&nameSeeder, "name", "n", "", "name of the seed")
	seedsCmd.PersistentFlags().IntVarP(&totalSeeds, "count", "c", 0, "total of the seed")
	seedsCmd.MarkFlagsRequiredTogether("name", "count")
	seedsCmd.PersistentFlags().StringSliceVar(&allSeeds, "list", listSeeds, "list of all seed")
}

func productSeeder(totalSeeds int) error {
	// Setup Configuration
	configuration := config.New()
	db, err := config.NewPostgreSQL(configuration)
	if err != nil {
		return err
	}

	productRepository := postgres.NewProductRepository()
	aprioriRepository := postgres.NewAprioriRepository()
	storageLibrary := aws.NewStorageS3(configuration)
	productService := service.NewProductService(&productRepository, &aprioriRepository, storageLibrary, db)

	seeder.RegisterSeeder(productService, totalSeeds)

	return nil
}
