// Package cmd contains a handful of commands that can be executed
// Most of the time we will run the command "podinate" with no arguments
// This will start the server
// If we run "podinate init" this will initialise the configuration
// This will create and print out the credentials for the initial user
// If the configuration is already initialised, this command will fail
// But we might also add commands in the future to do other things
// like "podinate reset" to reset the configuration
package cmd

import (
	"log"
	"net/http"

	router "github.com/johncave/podinate/api-backend/router"
	"github.com/spf13/cobra"
)

// Initialize performs a first time setup for the entire service.
// This means creating the initial admin user, the CA and the user's keypair.
// Also the first account.
// func Initialize() error {
// 	var uuid string

// 	err := config.DB.QueryRow()

// 	return nil
// }

var rootCmd = &cobra.Command{
	Use:   "podinate",
	Short: "Podinate is a simple containerisation solution",
	Long: `Podinate is a simple containerisation solution. This is the controller component that runs inside your cluster.
	
	To just start the server, simply run the command with no arguments.
	
	To initialise the configuration, run "podinate init", this will create and print out the credentials for the inital user.
	If the configuration is already initialised, this command will fail.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("Server initialising...")
		theRouter := router.GetRouter()
		return http.ListenAndServe(":3000", theRouter)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
