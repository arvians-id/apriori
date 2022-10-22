package setup

import (
	"database/sql"
	"github.com/arvians-id/apriori/cmd/config"
	"github.com/arvians-id/apriori/cmd/server"
	"github.com/gin-gonic/gin"
)

func ModuleSetup(configuration config.Config) (*gin.Engine, *sql.DB) {
	initialized, db := server.NewInitializedServer(configuration)
	return initialized, db
}
