package utils

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
)

func CommitOrRollback(tx *sql.Tx) {
	if err := recover(); err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			log.Fatal("cannot rollback ", err)
			return
		}
	} else {
		errCommit := tx.Commit()
		if errCommit != nil {
			log.Fatal("cannot commit ", err)
			return
		}
	}
}
