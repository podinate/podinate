package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pelletier/go-toml"
)

var (
	// Config - The configuration for this service.
	DB *sql.DB
	C  *toml.Tree
)

type configFile struct {
	DB DBConfig `toml:"database"`
}

type DBConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Database string `toml:"database"`
	User     string `toml:"user"`
}

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
	// Read in the config file
	confile, err := toml.LoadFile(os.Getenv("CONFIG_FILE"))
	if err != nil {
		log.Fatal(err)
	}

	// Make it globally accessible
	C = confile

	crdb := fmt.Sprintf("host=%s port=%d dbname=%s sslmode=disable user=%s password=%s", confile.Get("database.host"), confile.Get("database.port"), confile.Get("database.database"), confile.Get("database.user"), os.Getenv("POSTGRES_PASSWORD"))
	err = ConnectDatabase(crdb)
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
