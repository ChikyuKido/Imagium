package repo

import (
	"imagu/db"
	"imagu/db/model"
)

func InitSettingsRepo() error {
	if err := db.DB.AutoMigrate(&model.SettingsModel{}); err != nil {
		return err
	}

	err := EnsureDefaultSettings()
	if err != nil {
		return err
	}
	return nil
}
func EnsureDefaultSettings() error {
	var count int64

	if err := db.DB.Model(&model.SettingsModel{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		defaultSettings := model.SettingsModel{AdminRegister: false}
		if err := db.DB.Create(&defaultSettings).Error; err != nil {
			return err
		}
	}

	return nil
}

func UpdateAdminUser(adminUser bool) error {
	var settings model.SettingsModel
	if err := db.DB.First(&settings).Error; err != nil {
		return err
	}
	settings.AdminRegister = adminUser
	if err := db.DB.Save(&settings).Error; err != nil {
		return err
	}

	return nil
}

func GetAdminUser() (bool, error) {
	var settings model.SettingsModel
	if err := db.DB.First(&settings).Error; err != nil {
		return false, err
	}
	return settings.AdminRegister, nil
}
