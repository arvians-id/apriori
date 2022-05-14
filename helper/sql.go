package helper

import "database/sql"

func CommitOrRollback(tx *sql.Tx) {
	if err := recover(); err != nil {
		errRollback := tx.Rollback()
		PanicIfError(errRollback)
	} else {
		errCommit := tx.Commit()
		PanicIfError(errCommit)
	}
}
