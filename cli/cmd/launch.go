package cmd

import (
	"context"
	"fmt"

	"github.com/johncave/podinate/cli/apiclient"
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
		fmt.Println("Welcome to launch!")
		user, _, err := apiclient.C.UserApi.UserGet(context.Background()).Execute()
		if err != nil {
			fmt.Println("Error getting user:", err)
		}
		fmt.Printf("User: %+v\n", *user)

	},
}
