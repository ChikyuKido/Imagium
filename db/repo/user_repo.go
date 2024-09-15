package repo

import (
	"fmt"
	"imagu/db"
	"imagu/db/model"
	"imagu/util"
	"strings"
)

func InitUserRepo() error {
	if err := db.DB.AutoMigrate(&model.User{}); err != nil {
		return err
	}
	if !DoesUserByNameExists("guest") {
		guest := model.User{
			Username: "guest",
			Password: "nopassword",
			Roles:    "viewImage,uploadImage,viewStats",
		}
		if err := db.DB.Create(&guest).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateUser(username string, password string) error {
	if username == "guest" {
		return fmt.Errorf("guest user is preserved")
	}
	if DoesUserByNameExists(username) {
		return fmt.Errorf("user with that name already exists")
	}
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return err
	}

	user := model.User{
		Username: username,
		Password: hashedPassword,
		Roles:    "viewImage,uploadImage,viewStats,viewLibrary",
	}
	if err := db.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}
func GetUserByName(username string) (*model.User, error) {
	var user model.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func GetUserById(id string) (*model.User, error) {
	var user model.User
	if err := db.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func DoesUserByNameExists(username string) bool {
	var count int64
	if err := db.DB.Model(&model.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}
func HasRole(user *model.User, roleToCheck string) bool {
	roles := strings.Split(user.Roles, ",")
	for _, role := range roles {
		if strings.TrimSpace(role) == roleToCheck {
			return true
		}
	}
	return false
}
func GetAllUsers() ([]model.User, error) {
	var users []model.User
	if err := db.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func validRole(role string) bool {
	for _, r := range model.Roles {
		if r == role {
			return true
		}
	}
	return false
}
func UpdateRole(user *model.User) error {
	if err := db.DB.Model(&user).Update("Roles", user.Roles).Error; err != nil {
		return err
	}
	return nil
}

func GetAllImagesByUser(user *model.User, offset int, limit int) ([]model.ImageModel, error) {
	var images []model.ImageModel
	if err := db.DB.Where("author = ?", user.ID).Limit(limit).Offset(offset).Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func GetImageCountByUser(user *model.User) (int64, error) {
	var count int64
	if err := db.DB.Model(&model.ImageModel{}).Where("author = ?", user.ID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
