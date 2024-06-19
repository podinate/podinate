package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/cnrancher/autok3s/cmd"
	"github.com/podinate/podinate/cmd/cluster"
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

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.config/podinate/credentials.yaml)")
	viper.BindPFlag("config_file", rootCmd.PersistentFlags().Lookup("config"))
	//rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	//viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().String("project", "", "project to use")
	viper.BindPFlag("project", rootCmd.PersistentFlags().Lookup("project"))

	rootCmd.PersistentFlags().String("profile", "default", "Profile and credentials to use")
	viper.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile"))

	rootCmd.PersistentFlags().StringP("account", "a", "default", "account to use")
	viper.BindPFlag("account", rootCmd.PersistentFlags().Lookup("account"))

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug logging - LOTS of logs")

	// Add commnands related to autok3s
	kubectl := cmd.KubectlCommand()
	kubectl.Aliases = []string{"kubectl", "k"}
	rootCmd.AddCommand(kubectl, cluster.NewCommand())

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Bind the debug flag
		viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
		fmt.Println("Debug logging enabled", viper.GetBool("debug"))
		if viper.GetBool("debug") {
			logrus.SetLevel(logrus.TraceLevel)
			fmt.Println("Debug logging enabled")
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
