package setup

import (
	"database/sql"
	"time"
)

func SuiteSetup() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/apriori_test?parseTime=true")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db, nil
}

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
