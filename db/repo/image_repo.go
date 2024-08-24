package repo

import (
	"imagu/db"
)

func InitImageRepo() error {
	_, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS Images (
	    ID INTEGER PRIMARY KEY AUTOINCREMENT,
	    Name TEXT NOT NULL,
	    Author INTEGER NOT NULL,
	    UUID 	 TEXT NOT NULL
			);`)
	if err != nil {
		return err
	}
	return nil
}

func CreateImage(name string, author int, uuid string) error {
	_, err := db.DB.Exec(`INSERT INTO Images (Name, Author,UUID) VALUES (?, ?,?)`, name, author, uuid)
	if err != nil {
		return err
	}
	return nil
}
