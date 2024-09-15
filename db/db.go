package db

import (
	_ "github.com/mattn/go-sqlite3"
	gorm_logrus "github.com/onrik/gorm-logrus"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = gorm.Open(sqlite.Open("data/database.db"), &gorm.Config{
		Logger: gorm_logrus.New(),
	})
	if err != nil {
		logrus.Fatal("Could not open database: ", err)
	}
}

func CloseDB() {
	if DB != nil {
		db, err := DB.DB()
		if err != nil {
			return
		}
		db.Close()
	}
}
