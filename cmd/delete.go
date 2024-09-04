package cmd

import (
	"github.com/podinate/podinate/engine"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:          "delete",
	Short:        "Delete a podinate package",
	Long:         `Delete all the Kubernetes objects from a PodFile or Kubernetes YAML file. For example "podinate delete my-app.pf" or "podinate delete my-pod.yaml".`,
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		// I haven't figured out how to apply *.pod yet
		// path := "*.pod"
		// if len(args) > 0 {
		// 	path = args[0]
		// }

		pkg, err := engine.Parse(args)
		if err != nil {
			return err
		}

		err = pkg.Apply(cmd.Context(), true)
		if err != nil {
			return err
		}
		return nil
	},
}
