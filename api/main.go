package main

import (
	"apriori/config"
	"apriori/route"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func init() {
	config.SetupConfiguration()
}

func main() {
	configuration := config.New()
	initialized, _ := route.NewInitializedServer(configuration)

	// Start App
	addr := fmt.Sprintf(":%v", configuration.Get("APP_PORT"))
	err := initialized.Run(addr)
	if err != nil {
		log.Fatal("cannot run the server ", err)
		return
	}
}
