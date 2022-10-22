package util

import (
	"database/sql"
	"log"
)

func CommitOrRollback(tx *sql.Tx) {
	if tx != nil {
		if err := recover(); err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				log.Println("ERROR Rollback:", errRollback)
			}
		} else {
			errCommit := tx.Commit()
			if errCommit != nil {
				log.Println("ERROR Commit:", errCommit)
			}
		}
	}
}
