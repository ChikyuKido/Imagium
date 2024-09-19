package repo

import (
	"imagu/internal/db"
	"imagu/internal/db/model"
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
		defaultSettings := model.SettingsModel{
			AdminRegister:             false,
			AggregationTime:           15,
			AutomaticallyDeletionTime: 1440 * 7,
		}
		if err := db.DB.Create(&defaultSettings).Error; err != nil {
			return err
		}
	}

	return nil
}
func UpdateSettings(settings model.SettingsModel) error {
	if err := db.DB.Model(&model.SettingsModel{}).Where("id = ?", 1).Updates(settings).Error; err != nil {
		return err
	}
	return nil
}
func UpdateAdminUser(adminUser bool) error {
	if err := db.DB.Model(&model.SettingsModel{}).Where("id = ?", 1).Updates(model.SettingsModel{AdminRegister: adminUser}).Error; err != nil {
		return err
	}
	return nil
}
func GetSettings() (model.SettingsModel, error) {
	var settings model.SettingsModel
	if err := db.DB.First(&settings).Error; err != nil {
		return settings, err
	}
	return settings, nil
}
func GetAdminUser() (bool, error) {
	var settings, err = GetSettings()
	if err != nil {
		return false, err
	}
	return settings.AdminRegister, nil
}
func GetAggregationTime() (int, error) {
	var settings, err = GetSettings()
	if err != nil {
		return -1, err
	}
	return settings.AggregationTime, nil
}
func GetAutomaticallyDeletionTime() (int, error) {
	var settings, err = GetSettings()
	if err != nil {
		return -1, err
	}
	return settings.AutomaticallyDeletionTime, nil
}
