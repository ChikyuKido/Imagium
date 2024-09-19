package model

import "gorm.io/gorm"

type SettingsModel struct {
	gorm.Model
	AdminRegister             bool
	AggregationTime           int
	AutomaticallyDeletionTime int
}
