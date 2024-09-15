package repo

import (
	"imagu/db"
	"imagu/db/model"
)

func InitImageRepo() error {
	if err := db.DB.AutoMigrate(&model.ImageModel{}); err != nil {
		return err
	}
	return nil
}

func CreateImage(name string, author uint, uuid string, size int64) error {
	image := model.ImageModel{
		Name:      name,
		Author:    author,
		UUID:      uuid,
		Size:      size,
		SubImages: 0,
	}

	if err := db.DB.Create(&image).Error; err != nil {
		return err
	}
	return nil
}
func GetImageFromUUID(uuid string) (*model.ImageModel, error) {
	var image model.ImageModel
	if err := db.DB.Where("uuid = ?", uuid).First(&image).Error; err != nil {
		return nil, err
	}
	return &image, nil
}
func UpdateSizeAndCount(uuid string, sizeDiff int64, countDiff int) error {
	var image model.ImageModel

	if err := db.DB.Where("uuid = ?", uuid).First(&image).Error; err != nil {
		return err
	}
	if err := db.DB.Model(&image).Updates(map[string]interface{}{
		"Size":      image.Size + sizeDiff,
		"SubImages": image.SubImages + countDiff,
	}).Error; err != nil {
		return err
	}

	return nil
}

func GetAllImages() ([]model.ImageModel, error) {
	var images []model.ImageModel
	if err := db.DB.Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}
