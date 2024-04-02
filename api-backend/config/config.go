package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pelletier/go-toml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	lh "github.com/johncave/podinate/api-backend/loghandler"
)

var (
	// Config - The configuration for this service.
	DB     *sql.DB
	C      *toml.Tree
	Client *kubernetes.Clientset
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

		lh.Log.Errorw("Error connecting to database", "error", err)
		return err
	}

	err = DB.Ping()
	if err != nil {
		lh.Log.Errorw("Error connecting to database", "error", err)
		return err
	}
	return nil
}

// Init - Initialize the service.
func init() {
	// Read in the config file
	confile, err := toml.LoadFile(os.Getenv("CONFIG_FILE"))
	if err != nil {
		lh.Log.Fatalw("Error loading config file", "error", err)
		//log.Fatal(err)
	}

	// Make it globally accessible
	C = confile

	crdb := fmt.Sprintf("host=%s port=%d dbname=%s sslmode=disable user=%s password=%s", confile.Get("database.host"), confile.Get("database.port"), confile.Get("database.database"), confile.Get("database.user"), os.Getenv("POSTGRES_PASSWORD"))
	err = ConnectDatabase(crdb)
	if err != nil {
		lh.Log.Panicw("Error connecting to database", "error", err)
		//log.Fatal(err)
		//os.Exit(2)
	}

	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		lh.Log.Panicw("error getting Kubernetes config", "error", err)
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		lh.Log.Panicw("error connecting to kubernetes cluster", "error", err)
	}
	Client = clientset

	// Create admin account if not exists in the database
	// Horrible hack replace asap
	// TODO - replace this with a proper migration
	var uuid string
	err = DB.QueryRow("SELECT uuid FROM \"user\" WHERE id = 'administrator' AND main_provider IS NULL LIMIT 1").Scan(&uuid)

	lh.Log.Infow("Checking for administrator account", "error", err)
	if err == sql.ErrNoRows {
		// Create the admin account

	} else {
		lh.Log.Infow("Administrator account already exists")

	}

	log.Println("Connected to database")
}

// Cleanup - Cleanup any resources used by this service.
func Cleanup() {
	DB.Close()
	log.Println("Cleanup complete")
}
