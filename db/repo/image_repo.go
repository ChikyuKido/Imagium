package repo

import (
	"imagu/db"
	"imagu/db/model"
)

func InitImageRepo() error {
	_, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS Images (
	    ID INTEGER PRIMARY KEY AUTOINCREMENT,
	    Name TEXT NOT NULL,
	    Author INTEGER NOT NULL,
	    UUID TEXT NOT NULL,
	    Size INTEGER NOT NULL
			);`)
	if err != nil {
		return err
	}
	return nil
}

func CreateImage(name string, author int, uuid string, size int64) error {
	_, err := db.DB.Exec(`INSERT INTO Images (Name, Author,UUID,Size) VALUES (?, ?,?,?)`, name, author, uuid, size)
	if err != nil {
		return err
	}
	return nil
}

func GetImageFromUUID(uuid string) (*model.ImageModel, error) {
	var image model.ImageModel

	query := `SELECT ID, Name, Author,UUID,Size FROM Images WHERE UUID = ?`

	err := db.DB.QueryRow(query, uuid).Scan(&image.ID, &image.Name, &image.Author, &image.UUID, &image.Size)
	if err != nil {
		return nil, err
	}
	return &image, nil
}
func UpdateSize(uuid string, diff int64) error {
	query := `UPDATE Images SET Size = Size + ? WHERE UUID = ?`
	_, err := db.DB.Exec(query, diff, uuid)
	if err != nil {
		return err
	}
	return nil
}
