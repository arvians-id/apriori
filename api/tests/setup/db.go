package setup

import (
	"apriori/config"
	"database/sql"
	"fmt"
	"net/url"
)

func SuiteSetupMySQL(configuration config.Config) (*sql.DB, error) {
	conf := url.Values{}
	conf.Add("parseTime", "True")

	username := configuration.Get("DB_USERNAME")
	password := configuration.Get("DB_PASSWORD")
	host := configuration.Get("DB_HOST")
	port := configuration.Get("DB_PORT")
	database := configuration.Get("DB_DATABASE")

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?%v", username, password, host, port, database, conf.Encode())
	db, err := sql.Open(configuration.Get("DB_CONNECTION"), dsn)
	if err != nil {
		return nil, err
	}

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
