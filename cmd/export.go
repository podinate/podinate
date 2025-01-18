package cmd

import (
	"fmt"

	"github.com/podinate/podinate/engine"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(exportCmd)
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export a podinate definition to Kubernetes YAML",
	Long:  `Export a PodFile to Kubernetes YAML. This can be useful as part of a CI/CD process or if you want to chain Podinate with other tools. \n For example "podinate export my-app.pf" or "podinate export my-pod.yaml".`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg, err := engine.Parse(args)
		if err != nil {
			return err
		}

		out, err := pkg.Export(cmd.Context())

		if err != nil {
			return err
		}
		fmt.Print(*out)
		return nil
	},
}
