package config

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func NewDB() *sql.DB {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true", username, password, host, port, database)
	db, err := sql.Open(os.Getenv("DB_CONNECTION"), dsn)
	if err != nil {
		log.Fatal("cannot connected database", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("request Timeout ", err)
		return nil

	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	log.Info("Connected Database MySQL")

	return db
}
