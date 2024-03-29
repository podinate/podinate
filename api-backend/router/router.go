package router

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/johncave/podinate/api-backend/account"
	"github.com/johncave/podinate/api-backend/apierror"
	api "github.com/johncave/podinate/api-backend/go"
	"github.com/johncave/podinate/api-backend/iam"
	lh "github.com/johncave/podinate/api-backend/loghandler"
	"github.com/johncave/podinate/api-backend/project"
)

// GetRouter - Get the router for the API
func GetRouter() *mux.Router {
	ProjectAPIService := NewProjectAPIService()
	ProjectApiController := api.NewProjectApiController(ProjectAPIService)

	AccountAPIService := NewAccountAPIService()
	AccountApiController := api.NewAccountApiController(AccountAPIService)

	PodAPIService := NewPodAPIService()
	PodApiController := api.NewPodApiController(PodAPIService)

	UserAPIService := NewUserAPIService()
	UserApiController := api.NewUserApiController(UserAPIService)

	UserShimController := NewUserShimController(NewUserAPIService())
	PodShimController := NewPodShimController(NewPodAPIService())

	r := api.NewRouter(ProjectApiController, AccountApiController, PodShimController, PodApiController, UserShimController, UserApiController)
	r.Use(requestIDMiddleware)
	r.Use(authMiddleware)
	r.Use(loggingMiddleware)

	return r

}

// requestIDMiddleware - Add a request ID to the context
// This middleware runs first to make sure all others have access to the request ID
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := lh.NewRequestID()

		ctx := r.Context()
		ctx = context.WithValue(ctx, lh.ContextKey("request-id"), requestID)

		// Add the new context to the request
		r = r.Clone(ctx)

		next.ServeHTTP(w, r)
	})
}

// authMiddleware - Add the user to the context
// Checks if the path is one of the login paths, if so, skip authentication
// Otherwise, check the token and add the user to the context
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Check if the request is to the login endpoints
		currentRoute := mux.CurrentRoute(r).GetName()
		if currentRoute == "UserLoginInitiateGet" ||
			currentRoute == "UserLoginCallbackProviderGet" ||
			currentRoute == "UserLoginCompleteGet" ||
			currentRoute == "UserLoginRedirectTokenGet" ||
			currentRoute == "UserLoginPost" {

			next.ServeHTTP(w, r)
			return
		}

		r, err := iam.AddRequestorToRequest(r)

		if err != nil {
			//out := api.EncodeJSONResponse(responder.Response(http.StatusUnauthorized, "Invalid Key"), nil, w)

			out, err := json.Marshal(api.Error{
				Code:    http.StatusUnauthorized,
				Message: "Invalid API Key",
			})
			if err != nil {
				http.Error(w, "Invalid API Key", http.StatusUnauthorized)
				return
			}
			http.Error(w, string(out), http.StatusUnauthorized)

			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware - Log basic http information about the request
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := iam.GetFromContext(r.Context())
		requestor := "anonymous"

		if u != nil {
			requestor = u.GetResourceID()
		}

		requestID := lh.GetRequestID(r.Context())
		w.Header().Add("X-Request-ID", string(requestID))
		// body, err := r.GetBody()
		// if err != nil {
		// 	body = nil
		// 	lh.Error(r.Context(), "Error getting body", "error", err)
		// }
		lh.Log.Infow("request", "request_id", requestID, "method", r.Method, "url", r.URL, "remote", r.Header.Get("x-forwarded-for"), "user-agent", r.UserAgent(), "referer", r.Referer(), "requestor", requestor)

		next.ServeHTTP(w, r)
	})
}

// MakeRequest - Make a request to the API, useful for testing
func MakeRequest(method, url string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	//config.Init()
	log.Println("Configurinated")

	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	writer := httptest.NewRecorder()

	router := GetRouter()
	router.ServeHTTP(writer, request)
	return writer
}

func getAccountAndProject(accountID string, projectID string) (*account.Account, *project.Project, *apierror.ApiError) {
	// Get the account by ID
	theAccount, apiErr := account.GetByID(accountID)
	if apiErr != nil {
		return nil, nil, apiErr
	}

	// Get the project this pod lives in by ID
	theProject, apiErr := project.GetByID(theAccount, projectID)
	if apiErr != nil {
		return nil, nil, apiErr
	}

	return &theAccount, &theProject, nil
}
