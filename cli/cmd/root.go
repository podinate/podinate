package cmd

import (
	"fmt"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
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
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("license", "apache")
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigFile("podinate.yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//
	viper.AddConfigPath(home + "/.config/podinate/")
	viper.SetConfigName("credentials")
	if err := viper.MergeInConfig(); err != nil {
		fmt.Println("Can't read credentials file:", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "podinate",
	Short: "The podinate cli, easy to use hosting for containerised applications",
	Long:  "podinate is a command-line interface for managing your containers and applicatioions. \nIt provides an easy-to-use interface for deploying and managing your applications in Docker containers.",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		log.Println("Hello World")
		fmt.Println(viper.Get("project"))
		fmt.Println(viper.Get("user.key"))
	},
}

// Execute is the entry point for the cli
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
