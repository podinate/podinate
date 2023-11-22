package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/johncave/podinate/cli/apiclient"
	tui "github.com/johncave/podinate/cli/tui"
	"github.com/johncave/podinate/cli/tui/spinner"
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

		// Determine which login provider to use
		m := tui.NewList([]list.Item{
			//tui.ListItem("GitHub"),
			//tui.ListItem("GitLab"),
			tui.ListItem("Podinate's GitLab"),
		}, "Welcome to Podinate. How would you like to log in?")

		fm, err := tea.NewProgram(m).Run()
		if err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}

		// User has chosen how to log in, now start the login process
		provider := strings.ToLower(fm.(tui.ListModel).Choice())
		if provider == "podinate's gitlab" {
			provider = "podinate"
		}

		resp, r, err := apiclient.InitLogin(provider)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		//fmt.Printf("Response: %+v, r: %+v\n", resp, r)
		//fmt.Printf("URL: %s, Token: %s\n", *resp.Url, *resp.Token)
		err = open(*resp.Url)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		sp := tea.NewProgram(spinner.New("Awaiting login, check your browser...", "points"))
		go sp.Run()

		// go if _, err := tea.NewProgram(spinner.New("Awaiting login, check your browser...", "points")).Run(); err != nil {
		// 	fmt.Println("Error running program:", err)
		// 	os.Exit(1)
		// }

		// completeResp, r, err := client.UserApi.UserLoginCompleteGet(context.Background()).Token(*resp.Token).Client(clientId).Execute()
		// fmt.Printf("Complete response: %+v, r: %+v\n", completeResp, r)
		completeResp, r, err := apiclient.AwaitCompleteLogin(resp)

		if e := new(string); completeResp.ApiKey == e {
			fmt.Println("Error: Timed out getting api key")
			os.Exit(1)
		}

		sp.Quit()
		//fmt.Println("Finished retrying")
		// fmt.Printf("Done! API Key is: %s", *completeResp.ApiKey)
		fmt.Println("You are logged in")
	},
}

// https://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang
// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	//fmt.Printf("Opening %s in browser\n", url)
	//fmt.Printf("cmd: %s, args: %+v\n", cmd, args)
	return exec.Command(cmd, args...).Start()
}
