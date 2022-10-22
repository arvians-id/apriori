package main

import (
	"github.com/arvians-id/apriori/cmd/config"
	"github.com/arvians-id/apriori/cmd/library/scheduler"
	"github.com/arvians-id/apriori/cmd/library/seed"
	"github.com/arvians-id/apriori/cmd/server"
	"log"
)

func main() {
	// Command Line
	seed.Execute()

	// Scheduler
	scheduler.Execute()

	// Server
	configuration := config.New()
	initialized, _ := server.NewInitializedServer(configuration)

	port := configuration.Get("PORT")
	err := initialized.Run(":" + port)
	if err != nil {
		log.Fatal("cannot run the server ", err)
		return
	}
}
