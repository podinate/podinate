package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(tfCmd)
}

var tfCmd = &cobra.Command{
	Use:                "tofu",
	Aliases:            []string{"tf", "opentofu", "tfu", "tof", "terraform"},
	Short:              "Use OpenTofu to manage your infrastructure",
	Long:               `Passes through to the OpenTofu CLI with your user authentication configured`,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		os.Setenv("PODINATE_API_KEY", viper.GetString("profiles.0.api_key"))

		path, err := exec.LookPath("tofu")
		if err != nil {
			path, err = exec.LookPath("terraform")
			if err != nil {
				fmt.Println("Couldn't find OpenTofu or Terraform, please install one of them")
				os.Exit(1)
			}
		}

		e, err := exec.Command(path, args...).Output()
		if err != nil {
			fmt.Print(string(e))
			fmt.Println("Error running OpenTofu:", err)
			os.Exit(1)
		}
		fmt.Println(string(e))

	},
}
