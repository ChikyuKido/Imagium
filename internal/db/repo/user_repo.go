package repo

import (
	"fmt"
	"imagu/internal/db"
	model2 "imagu/internal/db/model"
	"imagu/internal/util"
	"strings"
)

func InitUserRepo() error {
	if err := db.DB.AutoMigrate(&model2.User{}); err != nil {
		return err
	}
	if !DoesUserByNameExists("guest") {
		guest := model2.User{
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

	user := model2.User{
		Username: username,
		Password: hashedPassword,
		Roles:    "viewImage,uploadImage,viewStats,viewLibrary",
	}
	if err := db.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}
func GetUserByName(username string) (*model2.User, error) {
	var user model2.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func GetUserById(id string) (*model2.User, error) {
	var user model2.User
	if err := db.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func DoesUserByNameExists(username string) bool {
	var count int64
	if err := db.DB.Model(&model2.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}
func HasRole(user *model2.User, roleToCheck string) bool {
	roles := strings.Split(user.Roles, ",")
	for _, role := range roles {
		if strings.TrimSpace(role) == roleToCheck {
			return true
		}
	}
	return false
}
func GetAllUsers() ([]model2.User, error) {
	var users []model2.User
	if err := db.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func validRole(role string) bool {
	for _, r := range model2.Roles {
		if r == role {
			return true
		}
	}
	return false
}
func UpdateRole(user *model2.User) error {
	if err := db.DB.Model(&user).Update("Roles", user.Roles).Error; err != nil {
		return err
	}
	return nil
}

func GetAllImagesByUser(user *model2.User, offset int, limit int) ([]model2.ImageModel, error) {
	var images []model2.ImageModel
	if err := db.DB.Where("author = ?", user.ID).Limit(limit).Offset(offset).Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func GetImageCountByUser(user *model2.User) (int64, error) {
	var count int64
	if err := db.DB.Model(&model2.ImageModel{}).Where("author = ?", user.ID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
