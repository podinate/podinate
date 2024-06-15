package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// func init() {
// 	rootCmd.AddCommand(launchCmd)
// }

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch an app on Podinate",
	Long:  `Logs in then provides options to get an app running on Podinate`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Don't run this command, it's not ready yet")
		os.Exit(1)

	},
}
