package model

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Roles    string `json:"roles"`
}

const (
	ROLE_ADMIN        = "admin"
	ROLE_VIEW_IMAGE   = "viewImage"
	ROLE_UPLOAD_IMAGE = "uploadImage"
	ROLE_REGISTER     = "register"
	ROLE_VIEW_STATS   = "viewStats"
)

var Roles = []string{
	ROLE_ADMIN,
	ROLE_VIEW_IMAGE,
	ROLE_UPLOAD_IMAGE,
	ROLE_REGISTER,
	ROLE_VIEW_STATS,
}
