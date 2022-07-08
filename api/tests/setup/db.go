package setup

import (
	"database/sql"
)

func TearDownTest(db *sql.DB) error {
	_, err := db.Exec(`TRUNCATE TABLE users;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`TRUNCATE TABLE password_resets;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`TRUNCATE TABLE apriories;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`TRUNCATE TABLE products;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`TRUNCATE TABLE transactions;`)
	if err != nil {
		return err
	}

	return nil
}
