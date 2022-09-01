package scheduler

import (
	"apriori/config"
	"apriori/helper"
	repository "apriori/repository/postgres"
	"apriori/route"
	"context"
	"log"
)

type ClearTablePasswordResetsSchedule struct {
}

func (scheduler *ClearTablePasswordResetsSchedule) Run() {
	repo := repository.NewPasswordResetRepository()

	ctx := context.Background()
	configuration := config.New()
	db, err := route.NewInitializedDatabase(configuration)
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
