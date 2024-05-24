package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(copyCmd)
}

var copyCmd = &cobra.Command{
	Use:          "copy",
	Short:        "Copy files or directories to or from a Pod",
	Long:         `Copy files or directories to or from a Pod`,
	Args:         cobra.MinimumNArgs(2),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Copy from %s to %s\n", args[0], args[1])
		return nil
	},
}
