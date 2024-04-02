// Package settings manages global settings that are stored in the database.
// See the config package for settings that are stored in the config file.
// An example of a setting stored in the database is the certificate authority for the cluster.
package settings

import "github.com/johncave/podinate/api-backend/config"

// func Get gets the value of a setting.
func Get(section string, key string) (string, error) {
	// Get the value from the database
	var value string
	err := config.DB.QueryRow("SELECT value FROM settings WHERE section = $1 AND key = $2", section, key).Scan(&value)
	return value, err
}

// func Set sets the value of a setting.
func Set(section string, key string, value string) error {
	// Set the value in the database
	_, err := config.DB.Exec("INSERT INTO settings (section, key, value) VALUES ($1, $2, $3) ON CONFLICT (section, key) DO UPDATE SET value = $3", section, key, value)
	return err
}

// func Delete deletes a setting.
func Delete(section string, key string) error {
	// Delete the value from the database
	_, err := config.DB.Exec("DELETE FROM settings WHERE section = $1 AND key = $2", section, key)
	return err
}

// func GetAll gets all settings in a section.
func GetAll(section string) (map[string]string, error) {
	// Get all the values from the database
	rows, err := config.DB.Query("SELECT key, value FROM settings WHERE section = $1", section)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		err = rows.Scan(&key, &value)
		if err != nil {
			return nil, err
		}
		settings[key] = value
	}

	return settings, nil
}

// func DeleteAll deletes all settings in a section.
func DeleteAll(section string) error {
	// Delete all the values from the database
	_, err := config.DB.Exec("DELETE FROM settings WHERE section = $1", section)
	return err
}
