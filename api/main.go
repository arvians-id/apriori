package main

import (
	"github.com/arvians-id/apriori/cmd"
	"github.com/arvians-id/apriori/config"
	"github.com/arvians-id/apriori/route"
	"github.com/arvians-id/apriori/scheduler"
	"log"
)

func main() {
	// Command Line
	cmd.Execute()

	// Scheduler
	scheduler.Execute()

	// Server
	configuration := config.New()
	initialized, _ := route.NewInitializedServer(configuration)

	port := configuration.Get("PORT")
	err := initialized.Run(":" + port)
	if err != nil {
		log.Fatal("cannot run the server ", err)
		return
	}
}
