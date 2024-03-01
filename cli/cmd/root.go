package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/johncave/podinate/cli/apiclient"
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

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.config/podinate/credentials.yaml)")
	viper.BindPFlag("config_file", rootCmd.PersistentFlags().Lookup("config"))
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().StringP("project", "p", "", "project to use")
	viper.BindPFlag("project", rootCmd.PersistentFlags().Lookup("project"))

	rootCmd.PersistentFlags().String("profile", "default", "Profile and credentials to use")
	viper.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile"))

	rootCmd.PersistentFlags().StringP("account", "a", "default", "account to use")
	viper.BindPFlag("account", rootCmd.PersistentFlags().Lookup("account"))

	cobra.OnInitialize(initConfig)

	//viper.BindPFlag("project", rootCmd.PersistentFlags().Lookup("project"))
	// viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	// viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
	// viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	// viper.SetDefault("license", "apache")
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
	// If error, that's fine, there may not be a config file

	// if err := viper.ReadInConfig(); err != nil {
	// 	fmt.Println("Can't read config:", err)
	// 	os.Exit(1)
	// }

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
		if err != nil {
			switch err.(type) {
			case viper.ConfigFileNotFoundError:
				apiclient.StartLogin()
			default:
				fmt.Println("Can't read credentials file:", err)
			}
		}
	}

	// I have no idea why this works.
	// If I print the string it's the word "login"
	// but inverting this condition has the opposite effect
	// Wtf is happening here?
	//log.Printf("called as %+v %T", loginCmd.CalledAs(), loginCmd.CalledAs())
	if loginCmd.CalledAs() == "" {
		apiclient.SetupUser()
	}

	if viper.GetBool("verbose") {
		log.Println("Verbose output enabled")
		fmt.Printf("Config: %+v\n", viper.AllSettings())
	}
	//fmt.Printf("Config: %+v\n", viper.AllSettings())
}

var rootCmd = &cobra.Command{
	Use:   "podinate",
	Short: "The podinate cli, easy to use hosting for containerised applications",
	Long:  "podinate is a command-line interface for managing your containers and applicatioions. \nIt provides an easy-to-use interface for deploying and managing containerised applications.",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		cobra.CheckErr(cmd.Help())
	},
}

// Execute is the entry point for the cli
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
