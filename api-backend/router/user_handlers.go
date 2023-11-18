package router

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/johncave/podinate/api-backend/config"
	eh "github.com/johncave/podinate/api-backend/errorhandler"
	api "github.com/johncave/podinate/api-backend/go"
	"github.com/johncave/podinate/api-backend/responder"
	myuser "github.com/johncave/podinate/api-backend/user"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/gitlab"
)

type UserAPIShim struct {
	service      api.UserApiServicer
	errorHandler api.ErrorHandler
}

func (c *UserAPIShim) Routes() api.Routes {
	return api.Routes{
		{
			Name:        "UserLoginRedirectTokenGet",
			Method:      "GET",
			Pattern:     "/v0/user/login/redirect/{token}",
			HandlerFunc: c.UserLoginRedirectTokenGet,
		},
		{
			Name:        "UserLoginCallbackProviderGet",
			Method:      "GET",
			Pattern:     "/v0/user/login/callback/{provider}",
			HandlerFunc: c.UserLoginCallbackProviderGet,
		},
	}
}

// UserLoginRedirectTokenGet - Redirect the user to the provider
func (c *UserAPIShim) UserLoginRedirectTokenGet(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Handling redirect")

	// Grab the token from the URL
	params := mux.Vars(r)
	token := params["token"]

	// Check the token is a valid uuidv4
	_, err := uuid.Parse(token)
	if err != nil {
		api.EncodeJSONResponse(responder.Response(400, "Invalid redirect token"), nil, w)
		return
	}

	code := 500 // Default to 500
	providerName, err := GetProviderFromSession(token)
	if err != nil {
		api.EncodeJSONResponse(responder.Response(500, err.Error()), &code, w)
		return
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		code = 400
		api.EncodeJSONResponse(responder.Response(400, err.Error()), &code, w)
		return
	}

	value, err := GetFromSession(token, providerName)
	if err != nil {
		api.EncodeJSONResponse(responder.Response(500, err.Error()), &code, w)
		return
	}

	session, err := provider.UnmarshalSession(value)
	if err != nil {
		api.EncodeJSONResponse(responder.Response(500, err.Error()), &code, w)
		return
	}

	// Redirect the user to the provider
	authUrl, err := session.GetAuthURL()
	if err != nil {
		api.EncodeJSONResponse(responder.Response(500, err.Error()), &code, w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "login-token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   3600,
	})
	http.Redirect(w, r, authUrl, http.StatusTemporaryRedirect)
	w.Write([]byte("You should now be redirected. If not click <a href=\"" + authUrl + "\">here</a>"))
}

// UserLoginCallbackProviderGet - Handle the callback from the provider
func (c *UserAPIShim) UserLoginCallbackProviderGet(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Grab the provider from the URL
	params := mux.Vars(r)
	providerName := params["provider"]

	fmt.Printf("Handling callback for %+v %s\n", params, providerName)
	//result, err := c.service.UserLoginCallbackProviderGet(r, providerName)

	SetUpGoth()

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		api.EncodeJSONResponse(responder.Response(400, err.Error()), nil, w)
		return
	}

	cookie, err := r.Cookie("login-token")
	if err != nil {
		api.EncodeJSONResponse(responder.Response(403, "No login cookie found"), nil, w)
		return
	}

	value, err := GetFromSession(cookie.Value, providerName)
	if err != nil {
		api.EncodeJSONResponse(responder.Response(500, err.Error()), nil, w)
		return
	}

	session, err := provider.UnmarshalSession(value)
	if err != nil {
		api.EncodeJSONResponse(responder.Response(500, err.Error()), nil, w)
		return
	}
	if session == nil {
		api.EncodeJSONResponse(responder.Response(500, "No session"), nil, w)
		return
	}

	err = validateState(r, session)
	if err != nil {
		api.EncodeJSONResponse(responder.Response(500, err.Error()), nil, w)
		return
	}

	user, err := provider.FetchUser(session)
	if err == nil {
		// We already found the user information, so it can be displayed
		fmt.Printf("Success %+v\n", user)
	} else {
		// We need to fetch the user information from the provider
		_, err = session.Authorize(provider, query)
		if err != nil {
			api.EncodeJSONResponse(responder.Response(500, err.Error()), nil, w)
			return
		}
		user, err = provider.FetchUser(session)
		if err != nil {
			api.EncodeJSONResponse(responder.Response(500, err.Error()), nil, w)
			return
		}
	}

	fmt.Printf("User: %+v\n", user)

	// We've now finished all the oauth stuff, so we can store and setup the user

	// See if the user is already registered in the oauth_login table
	var providerID string
	var authorisedUser string
	err = config.DB.QueryRow("SELECT provider_id, authorised_user FROM oauth_login WHERE provider = $1 AND provider_id = $2", providerName, user.UserID).Scan(&providerID, &authorisedUser)
	if err != nil {
		eh.Log.Infow("User not found, inserting new user", "error", err, "user", user)
		fmt.Println("User not found, inserting new user")
		// The user is not registered, so we need to register them
		// Add the user to the user table
		err = config.DB.QueryRow("INSERT INTO \"user\" (main_provider, id, display_name, email, avatar_url) VALUES ($1, $2, $3, $4, $5) RETURNING uuid", user.Provider, user.NickName, user.Name, user.Email, user.AvatarURL).Scan(&authorisedUser)
		if err != nil {
			eh.Log.Errorw("Error inserting new user into DB", "error", err)

			api.EncodeJSONResponse(responder.Response(500, err.Error()), nil, w)
			return
		}

		// Insert into the oauth_login table
		_, err = config.DB.Exec("INSERT INTO oauth_login (provider, provider_id, provider_username, access_token, refresh_token, authorised_user) VALUES ($1, $2, $3, $4, $5, $6)", providerName, user.UserID, user.NickName, user.AccessToken, user.RefreshToken, authorisedUser)
		if err != nil {
			eh.Log.Errorw("Error inserting oauth_login", "error", err)
			api.EncodeJSONResponse(responder.Response(500, err.Error()), nil, w)
			return
		}

		eh.Log.Infow("User registered", "user", authorisedUser)

	} else {
		// The user is already registered, so we just need to update the oauth_login table
		_, err = config.DB.Exec("UPDATE oauth_login SET provider_username = $1, access_token = $2, refresh_token = $3 WHERE provider = $4 AND provider_id = $5", user.NickName, user.AccessToken, user.RefreshToken, providerName, user.UserID)
		if err != nil {
			eh.Log.Errorw("Error updating oauth_login", "error", err, "user", authorisedUser)
			api.EncodeJSONResponse(responder.Response(500, err.Error()), nil, w)
			return
		}
		// Update the user table
		_, err = config.DB.Exec("UPDATE \"user\" SET display_name = $1, email = $2, avatar_url = $3 WHERE uuid = $4", user.Name, user.Email, user.AvatarURL, authorisedUser)
		if err != nil {
			eh.Log.Errorw("Error updating user", "error", err, "user", authorisedUser)
			api.EncodeJSONResponse(responder.Response(500, err.Error()), nil, w)
			return
		}
		eh.Log.Infow("User updated", "reason", "new login, updated details from oauth", "user", authorisedUser)
	}

	theUser, err := myuser.GetByUUID(authorisedUser)
	if err != nil {
		eh.Log.Errorw("Error getting user", "error", err, "user", theUser, "uuid", authorisedUser)
		api.EncodeJSONResponse(responder.Response(500, err.Error()), nil, w)
		return
	}
	eh.Log.Infow("User registered", "user", theUser)

	// Store which user to authorise in the session
	err = StoreInSession(cookie.Value, "authorised_user", authorisedUser)
	if err != nil {
		eh.Log.Errorw("Error storing authorised user in session", "error", err, "user", theUser, "uuid", authorisedUser)
		api.EncodeJSONResponse(responder.Response(500, err.Error()), nil, w)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<html><body>Your session is now authorised. You can now close this window and return to what you were doing.</body> <script> setTimeout(function(){window.close()},5000);</script></html>"))
}

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

