package cmd

import (
	"fmt"

	"github.com/johncave/podinate/cli/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(execCmd)
}

var execCmd = &cobra.Command{
	Use:     "exec",
	Aliases: []string{"execute"},
	Short:   "Execute a command in a pod",
	Long: `Execute a command in a pod
	For example to execute a command in a pod:
	podinate exec <pod_id> --project <project_id>ls -l`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		podID := args[0]
		projectID := viper.GetString("project")
		theProject, err := sdk.GetProjectByID(projectID)
		if err != nil {
			return err
		}
		thePod, err := theProject.GetPodByID(podID)
		if err != nil {
			return err
		}
		command := args[1:]
		result, err := thePod.Exec(command)
		if err != nil {
			return err
		}
		fmt.Print(result)
		return nil
	},
}
