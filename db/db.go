package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

var DB *sql.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		logrus.Fatal("Could not open database: ", err)
	}
	if err = DB.Ping(); err != nil {
		logrus.Fatal("Could not open database: ", err)
	}
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
