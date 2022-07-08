package setup

import (
	"apriori/config"
	"apriori/route"
	"database/sql"
	"github.com/gin-gonic/gin"
)

func ModuleSetup(configuration config.Config) (*gin.Engine, *sql.DB) {
	initialized, db := route.NewInitializedServer(configuration)
	return initialized, db
}
