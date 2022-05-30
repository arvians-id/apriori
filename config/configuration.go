package config

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func SetupConfiguration() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}
