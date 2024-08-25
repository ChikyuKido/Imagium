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
		_, err = db.DB.Exec(`INSERT INTO Users (Username, Password, Roles) VALUES (?, ?,?)`, "guest", "nopassword", "viewImage,uploadImage,viewStats")
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
	_, err = db.DB.Exec(`INSERT INTO Users (Username, Password,Roles) VALUES (?, ?,?)`, username, hashedPassword, "viewImage,uploadImage,viewStats")
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
func HasRole(user *model.User, roleToCheck string) bool {
	roles := strings.Split(user.Roles, ",")

	for _, role := range roles {
		if strings.TrimSpace(role) == roleToCheck {
			return true
		}
	}

	return false
}

func AddRole(user *model.User, roleToCheck string) error {
	if !validRole(roleToCheck) {
		return fmt.Errorf("Invalid role: %s", roleToCheck)
	}
	roles := strings.Split(user.Roles, ",")
	contains := false
	for _, role := range roles {
		if strings.TrimSpace(role) == roleToCheck {
			contains = true
		}
	}
	if !contains {
		roles = append(roles, roleToCheck)
	}
	user.Roles = strings.Join(roles, ",")
	err := updateRole(user)
	if err != nil {
		return err
	}
	return nil
}

func RemoveRole(user *model.User, roleToCheck string) error {
	if !validRole(roleToCheck) {
		return fmt.Errorf("Invalid role: %s", roleToCheck)
	}
	roles := strings.Split(user.Roles, ",")
	newRoles := make([]string, 0)
	for _, role := range roles {
		if strings.TrimSpace(role) != roleToCheck {
			newRoles = append(newRoles, role)
		}
	}
	user.Roles = strings.Join(newRoles, ",")
	err := updateRole(user)
	if err != nil {
		return err
	}
	return nil
}
func GetAllUsers() ([]model.User, error) {
	var users []model.User

	query := `SELECT ID, Username, Password, Roles FROM Users`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Roles)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
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
func updateRole(user *model.User) error {
	query := `UPDATE Users SET Roles = ? WHERE ID = ?`
	_, err := db.DB.Exec(query, user.Roles, user.ID)
	if err != nil {
		return err
	}
	return nil
}
