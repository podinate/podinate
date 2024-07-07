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
	Long:         `Apply a podinate definition. Defaults to reading *.pod or package/*.pod files in the current directory.`,
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

		err = pkg.Apply(cmd.Context())
		if err != nil {
			return err
		}
		return nil
	},
}
