package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var dbPool *sql.DB

func dbPoolInit(sc *ServerConfig) (*sql.DB, error) {
	// open a db pool
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", sc.DBUsername, sc.DBPassword, sc.DBHostAddress, sc.DBName)
	db, dbErr := sql.Open("mysql", dbDSN)
	if dbErr != nil {
		return nil, dbErr
	}

	// setup db pool
	db.SetMaxOpenConns(0)
	db.SetMaxIdleConns(3)
	db.SetConnMaxLifetime(0)

	// check db connection
	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	return db, nil
}
