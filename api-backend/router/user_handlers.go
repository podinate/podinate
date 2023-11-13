package router

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/johncave/podinate/api-backend/config"
	api "github.com/johncave/podinate/api-backend/go"
	"github.com/johncave/podinate/api-backend/responder"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/gitlab"
)

type UserApiController struct {
	*api.UserApiController
}

// Import the default service struct
type UserAPIService struct {
	api.UserApiService
}

// NewUserAPIService creates a new service for handling requests in the default API
func NewUserAPIService() api.UserApiServicer {
	return &UserAPIService{}
}

// Build the service methods

// UserGet - Get information about the current user
func (s *UserAPIService) UserGet(ctx context.Context, account string) (api.ImplResponse, error) {
	// TODO - update UserGet with the required logic for this service method.
	return api.Response(501, nil), nil
}

// UserLoginCompleteGet - Complete the login process
func (s *UserAPIService) UserLoginCompleteGet(ctx context.Context, token string) (api.ImplResponse, error) {
	SetUpGoth()
	// TODO - update UserLoginCompleteGet with the required logic for this service method.

	// err =

	return api.Response(501, nil), nil
}

// UserLoginCallbackProviderGet - Handle the callback from the provider
func (s *UserAPIService) UserLoginCallbackProviderGet(ctx context.Context, providerName string) (api.ImplResponse, error) {
	SetUpGoth()

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return responder.Response(400, err.Error()), err
	}

	value, err := GetFromSession(token, providerName)
	if err != nil {
		return responder.Response(500, err.Error()), err
	}

	session, err := provider.UnmarshalSession(value)
	if err != nil {
		return responder.Response(500, err.Error()), err
	}
	if session == nil {
		return responder.Response(500, "Session is nil"), err
	}

}

// UserLoginRedirectTokenGet - Redirect the user to the provider
func (s *UserAPIService) UserLoginRedirectTokenGet(ctx context.Context, token string) (api.ImplResponse, error) {
	SetUpGoth()
	// TODO - update UserLoginRedirectTokenGet with the required logic for this service method.
	providerName, err := GetProviderFromSession(token)
	if err != nil {
		return responder.Response(500, err.Error()), nil
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return responder.Response(400, err.Error()), nil
	}

	value, err := GetFromSession(token, providerName)
	if err != nil {
		return responder.Response(500, err.Error()), nil
	}

	session, err := provider.UnmarshalSession(value)
	if err != nil {
		return responder.Response(500, err.Error()), nil
	}

}

// UserLoginInitiateGet - Initiate the login process
func (s *UserAPIService) UserLoginInitiateGet(ctx context.Context, providerName string) (api.ImplResponse, error) {
	SetUpGoth()
	// TODO - update UserLoginInitiateGet with the required logic for this service method.
	//gothic.BeginAuthHandler()

	// Create a session ID for the login
	sessionID := uuid.New().String()

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return responder.Response(400, err.Error()), nil
	}

	session, err := provider.BeginAuth(sessionID)
	if err != nil {
		return responder.Response(500, err.Error()), nil
	}

	url, err := session.GetAuthURL()
	if err != nil {
		return responder.Response(500, err.Error()), nil
	}

	err = StoreInSession(sessionID, providerName, session.Marshal())
	if err != nil {
		return responder.Response(500, err.Error()), nil
	}

	out := struct {
		Url   string `json:"url"`
		Token string `json:"token"`
	}{
		Url:   url,
		Token: sessionID,
	}

	return responder.Response(200, out), nil
}

// SetUpGoth sets up the Goth library and registers the providers
func SetUpGoth() {

	gitlab.AuthURL = "http://gitlab.dev-git.podinate.com:8081/oauth/authorize"
	gitlab.TokenURL = "http://gitlab-webservice-default.gitlab:8080/oauth/token"
	gitlab.ProfileURL = "http://gitlab-webservice-default.gitlab:8080/api/v3/user"

	// TODO: Make callback URL configurable
	podinateProvider := gitlab.New("b6f324560f50b728493a021a4cb529d52c0aa4c49af7a54c2248ed7ddc9d7e68", "f769bd4869b6d9d81fbe287fd99ff4db6601aca9214463a246193194a1cf5858", "http://localhost:3001/user/login/callback/podinate", "read_user", "read_repository", "read_registry", "profile", "email")
	podinateProvider.SetName("podinate")

	goth.UseProviders(podinateProvider)
	fmt.Printf("%+v", goth.GetProviders())
}

// StoreInSession stores a value in the session
func StoreInSession(sessionID string, key string, value string) error {
	// Store in postgres
	_, err := config.DB.Exec("INSERT INTO login_session (session_id, key, value) VALUES ($1, $2, $3) ON CONFLICT ON CONSTRAINT composite_primary DO UPDATE SET value = EXCLUDED.value", sessionID, key, value)
	if err != nil {
		return err
	}
	return nil
}

// GetFromSession gets a value from the session
func GetFromSession(sessionID string, key string) (string, error) {
	var value string
	err := config.DB.QueryRow("SELECT value FROM login_session WHERE session_id = $1 AND key = $2", sessionID, key).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

// GetProviderFromSession gets a provider from the session
func GetProviderFromSession(sessionID string) (string, error) {
	var value string
	err := config.DB.QueryRow("SELECT key FROM login_session WHERE session_id = $1", sessionID).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}
