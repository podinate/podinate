package cmd

import (
	"fmt"

	sdk "github.com/johncave/podinate/cli/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	describeCmd.AddCommand(describePodCmd)
	describeCmd.AddCommand(describePackageCmd)
	rootCmd.AddCommand(describeCmd)
}

var describeCmd = &cobra.Command{
	Use:     "describe",
	Aliases: []string{"desc", "inspect"},
	Short:   "Describe a thing on Podinate",
	Long: `Describes a thing on Podinate
	Examples include:
	- describe pod my-pod
	- describe project my-project`,
	// Make the "pod" optional, ie "podinate describe abc123" should be the same as "podinate describe pod abc123"
	RunE: describePodCmd.RunE,
}

var describePodCmd = &cobra.Command{
	Use:   "pod",
	Short: "Describe a pod on Podinate",
	RunE: func(cmd *cobra.Command, args []string) error {
		project, sdkerr := sdk.GetProjectByID(viper.GetString("project"))
		if sdkerr != nil {
			return sdkerr
		}

		pod, sdker := project.GetPodByID(args[0])
		if sdker != nil {
			return sdker
		}

		// Print the pod
		fmt.Println(pod.Describe())
		return nil
	},
}

var describePackageCmd = &cobra.Command{
	Use:   "package",
	Short: "Describe a package on Podinate",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
