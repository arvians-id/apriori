package setup

import (
	"database/sql"
	"github.com/arvians-id/apriori/config"
	"github.com/arvians-id/apriori/route"
	"github.com/gin-gonic/gin"
)

func ModuleSetup(configuration config.Config) (*gin.Engine, *sql.DB) {
	initialized, db := route.NewInitializedServer(configuration)
	return initialized, db
}
