// File: apiclient.go
// Package: apiclient
// Abstracts away the API client so that it can be used in the CLI
package apiclient

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/avast/retry-go"
	api "github.com/johncave/podinate/lib/api_client"
)

var (
	C *api.APIClient
)

func init() {
	config := api.NewConfiguration()
	config.Host = "localhost:3001"
	config.Scheme = "http"
	config.UserAgent = "Podinate CLI"

	client := api.NewAPIClient(config)
	C = client
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
	fmt.Printf("%+v\n", viper.)
	return resp, r, err
}

type GlobalConfigFile struct {
	Profiles: []Profile
}

type Profile struct {
	Name string	`yaml:"name"`
	APIKey string `yaml:"api_key"`