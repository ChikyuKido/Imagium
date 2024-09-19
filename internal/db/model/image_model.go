package model

import "gorm.io/gorm"

type ImageModel struct {
	gorm.Model
	Name      string
	Author    uint
	UUID      string
	Size      int64
	SubImages int
}
