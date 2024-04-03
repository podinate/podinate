// File: apiclient.go
// Package: apiclient
// Abstracts away the API client so that it can be used in the CLI
package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/avast/retry-go"
	tea "github.com/charmbracelet/bubbletea"
	textarea "github.com/johncave/podinate/cli/tui/textarea"
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

	// Find the profile
	// fmt.Printf("Config: %+v\n", viper.GetStringMap("profiles"))
	//m := viper.GetStringMap("profiles")
	//m := viper.GetStringMapString("profiles")
	//fmt.Printf("Config: %+v\n", m)

	// Find home directory.
	configFile, err := readConfigFile()
	if err != nil {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}

	for i, p := range configFile.Profiles {
		if p.Name == viper.GetString("profile") {
			if viper.GetBool("verbose") {
				log.Printf("Using profile: %+v\n", p)
			}
			// I don't know why tf the apikey isn't in the profile, but this works
			apiKey := viper.GetString("profiles." + strconv.Itoa(i) + ".api_key")
			config.DefaultHeader["Authorization"] = apiKey

			// Set the API URL
			u, err := url.Parse(p.APIUrl)
			if err != nil {
				fmt.Println("Error parsing API URL in profile:", err)
				os.Exit(1)
			}
			config.Host = u.Host
			config.Scheme = u.Scheme

			viper.Set("api_key", apiKey)
			return
		}
	}

	log.Printf("Profile not found: %s\n", viper.GetString("profile"))
	os.Exit(1)
	return

	// baseURL := viper.GetString("api_url")
	// u, err := url.Parse(baseURL)
	// if err != nil {
	// 	fmt.Println("Error parsing API URL in profile:", err)
	// 	os.Exit(1)
	// }

	// fmt.Printf("Config: %+v\n", viper.AllKeys())
	// fmt.Printf("Config: %+v\n", )

	// fmt.Printf("API Key: %s\n", viper.GetString("profiles.0.api_key"))
	// config.DefaultHeader["Authorization"] = viper.GetString("profiles.0.api_key")
	// viper.Set("api_key", viper.GetString("profiles.0.api_key"))
}

// InitLogin starts the login process
func InitLogin(provider string) (*api.UserLoginInitiateGet200Response, *http.Response, error) {
	resp, r, err := C.UserApi.UserLoginInitiateGet(nil).Provider(provider).Execute()
	return resp, r, err
}

// AwaitCompleteLogin waits for the user to complete the login process
func AwaitCompleteLogin(initresp *api.UserLoginInitiateGet200Response) (*api.UserLoginPost200Response, *http.Response, error) {
	hostname, _ := os.Hostname()
	clientId := "podinate-cli on " + hostname
	var resp *api.UserLoginPost200Response
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
	err = SaveProfile("http://localhost:3001", "default", *resp.ApiKey)
	if err != nil {
		fmt.Println("Error saving profile:", err)
		os.Exit(1)
	}
	viper.Set("api_key", *resp.ApiKey)
	config.DefaultHeader["Authorization"] = *resp.ApiKey

	return resp, r, err

}

type GlobalConfigFile struct {
	Profiles []Profile `yaml:"profiles"`
}

type Profile struct {
	Name   string `yaml:"name"`
	APIKey string `yaml:"api_key"`
	APIUrl string `yaml:"api_url"`
}

func SaveYamlProfile(in string) error {
	var profile Profile
	err := yaml.Unmarshal([]byte(in), &profile)
	if err != nil {
		return err
	}
	return SaveProfile(profile.APIUrl, profile.Name, profile.APIKey)
}

// SaveProfile saves the profile to the config file
func SaveProfile(apiURL string, profileName string, apiKey string) error {
	// TODO: Check for existing profiles

	currentConfig, err := readConfigFile()
	if err != nil && !errors.Is(err, os.ErrNotExist) { // If the file doesn't exist, that's fine
		return err
	}

	newProfile := Profile{
		Name:   profileName,
		APIKey: apiKey,
		APIUrl: apiURL,
	}

	overwrote := false
	// Skip if no current config
	if currentConfig != nil {
		for i, p := range currentConfig.Profiles {
			if p.Name == profileName {
				if viper.GetBool("verbose") {
					log.Printf("Profile already exists: %+v, overwriting\n", p)
				}
				currentConfig.Profiles[i] = newProfile
				overwrote = true
			}
		}
	} else {
		currentConfig = &GlobalConfigFile{}
	}

	if !overwrote {
		currentConfig.Profiles = append(currentConfig.Profiles, newProfile)
	}

	if viper.GetBool("verbose") {
		log.Printf("Saving profile: %+v\n", currentConfig)
	}

	yamlContent, err := yaml.Marshal(currentConfig)
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

		err = os.MkdirAll(home+"/.config/podinate", os.ModePerm)
		if err != nil {
			return err
		}
		filePath = home + "/.config/podinate/credentials.yaml"
	}
	fmt.Printf("Saving profile to %s\n", filePath)
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
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

