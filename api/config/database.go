package config

import (
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func NewMySQL(configuration Config) (*sql.DB, error) {
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

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Limit connection with db pooling
	setMaxIdleConns, err := strconv.Atoi(configuration.Get("DATABASE_POOL_MIN"))
	if err != nil {
		panic(err)
	}
	setMaxOpenConns, err := strconv.Atoi(configuration.Get("DATABASE_POOL_MAX"))
	if err != nil {
		panic(err)
	}
	setConnMaxIdleTime, err := strconv.Atoi(configuration.Get("DATABASE_MAX_IDLE_TIME_SECOND"))
	if err != nil {
		panic(err)
	}
	setConnMaxLifetime, err := strconv.Atoi(configuration.Get("DATABASE_MAX_LIFE_TIME_SECOND"))
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(setMaxIdleConns)                                    // minimal connection
	db.SetMaxOpenConns(setMaxOpenConns)                                    // maximal connection
	db.SetConnMaxLifetime(time.Duration(setConnMaxIdleTime) * time.Minute) // unused connections will be deleted
	db.SetConnMaxIdleTime(time.Duration(setConnMaxLifetime) * time.Minute) // connection that can be used

	return db, nil
}

func NewPostgreSQL(configuration Config) (*sql.DB, error) {
	username := configuration.Get("DB_USERNAME")
	password := configuration.Get("DB_PASSWORD")
	host := configuration.Get("DB_HOST")
	port := configuration.Get("DB_PORT")
	database := configuration.Get("DB_DATABASE")
	sslMode := configuration.Get("DB_SSL_MODE")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, username, password, database, sslMode)
	db, err := sql.Open(configuration.Get("DB_CONNECTION"), dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Limit connection with db pooling
	setMaxIdleConns, err := strconv.Atoi(configuration.Get("DATABASE_POOL_MIN"))
	if err != nil {
		panic(err)
	}
	setMaxOpenConns, err := strconv.Atoi(configuration.Get("DATABASE_POOL_MAX"))
	if err != nil {
		panic(err)
	}
	setConnMaxIdleTime, err := strconv.Atoi(configuration.Get("DATABASE_MAX_IDLE_TIME_SECOND"))
	if err != nil {
		panic(err)
	}
	setConnMaxLifetime, err := strconv.Atoi(configuration.Get("DATABASE_MAX_LIFE_TIME_SECOND"))
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(setMaxIdleConns)                                    // minimal connection
	db.SetMaxOpenConns(setMaxOpenConns)                                    // maximal connection
	db.SetConnMaxLifetime(time.Duration(setConnMaxIdleTime) * time.Minute) // unused connections will be deleted
	db.SetConnMaxIdleTime(time.Duration(setConnMaxLifetime) * time.Minute) // connection that can be used

	return db, nil
}
