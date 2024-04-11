package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var (
	Version = "Development"
	Commit  = "Development"
	Date    = time.Now().Format("2006-01-02 15:04:05")
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Podinate CLI",
	Long:  `All software has versions. This is Podinate CLI's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Podinate CLI " + Version + "\nCommit " + Commit + "\nBuilt " + Date)
	},
}