func StartLogin() {
	// This is all related to oauth login which we are canning for now
	// // Determine which login provider to use
	// m := tui.NewList([]list.Item{
	// 	//tui.ListItem("GitHub"),
	// 	//tui.ListItem("GitLab"),
	// 	tui.ListItem("Podinate's GitLab"),
	// }, "Welcome to Podinate. How would you like to log in?", "Logging in with %s")

	// fm, err := tea.NewProgram(m).Run()
	// if err != nil {
	// 	fmt.Println("Error running program:", err)
	// 	os.Exit(1)
	// }

	// // User has chosen how to log in, now start the login process
	// provider := strings.ToLower(fm.(tui.ListModel).Choice())
	// if provider == "podinate's gitlab" {
	// 	provider = "podinate"
	// }

	// resp, r, err := InitLogin(provider)
	// if err != nil {
	// 	fmt.Println("Error:", err, r)
	// 	os.Exit(1)
	// }

	// err = openBrowser(*resp.Url)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	os.Exit(1)
	// }

	// sp := tea.NewProgram(spinner.New("Awaiting login, check your browser...", "points"))
	// go sp.Run()

	// completeResp, r, err := AwaitCompleteLogin(resp)

	// if e := new(string); completeResp.ApiKey == e {
	// 	fmt.Println("Error: Timed out getting api key")
	// 	os.Exit(1)
	// }

	// sp.Quit()
	//fmt.Println("Finished retrying")
	// fmt.Printf("Done! API Key is: %s", *completeResp.ApiKey)

	// input := tea.NewProgram(textinput.New("Podinate API URL", "https://localhost:31337"))
	// input.Run()

	// input = tea.NewProgram(textinput.New("Username", ""))

	// var value *string

	var profile string

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	// Checking to see if there is data being piped in
	if fi.Mode()&os.ModeNamedPipe == 0 {
		// No data being piped in, so we need to get the data from the user

		model := textarea.New()

		input := tea.NewProgram(model)
		m, err := input.Run()
		if err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}

		profile = m.(textarea.Model).Value

	} else {
		// Read the piped-in profile
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Error reading profile from stdin:", err)
			os.Exit(1)
		}

		profile = string(bytes)
	}

	if viper.GetBool("verbose") {
		fmt.Println(profile)
	}

	if profile == "" {
		fmt.Println("No profile provided")
		os.Exit(1)
	}

	err = SaveYamlProfile(profile)
	if err != nil {
		fmt.Println("Not a valid profile:", err)
		os.Exit(1)
	}

	// var profileString string
	// // var username string

	// fmt.Print("Paste your Podinate profile file:\n\n")
	// fmt.Scanln(&profileString)

	// fmt.Print("Username: ")
	// fmt.Scanln(&username)

	// fmt.Print("Password: ")
	// pass, err := terminal.ReadPassword(0)
	// if err != nil {
	// 	fmt.Println("Error reading password:", err)
	// 	os.Exit(1)
	// }

	// doLogin(apiURL, username, string(pass))

	fmt.Println("You are logged in")
}

func doLogin(apiURL string, username string, password string) {
	u, err := url.Parse(apiURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		os.Exit(1)
	}
	config = api.NewConfiguration()
	config.Host = u.Host
	config.Scheme = u.Scheme
	config.UserAgent = "Podinate CLI"

	client := api.NewAPIClient(config)
	C = client

	userLoginRequest := *api.NewUserLoginPostRequest(username, password)
	hostname, _ := os.Hostname()
	clientId := "podinate-cli on " + hostname
	userLoginRequest.Client = &clientId

	resp, _, err := C.UserApi.UserLoginPost(nil).UserLoginPostRequest(userLoginRequest).Execute()
	if err != nil {
		fmt.Println("Error signing in", err)
		os.Exit(1)
	}

	SaveProfile(apiURL, viper.GetString("profile"), *resp.ApiKey)

	// resp, r, err := C.UserApi.UserLoginPostRequest(context.Background()).UserLoginPostRequest(api.UserLoginPostRequest{
	// 	Username: username,
	// 	Password: password,
	// }).Execute()

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

// readConfigFile reads the config file and returns the GlobalConfigFile struct
func readConfigFile() (*GlobalConfigFile, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	var configFile string
	if viper.GetString("config_file") != "" {
		configFile = viper.GetString("config_file")
	} else {
		configFile = home + "/.config/podinate/credentials.yaml"
	}

	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var profiles GlobalConfigFile
	err = yaml.Unmarshal(yamlFile, &profiles)
	if err != nil {
		return nil, err
	}
	return &profiles, nil
}

// handleAPIError handles API errors
func handleAPIError(r *http.Response, err error) *SDKError {
	//fmt.Println("handleAPIError", r, err)
	if r == nil {
		return newError(0, "Could not connect to Podinate API: "+err.Error())
	}
	if err != nil {
		var sdkerr *SDKError
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return newError(http.StatusTeapot, "Couldn't read response from API! "+string(body))
		}
		//jsonError := json.NewDecoder(r.Body).Decode(&sdkerr)
		jsonError := json.Unmarshal(body, &sdkerr)
		if jsonError != nil {
			if r != nil {
				return newError(r.StatusCode, "Got invalid error response from API! "+string(body))
			} else {
				return newError(0, "Got invalid error response from API! "+string(body))
			}
		}
		return sdkerr
	}
	return nil
}

type SDKError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error returns the error message - satisfy the error interface
func (e *SDKError) Error() string {
	return strconv.Itoa(e.Code) + ": " + e.Message
}

// NewError creates a new SDKError
func newError(code int, message string) *SDKError {
	return &SDKError{
		Code:    code,
		Message: message,
	}
}
