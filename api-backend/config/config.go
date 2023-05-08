package config

import (
	"database/sql"
	"fmt"
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

// Init - Initialize the service.
func Init() error {
	// TODO - read in the config file

	crdb := fmt.Sprintf("host=%s port=%d dbname=%s sslmode=verify-full sslrootcert=/cockroach/cockroach-certs/ca.crt sslkey=/cockroach/cockroach-certs/client.root.key sslcert=/cockroach/cockroach-certs/client.root.crt", "masterdb-public", 26257, "podinate")
	err := ConnectDatabase(crdb)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")
	return nil
}

// Cleanup - Cleanup any resources used by this service.
func Cleanup() {
	DB.Close()
	log.Println("Cleanup complete")
}
