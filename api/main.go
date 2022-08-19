package main

import (
	"apriori/cmd"
	"apriori/config"
	"apriori/route"
	"log"
)

func main() {
	cmd.Execute()
	configuration := config.New()
	initialized, _ := route.NewInitializedServer(configuration)

	// Start App
	port := configuration.Get("PORT")
	err := initialized.Run(":" + port)
	if err != nil {
		log.Fatal("cannot run the server ", err)
		return
	}
}
