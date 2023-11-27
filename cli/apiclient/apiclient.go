// File: apiclient.go
// Package: apiclient
// Abstracts away the API client so that it can be used in the CLI
package apiclient

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/johncave/podinate/cli/tui"
	"github.com/johncave/podinate/cli/tui/spinner"
	api "github.com/johncave/podinate/lib/api_client"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	C      *api.APIClient
	config *api.Configuration
)

func init() {
	config = api.NewConfiguration()
	config.Host = "localhost:3001"
	config.Scheme = "http"
	config.UserAgent = "Podinate CLI"

	client := api.NewAPIClient(config)
	C = client
}

func SetupUser() {
	key, ex := os.LookupEnv("PODINATE_API_KEY")
	if ex {
		config.DefaultHeader["Authorization"] = key
		return
	}

	// fmt.Printf("Config: %+v\n", viper.AllKeys())
	// fmt.Printf("Config: %+v\n", )

	config.DefaultHeader["Authorization"] = viper.GetString("profiles.0.api_key")
}

// InitLogin starts the login process
func InitLogin(provider string) (*api.UserLoginInitiateGet200Response, *http.Response, error) {
	resp, r, err := C.UserApi.UserLoginInitiateGet(nil).Provider(provider).Execute()
	return resp, r, err
}

// AwaitCompleteLogin waits for the user to complete the login process
func AwaitCompleteLogin(initresp *api.UserLoginInitiateGet200Response) (*api.UserLoginCompleteGet200Response, *http.Response, error) {
	hostname, _ := os.Hostname()
	clientId := "podinate-cli on " + hostname
	var resp *api.UserLoginCompleteGet200Response
	var r *http.Response
	var err error

	retry.Do(
		func() error {
			//fmt.Println("Retrying...")
			resp, r, err = C.UserApi.UserLoginCompleteGet(context.Background()).Token(*initresp.Token).Client(clientId).Execute()
			if err != nil {
				return err
			}
			//fmt.Printf("Login not complete: %+v, r: %+v\n", completeResp, r)
			if r.StatusCode != 200 {
				return fmt.Errorf("status code was not 200")
			}
			//fmt.Printf("Complete response: %+v, r: %+v\n", completeResp, r)
			return nil
		},
		retry.Delay(time.Duration(2*time.Second)),
		retry.MaxJitter(time.Duration(1*time.Second)),
		retry.MaxDelay(time.Duration(3*time.Second)),
		retry.Attempts(100))

	// Save the token to the config file
	fmt.Printf("%+v\n", viper.AllKeys())

	// Save the token to the config file
	err = SaveProfile("default", *resp.ApiKey)
	if err != nil {
		fmt.Println("Error saving profile:", err)
		os.Exit(1)
	}
	viper.Set("user.key", *resp.ApiKey)
	config.DefaultHeader["Authorization"] = *resp.ApiKey

	return resp, r, err

}

type GlobalConfigFile struct {
	Profiles []Profile `yaml:"profiles"`
}

type Profile struct {
	Name   string `yaml:"name"`
	APIKey string `yaml:"api_key"`
}

// SaveProfile saves the profile to the config file
func SaveProfile(profileName string, apiKey string) error {
	// TODO: Check for existing profiles

	yamlContent, err := yaml.Marshal(GlobalConfigFile{
		Profiles: []Profile{
			{
				Name:   profileName,
				APIKey: apiKey,
			},
		},
	})
	if err != nil {
		return err
	}
	//fmt.Println("Saving profile to", viper.ConfigFileUsed())
	filePath := viper.ConfigFileUsed()
	if filePath == "" {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		filePath = home + "/.config/podinate/credentials.yaml"
	}
	fmt.Printf("Saving profile to %s\n", filePath)
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(yamlContent)
	if err != nil {
		return err
	}
	return nil
}

func DoLogin() {
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

	resp, r, err := InitLogin(provider)
	if err != nil {
		fmt.Println("Error:", err, r)
		os.Exit(1)
	}

	err = openBrowser(*resp.Url)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	sp := tea.NewProgram(spinner.New("Awaiting login, check your browser...", "points"))
	go sp.Run()

	completeResp, r, err := AwaitCompleteLogin(resp)

	if e := new(string); completeResp.ApiKey == e {
		fmt.Println("Error: Timed out getting api key")
		os.Exit(1)
	}

	sp.Quit()
	//fmt.Println("Finished retrying")
	// fmt.Printf("Done! API Key is: %s", *completeResp.ApiKey)
	fmt.Println("You are logged in")
}

// https://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang
// open opens the specified URL in the default browser of the user.
func openBrowser(url string) error {
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
	//e := exec.Command(cmd, args...)

	return exec.Command(cmd, args...).Start()
}
