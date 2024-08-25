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
	    Size INTEGER NOT NULL,
	    SubImages INTEGER NOT NULL
			);`)
	if err != nil {
		return err
	}
	return nil
}

func CreateImage(name string, author int, uuid string, size int64) error {
	_, err := db.DB.Exec(`INSERT INTO Images (Name, Author,UUID,Size,SubImages) VALUES (?, ?,?,?,0)`, name, author, uuid, size)
	if err != nil {
		return err
	}
	return nil
}

func GetImageFromUUID(uuid string) (*model.ImageModel, error) {
	var image model.ImageModel

	query := `SELECT ID, Name, Author,UUID,Size,SubImages FROM Images WHERE UUID = ?`

	err := db.DB.QueryRow(query, uuid).Scan(&image.ID, &image.Name, &image.Author, &image.UUID, &image.Size, &image.SubImages)
	if err != nil {
		return nil, err
	}
	return &image, nil
}
func UpdateSizeAndCount(uuid string, sizeDiff int64, countDiff int) error {
	query := `UPDATE Images SET Size = Size + ?,SubImages = SubImages + ? WHERE UUID = ?`
	_, err := db.DB.Exec(query, sizeDiff, countDiff, uuid)
	if err != nil {
		return err
	}
	return nil
}
func GetAllImages() ([]model.ImageModel, error) {
	var images []model.ImageModel
	query := `SELECT ID, Name, Author,UUID,Size,SubImages FROM Images`

	rows, err := db.DB.Query(query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var image model.ImageModel
		err := rows.Scan(&image.ID, &image.Name, &image.Author, &image.UUID, &image.Size, &image.SubImages)
		if err != nil {
			panic(err)
		}
		images = append(images, image)
	}
	return images, nil
}
