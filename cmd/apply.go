package cmd

import (
	"github.com/podinate/podinate/engine"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:          "apply",
	Short:        "Apply a podinate definition",
	Long:         `Apply a PodFile or Kubernetes YAML file. For example "podinate apply my-app.pf" or "podinate apply my-pod.yaml".`,
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

		return pkg.Apply(cmd.Context(), false)

	},
}
