package main

import (
	"apriori/app/scheduler"
	"apriori/cmd"
	"apriori/config"
	"apriori/route"
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
