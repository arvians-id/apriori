package setup

import (
	"database/sql"
	"github.com/arvians-id/apriori/cmd/server"
	"github.com/arvians-id/apriori/config"
	"github.com/gin-gonic/gin"
)

func ModuleSetup(configuration config.Config) (*gin.Engine, *sql.DB) {
	initialized, db := server.NewInitializedServer(configuration)
	return initialized, db
}
