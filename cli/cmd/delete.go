package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/johncave/podinate/cli/package_parser"
	"github.com/johncave/podinate/cli/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	deleteCmd.AddCommand(deletePodCmd)
	deleteCmd.AddCommand(deleteProjectCmd)
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.PersistentFlags().Bool("y", false, "Delete without asking for confirmation")
}

var deleteCmd = &cobra.Command{
	Use:          "delete",
	Short:        "Delete things on Podinate",
	Long:         `Deletes things on Podinate`,
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg, err := package_parser.Parse(args[0])
		if err != nil {
			return err
		}
		return pkg.Delete()

	},
}

var deletePodCmd = &cobra.Command{
	Use:     "pod",
	Aliases: []string{"pods", "po"},
	Short:   "Delete pods on Podinate",
	Long:    `Deletes pods on Podinate`,

	Run: func(cmd *cobra.Command, args []string) {
		theProject, err := sdk.ProjectGetByID(viper.GetString("project"))
		if err != nil {
			fmt.Print("Could not find this project %s", viper.GetString("project"))
			os.Exit(1)
		}

		thePod, err := theProject.GetPodByID(args[0])

		//_, _, err := apiclient.C.PodApi.ProjectProjectIdPodPodIdGet(cmd.Context(), viper.GetString("project"), args[0]).Account(viper.GetString("account")).Execute()

		if err != nil {
			fmt.Println("Could not find this pod")
			os.Exit(1)
		}

		if !viper.GetBool("y") {
			fmt.Printf("Are you sure you want to delete pod %s? [y/N]: ", args[0])
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				fmt.Println("Aborting")
				return
			}
		}

		// Delete the pod
		//_, err = apiclient.C.PodApi.ProjectProjectIdPodPodIdDelete(cmd.Context(), viper.GetString("project"), args[0]).Account(viper.GetString("account")).Execute()

		err = thePod.Delete()
		if err != nil {
			fmt.Println("Could not delete pod", err.Error())
			os.Exit(1)
		}
	},
}

// Delete a project
var deleteProjectCmd = &cobra.Command{
	Use:     "project",
	Aliases: []string{"projects", "proj", "projs"},
	Short:   "Delete projects on Podinate",
	Long:    `Deletes projects on Podinate`,
	Run: func(cmd *cobra.Command, args []string) {
		p, err := sdk.ProjectGetByID(args[0])
		if err != nil {
			fmt.Println("Could not find this project")
			os.Exit(1)
		}

		if !viper.GetBool("y") {
			fmt.Printf("Are you sure you want to delete project %s? [y/N]: ", p.Name)
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				fmt.Println("Aborting")
				return
			}
		}

		// Delete the project
		err = p.Delete()
		if err != nil {
			fmt.Println("Could not delete project", err.Error())
			os.Exit(1)
		}

	},
}
