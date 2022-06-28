package setup

import (
	"apriori/config"
	"apriori/route"
	"github.com/gin-gonic/gin"
)

func ModuleSetup(configuration config.Config) *gin.Engine {
	initialized := route.NewInitializedServer(configuration)
	return initialized
}
