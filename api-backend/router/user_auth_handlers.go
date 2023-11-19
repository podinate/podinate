package router

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	eh "github.com/johncave/podinate/api-backend/errorhandler"
	api "github.com/johncave/podinate/api-backend/go"
	"github.com/johncave/podinate/api-backend/responder"
	myuser "github.com/johncave/podinate/api-backend/user"
	"github.com/markbates/goth"
)

type UserAPIShim struct {
	service      api.UserApiServicer
	errorHandler api.ErrorHandler
}

// NewUserShimController creates a new shim controller for handling requests in the default API
func NewUserShimController(s api.UserApiServicer, opts ...api.UserApiOption) api.Router {
	controller := &UserAPIShim{
		service: s,
	}

	return controller
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

// UserLoginRedirectTokenGet - Redirect the user to the provider.
// This is a shim so that the user can be redirected to the provider from this code,
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
// This is shimmed so that Goth can have full access to all the variables
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
	theUser, err := myuser.RegisterUser(user)
	if err != nil {
		eh.Log.Errorw("Error registering user", "error", err, "user", theUser, "uuid", theUser.UUID)
		api.EncodeJSONResponse(responder.Response(500, err.Error()), nil, w)
		return
	}

	// Store which user to authorise in the session
	err = StoreInSession(cookie.Value, "authorised_user", theUser.UUID)
	if err != nil {
		eh.Log.Errorw("Error storing authorised user in session", "error", err, "user", theUser, "uuid", theUser.UUID)
		api.EncodeJSONResponse(responder.Response(500, err.Error()), nil, w)
		return
	}

	// Delete the login-token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "login-token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1,
	})

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<html><body>Your session is now authorised. You can now close this window and return to what you were doing.</body> <script> setTimeout(function(){window.close()},5000);</script></html>"))
}
