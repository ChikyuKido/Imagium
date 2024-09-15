package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string
	Password string
	Roles    string
}

const (
	ROLE_ADMIN        = "admin"
	ROLE_VIEW_IMAGE   = "viewImage"
	ROLE_UPLOAD_IMAGE = "uploadImage"
	ROLE_REGISTER     = "register"
	ROLE_VIEW_STATS   = "viewStats"
	ROLE_VIEW_LIBRARY = "viewLibrary"
)

var Roles = []string{
	ROLE_ADMIN,
	ROLE_VIEW_IMAGE,
	ROLE_UPLOAD_IMAGE,
	ROLE_REGISTER,
	ROLE_VIEW_STATS,
	ROLE_VIEW_LIBRARY,
}
