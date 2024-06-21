package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string
	verbose bool
)

// Config file stuff
func init() {

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug logging - LOTS of logs")
	rootCmd.PersistentFlags().String("kubeconfig", "", "Path to the kubeconfig file to use for CLI requests")

	// Add commnands related to autok3s
	// kubectl := cmd.KubectlCommand()
	// kubectl.Aliases = []string{"kubectl", "k"}
	// rootCmd.AddCommand(kubectl, cluster.NewCommand())

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {

		// Bind the kubeconfig flag
		viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))

		// Bind the debug flag
		viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
		//fmt.Println("Debug logging enabled", viper.GetBool("debug"))
		if viper.GetBool("debug") {
			logrus.SetLevel(logrus.TraceLevel)
			//fmt.Println("Debug logging enabled")
			logrus.Debug("Debug logging enabled")
		}

		return nil
	}
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigFile("podinate.yaml")
	}

	viper.ReadInConfig()

	if viper.GetBool("verbose") {
		log.Println("Verbose output enabled")
		fmt.Printf("Config: %+v\n", viper.AllSettings())
		//zap.NewAtomicLevel().SetLevel(zap.DebugLevel)
	}
	//fmt.Printf("Config: %+v\n", viper.AllSettings())
}

var rootCmd = &cobra.Command{
	Use:          "podinate",
	Short:        "The podinate cli, easy to use hosting for containerised applications",
	Long:         "podinate is a command-line interface for managing your containers and applicatioions. \nIt provides an easy-to-use interface for deploying and managing containerised applications.",
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here

		cobra.CheckErr(cmd.Help())
	},
}

// Execute is the entry point for the cli
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
