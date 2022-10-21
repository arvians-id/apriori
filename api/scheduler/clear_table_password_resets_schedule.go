package scheduler

import (
	"context"
	"github.com/arvians-id/apriori/cmd/server"
	"github.com/arvians-id/apriori/config"
	"github.com/arvians-id/apriori/helper"
	repository "github.com/arvians-id/apriori/repository/postgres"
	"log"
)

type ClearTablePasswordResetsSchedule struct {
}

func (scheduler *ClearTablePasswordResetsSchedule) Run() {
	repo := repository.NewPasswordResetRepository()

	ctx := context.Background()
	configuration := config.New()
	db, err := server.NewInitializedDatabase(configuration)
	if err != nil {
		log.Fatalln(err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer helper.CommitOrRollback(tx)

	err = repo.Truncate(ctx, tx)
	if err != nil {
		log.Fatalln(err)
	}
}
