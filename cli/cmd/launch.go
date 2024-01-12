package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/johncave/podinate/cli/tui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(launchCmd)
}

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch an app on Podinate",
	Long:  `Logs in then provides options to get an app running on Podinate`,
	Run: func(cmd *cobra.Command, args []string) {

		//fmt.Printf("Response: %+v, r: %+v\n", resp, r)
		//fmt.Printf("URL: %s, Token: %s\n", *resp.Url, *resp.Token)

		// go if _, err := tea.NewProgram(spinner.New("Awaiting login, check your browser...", "points")).Run(); err != nil {
		// 	fmt.Println("Error running program:", err)
		// 	os.Exit(1)
		// }

		// completeResp, r, err := client.UserApi.UserLoginCompleteGet(context.Background()).Token(*resp.Token).Client(clientId).Execute()
		// fmt.Printf("Complete response: %+v, r: %+v\n", completeResp, r)
		// fmt.Println("Welcome to launch!")
		// user, _, err := apiclient.C.UserApi.UserGet(context.Background()).Execute()
		// if err != nil {
		// 	fmt.Println("Error getting user:", err)
		// }
		// fmt.Printf("User: %+v\n", *user)

		//list.New()

		m := tui.NewList([]list.Item{
			tui.ListItem("An existing app"),
			tui.ListItem("My app"),
			tui.ListItem("A static website"),
		}, "Welcome to Podinate. What would you like to launch?", "Launching %s")
		_, err := tea.NewProgram(m).Run()
		if err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}

		// fmt.Println("You chose:", fm.(tui.ListModel).Choice())
		// MVP, I don't give a fuck what you choose

		m = tui.NewList([]list.Item{
			tui.ListItem("WordPress"),
			tui.ListItem("Gitea"),
			tui.ListItem("Visual Stuio Code Web"),
			tui.ListItem("Matomo Analytics"),
			tui.ListItem("Mattermost"),
			tui.ListItem("Nextcloud"),
		}, "What app would you like to launch?", "Launching %s")
		_, err = tea.NewProgram(m).Run()
		if err != nil {
			fmt.Println("Error choosing app")
			os.Exit(1)
		}
		wptf := []byte(`
terraform {
	required_providers {
		podinate = {
		source  = "podinate/podinate"
		version = "0.0.1"
		}
	}
}

provider "podinate" {
	# Configuration options
	server_url   = "http://localhost:3001/v0"
	api_key_auth = var.podinate_api_key
}

resource "podinate_project" "wordpress_project" {
	account = var.account_id
	id      = replace(lower(var.project_name), " ", "-")
	name    = var.project_name
}

resource "podinate_pod" "wordpress_pod" {
	account    = podinate_project.wordpress_project.account
	project_id = podinate_project.wordpress_project.id
	id         = "wordpress"
	image      = "wordpress"
	name       = "WordPress"
	tag        = "6"
	environment = [
		{
		key   = "WORDPRESS_DB_HOST"
		value = "mysql"
		},
		{
		key   = "WORDPRESS_DB_USER"
		value = "wordpress"
		},
		{
		key    = "WORDPRESS_DB_PASSWORD"
		value  = var.db_pass
		secret = true
		},
		{
			key = "WORDPRESS_DB_NAME"
			value = "wordpress"
		}
	]
	services = [
		{
			name     = "wordpress"
			port     = 80
			protocol = "http"
			domain_name = "john-blog.podinate.app"
		}
	]
}

resource "podinate_pod" "database_pod" {
	account    = podinate_project.wordpress_project.account
	project_id = podinate_project.wordpress_project.id
	id         = "mariadb"
	image      = "mariadb"
	name       = "MySQL"
	tag        = "10"
	environment = [
		{
		key   = "MARIADB_USER"
		value = "wordpress"
		},
		{
		key    = "MARIADB_PASSWORD"
		value  = var.db_pass
		secret = true
		},
		{
		key   = "MARIADB_DATABASE"
		value = "wordpress"
		},
		{
			key = "MARIADB_RANDOM_ROOT_PASSWORD"
			value = "yesn't'd've"
		}
	]
	services = [
		{
			name     = "mysql"
			port     = 3306
			protocol = "tcp"
		}
	]
}
		  
		`)

		tfvars := `
variable "podinate_api_key" {
	type = string 
	description = "API key for Podinate"
}

variable "account_id" {
	type = string
	description = "Account ID for Podinate"
}

variable "project_name" {
	type = string
	default = "WordPress Blog"
	description = "Project name for Podinate"
}
variable "db_pass" {
	type = string
	description = "Password for the WordPress database"
}
		`

		tfvals := "project_name = \"WordPress Blog\"\ndb_pass = \"awD7iw!vNTEgP^!qa%fq\""
		podinateYaml := "project: \"wordpress-blog\""

		os.WriteFile("./podinate.tf", wptf, 0644)
		os.WriteFile("./variables.tf", []byte(tfvars), 0644)
		os.WriteFile("./terraform.tfvars", []byte(tfvals), 0644)
		os.WriteFile("./podinate.yaml", []byte(podinateYaml), 0644)

		fmt.Println("Wrote Terraform files, when you're ready, run 'podinate tf apply' to launch your app")

	},
}
