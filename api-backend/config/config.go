package config

import (
	"database/sql"
	"log"
)

var (
	// Config - The configuration for this service.
	DB *sql.DB
)

// Connect to the database
func ConnectDatabase(connectionString string) error {
	// TODO - update Connect with the required logic for this service method.
	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}
	return nil
}

// Cleanup - Cleanup any resources used by this service.
func Cleanup() {
	DB.Close()
	log.Println("Cleanup complete")
}
