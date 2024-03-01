package cmd

import (
	"github.com/johncave/podinate/cli/apiclient"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Connect to a Podinate instance",
	Long:  `Login to a Podinate instance and set it up as a profile in your config file. To set up a different profile, use the --profile flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiclient.StartLogin()
	},
}