// NewUserShimController creates a new shim controller for handling requests in the default API
func NewUserShimController(s api.UserApiServicer, opts ...api.UserApiOption) api.Router {
	controller := &UserAPIShim{
		service: s,
	}

	return controller
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
func (s *UserAPIService) UserLoginCompleteGet(ctx context.Context, token string, clientName string) (api.ImplResponse, error) {
	// Check the name isn't 2TB long
	if len(clientName) > 2048 {
		eh.Log.Errorw("Client gave a name too long", "name", clientName)
		return responder.Response(400, "Client name too long"), nil
	}

	userid, err := GetFromSession(token, "authorised_user")
	if err != nil {
		// TODO: Separate handling for token not found vs other errors
		if err == sql.ErrNoRows {
			eh.Log.Errorw("No authorised user found in session", "token", token)
			return responder.Response(403, "Invalid login token"), nil
		}

		eh.Log.Errorw("Error getting authorised user from session", "error", err, "token", token)
		return responder.Response(500, err.Error()), nil
	}

	user, err := myuser.GetByUUID(userid)
	if err != nil {
		eh.Log.Errorw("Error getting user for token", "error", err, "user", user, "uuid", userid)
		return responder.Response(500, err.Error()), nil
	}

	// Issue an API key for the user
	apiKey, err := user.IssueAPIKey(clientName)
	if err != nil {
		eh.Log.Errorw("Error issuing API key in exchange for token", "error", err, "user", user, "uuid", userid)
		return responder.Response(500, err.Error()), nil
	}

	// API key issued! Remove the reference to the token from the session
	err = StoreInSession(token, "authorised_user", "")
	if err != nil {
		eh.Log.Errorw("Error removing authorised user from session", "error", err, "user", user, "uuid", userid)
		return responder.Response(500, err.Error()), nil
	}

	resp := api.UserLoginCompleteGet200Response{
		ApiKey:   apiKey,
		LoggedIn: true,
	}

	eh.Log.Infow("User logged in", "user", user)

	return responder.Response(200, resp), nil
}

// UserLoginCallbackProviderGet - Handle the callback from the provider
func (s *UserAPIService) UserLoginCallbackProviderGet(ctx context.Context, providerName string) (api.ImplResponse, error) {
	return api.Response(501, nil), nil
}

// // UserLoginRedirectTokenGet - Redirect the user to the provider
// func (s *UserAPIService) UserLoginRedirectTokenGet(ctx context.Context, token string) (api.ImplResponse, error) {
// 	SetUpGoth()
// 	// TODO - update UserLoginRedirectTokenGet with the required logic for this service method.
// 	providerName, err := GetProviderFromSession(token)
// 	if err != nil {
// 		return responder.Response(500, err.Error()), nil
// 	}

// 	provider, err := goth.GetProvider(providerName)
// 	if err != nil {
// 		return responder.Response(400, err.Error()), nil
// 	}

// 	value, err := GetFromSession(token, providerName)
// 	if err != nil {
// 		return responder.Response(500, err.Error()), nil
// 	}

// 	session, err := provider.UnmarshalSession(value)
// 	if err != nil {
// 		return responder.Response(500, err.Error()), nil
// 	}

// }

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
