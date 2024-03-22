package cmd

import (
	"fmt"
	"os"

	bubbletable "github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/johncave/podinate/cli/sdk"
	"github.com/johncave/podinate/cli/tui/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	getCmd.AddCommand(getPodsCmd)
	getCmd.AddCommand(getProjectsCmd)
	rootCmd.AddCommand(getCmd)
	getCmd.PersistentFlags().StringP("format", "f", "table", "output format, pick from table, json, yaml")
}

var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"ls", "l", "list", "view"},
	Short:   "List things on Podinate",
	Long:    `Lists things on Podinate`,
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(cmd.Help())
	},
}

var getPodsCmd = &cobra.Command{
	Use:     "pods",
	Aliases: []string{"pod", "po"},
	Short:   "List pods on Podinate",
	Long:    `Lists all pods on Podinate`,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("verbose") {
			fmt.Printf("Getting pods from project %s and account %s\n", viper.GetString("project"), viper.GetString("account"))
		}

		if viper.GetString("project") == "" {
			fmt.Println("Please specify a project with --project or -p")
			os.Exit(1)
		}

		resp, _, err := sdk.C.PodApi.ProjectProjectIdPodGet(cmd.Context(), viper.GetString("project")).Account(viper.GetString("account")).Execute()

		if err != nil && err.Error() == "404 Not Found" {
			fmt.Println("No pods found")
			os.Exit(0)
		}
		//fmt.Printf("Error: %T\n", err)
		cobra.CheckErr(err)
		//out, _ := json.Marshal(resp)
		//fmt.Printf("Response: %s, r: %+v\n", out, r)

		if len(resp.Items) < 1 {
			fmt.Println("No pods found")
			os.Exit(0)
		}

		columns := []bubbletable.Column{
			{Title: "ID", Width: 15},
			{Title: "Name", Width: 20},
			{Title: "Status", Width: 10},
			{Title: "Image", Width: 20},
		}

		var rows []bubbletable.Row

		for _, i := range resp.Items {
			p := i.Pod

			rows = append(rows, bubbletable.Row{
				p.Id, p.Name, *p.Status, p.Image + ":" + p.Tag,
			})
		}

		m := table.New(columns, rows)

		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}

var getProjectsCmd = &cobra.Command{
	Use:     "projects",
	Aliases: []string{"project", "proj", "projs"},
	Short:   "List projects on Podinate",
	Long:    `Lists all projects on Podinate account`,
	Run: func(cmd *cobra.Command, args []string) {
		projects, err := sdk.GetAllProjects(viper.GetString("account"))

		//fmt.Printf("Response: %+v, r: %+v\n", resp, r)
		cobra.CheckErr(err)

		if len(projects) < 1 {
			fmt.Println("No projects found")
			os.Exit(0)
		}

		columns := []bubbletable.Column{
			{Title: "ID", Width: 15},
			{Title: "Name", Width: 20},
		}

		var rows []bubbletable.Row

		for _, p := range projects {

			rows = append(rows, bubbletable.Row{
				p.ID, p.Name,
			})
		}

		m := table.New(columns, rows)

		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}

type Showable interface {
	GetItems() map[string]string
}

// func Shower()
