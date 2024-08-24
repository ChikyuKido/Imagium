package repo

import (
	"fmt"
	"imagu/db"
	"imagu/db/model"
	"imagu/util"
	"strings"
)

func InitUserRepo() error {
	_, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS Users (
	    ID INTEGER PRIMARY KEY AUTOINCREMENT,
	    Username TEXT NOT NULL,
	    Password TEXT NOT NULL,
	    Roles 	 TEXT NOT NULL
			);`)
	if !DoesUserByNameExists("guest") {
		_, err = db.DB.Exec(`INSERT INTO Users (Username, Password, Roles) VALUES (?, ?,?)`, "guest", "nopassword", "viewImage,uploadImage")
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
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
		return nil
	}
	_, err = db.DB.Exec(`INSERT INTO Users (Username, Password,Roles) VALUES (?, ?,?)`, username, hashedPassword, "viewImage,uploadImage")
	if err != nil {
		return err
	}
	return nil
}
func GetUserByName(username string) (*model.User, error) {
	var user model.User

	query := `SELECT ID, Username, Password,Roles FROM Users WHERE Username = ?`

	err := db.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Roles)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func DoesUserByNameExists(username string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM Users WHERE Username = ?)`
	err := db.DB.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}
func HasPermission(user *model.User, roleToCheck string) bool {
	roles := strings.Split(user.Roles, ",")

	for _, role := range roles {
		if strings.TrimSpace(role) == roleToCheck {
			return true
		}
	}

	return false
}
