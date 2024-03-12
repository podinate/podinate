package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/johncave/podinate/cli/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	logsCmd.Flags().BoolP("follow", "f", false, "Follow the logs")
	logsCmd.Flags().IntP("tail", "t", 25, "Number of lines to show from the end of the logs")
	rootCmd.AddCommand(logsCmd)
}

var logsCmd = &cobra.Command{
	Use:     "logs",
	Aliases: []string{"log"},
	Short:   "Get logs for a pod",
	Long:    `Get the logs for a pod, use the -f flag to follow the logs`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		podID := args[0]
		projectID := viper.GetString("project")
		theProject, err := sdk.ProjectGetByID(projectID)
		cobra.CheckErr(err)
		thePod, err := theProject.GetPodByID(podID)
		cobra.CheckErr(err)
		//fmt.Printf("%+v\n", thePod)
		lines, err := cmd.Flags().GetInt("tail")
		cobra.CheckErr(err)
		follow, err := cmd.Flags().GetBool("follow")
		cobra.CheckErr(err)

		if !follow {
			logs, err := thePod.GetLogs(lines, follow)
			cobra.CheckErr(err)
			fmt.Print(logs)
		} else {
			buf, err := thePod.GetLogsBuffer(lines, follow)
			cobra.CheckErr(err)
			_, err = io.Copy(os.Stdout, buf)
			cobra.CheckErr(err)

		}
	},
}
