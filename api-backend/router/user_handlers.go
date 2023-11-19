package router

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/johncave/podinate/api-backend/config"
	api "github.com/johncave/podinate/api-backend/go"
	"github.com/johncave/podinate/api-backend/iam"
	lh "github.com/johncave/podinate/api-backend/loghandler"
	"github.com/johncave/podinate/api-backend/responder"
	myuser "github.com/johncave/podinate/api-backend/user"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/gitlab"
)

// validateState ensures that the state token param from the original
// AuthURL matches the one included in the current (callback) request.
func validateState(req *http.Request, sess goth.Session) error {
	rawAuthURL, err := sess.GetAuthURL()
	if err != nil {
		return err
	}

	authURL, err := url.Parse(rawAuthURL)
	if err != nil {
		return err
	}

	reqState := GetState(req)

	originalState := authURL.Query().Get("state")
	if originalState != "" && (originalState != reqState) {
		return errors.New("state token mismatch")
	}
	return nil
}

// GetState gets the state returned by the provider during the callback.
// This is used to prevent CSRF attacks, see
// http://tools.ietf.org/html/rfc6749#section-10.12
var GetState = func(req *http.Request) string {
	params := req.URL.Query()
	if params.Encode() == "" && req.Method == http.MethodPost {
		return req.FormValue("state")
	}
	return params.Get("state")
}

// Import the default service struct
type UserAPIService struct {
	api.UserApiService
}

// NewUserAPIService creates a new service for handling requests in the default API
func NewUserAPIService() api.UserApiServicer {
	return &UserAPIService{}
}

// UserGet - Get information about the logged in user
func (s *UserAPIService) UserGet(ctx context.Context) (api.ImplResponse, error) {
	u := iam.GetFromContext(ctx).(*myuser.User)
	return api.Response(200, u.ToAPI()), nil
}

// UserLoginCompleteGet - Complete the login process
func (s *UserAPIService) UserLoginCompleteGet(ctx context.Context, token string, clientName string) (api.ImplResponse, error) {
	// Check the name isn't 2TB long
	if len(clientName) > 2048 {
		lh.Log.Errorw("Client gave a name too long", "name", clientName)
		return responder.Response(400, "Client name too long"), nil
	}

	userid, err := GetFromSession(token, "authorised_user")
	if err != nil {
		// TODO: Separate handling for token not found vs other errors
		if err == sql.ErrNoRows {
			lh.Log.Errorw("No authorised user found in session", "token", token)
			return responder.Response(403, "Invalid login token"), nil
		}

		lh.Log.Errorw("Error getting authorised user from session", "error", err, "token", token)
		return responder.Response(500, err.Error()), nil
	}

	user, err := myuser.GetByUUID(userid)
	if err != nil {
		lh.Log.Errorw("Error getting user for token", "error", err, "user", user, "uuid", userid)
		return responder.Response(500, err.Error()), nil
	}

	// Issue an API key for the user
	apiKey, err := user.IssueAPIKey(clientName)
	if err != nil {
		lh.Log.Errorw("Error issuing API key in exchange for token", "error", err, "user", user, "uuid", userid)
		return responder.Response(500, err.Error()), nil
	}

	// API key issued! Remove the reference to the token from the session
	err = StoreInSession(token, "authorised_user", "")
	if err != nil {
		lh.Log.Errorw("Error removing authorised user from session", "error", err, "user", user, "uuid", userid)
		return responder.Response(500, err.Error()), nil
	}

	resp := api.UserLoginCompleteGet200Response{
		ApiKey:   apiKey,
		LoggedIn: true,
	}

	lh.Log.Infow("User logged in", "user", user)

	return responder.Response(200, resp), nil
}

// UserLoginCallbackProviderGet - Handle the callback from the provider
func (s *UserAPIService) UserLoginCallbackProviderGet(ctx context.Context, providerName string) (api.ImplResponse, error) {
	return api.Response(501, nil), nil
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

	// url, err := session.GetAuthURL()
	// if err != nil {
	// 	return responder.Response(500, err.Error()), nil
	// }

	url := fmt.Sprintf("http://localhost:3001/v0/user/login/redirect/%s", sessionID)

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
	podinateProvider := gitlab.New("b6f324560f50b728493a021a4cb529d52c0aa4c49af7a54c2248ed7ddc9d7e68", "f769bd4869b6d9d81fbe287fd99ff4db6601aca9214463a246193194a1cf5858", "http://localhost:3001/v0/user/login/callback/podinate", "read_user", "read_repository", "read_registry", "profile", "email")
	podinateProvider.SetName("podinate")

	goth.UseProviders(podinateProvider)
	fmt.Printf("%+v", goth.GetProviders())
}

// StoreInSession stores a value in the session
func StoreInSession(sessionID string, key string, value string) error {
	// If value is empty string, remove from key value store
	if value == "" {
		_, err := config.DB.Exec("DELETE FROM login_session WHERE session_id = $1 AND key = $2", sessionID, key)
		if err != nil {
			return err
		}
		return nil
	}

	// If value not empty, store in postgres
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
