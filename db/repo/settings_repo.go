package repo

import (
	"imagu/db"
)

func InitSettingsRepo() error {
	_, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS Settings (
        ID INTEGER PRIMARY KEY AUTOINCREMENT,
        AdminUser BOOLEAN NOT NULL
    );`)
	if err != nil {
		return err
	}
	err = EnsureDefaultSettings()
	if err != nil {
		return err
	}
	return nil
}
func EnsureDefaultSettings() error {
	var count int
	row := db.DB.QueryRow("SELECT COUNT(*) FROM Settings")
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		_, err := db.DB.Exec("INSERT INTO Settings (AdminUser) VALUES (false)")
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateAdminUser(adminUser bool) error {
	_, err := db.DB.Exec("UPDATE Settings SET AdminUser = ? WHERE ID = 1", adminUser)
	if err != nil {
		return err
	}
	return nil
}

func GetAdminUser() (bool, error) {
	var adminUser bool
	row := db.DB.QueryRow("SELECT AdminUser FROM Settings WHERE ID = 1")
	err := row.Scan(&adminUser)
	if err != nil {
		return false, err
	}
	return adminUser, nil
}
