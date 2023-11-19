package router

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/johncave/podinate/api-backend/config"
	api "github.com/johncave/podinate/api-backend/go"
	"github.com/johncave/podinate/api-backend/iam"
	lh "github.com/johncave/podinate/api-backend/loghandler"
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

	r := api.NewRouter(ProjectApiController, AccountApiController, PodApiController, UserShimController, UserApiController)
	r.Use(authMiddleware)
	r.Use(loggingMiddleware)

	return r

}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Check if the request is to the login endpoints
		currentRoute := mux.CurrentRoute(r).GetName()
		if currentRoute == "UserLoginInitiateGet" || currentRoute == "UserLoginCallbackProviderGet" || currentRoute == "UserLoginCompleteGet" || currentRoute == "UserLoginRedirectTokenGet" {
			next.ServeHTTP(w, r)
			return
		}

		// Get the token from the header
		keyin := r.Header.Get("Authorization")
		if keyin == "" {
			lh.Log.Errorw("No API Key provided", "request", r)
			http.Error(w, "No API Key provided", http.StatusUnauthorized)
			return

		}

		requestor, err := iam.GetFromAuthorizationHeader(keyin)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Set the user in the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, iam.ContextKey("requestor"), requestor)

		// Add the new context to the request
		r = r.Clone(ctx)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := iam.GetFromContext(r.Context())
		requestor := "anonymous"

		if u != nil {
			requestor = u.GetRID()
		}

		requestID := lh.NewRequestID()

		ctx := r.Context()
		ctx = context.WithValue(ctx, lh.ContextKey("request-id"), requestID)

		// Add the new context to the request
		r = r.Clone(ctx)

		lh.Log.Infow("request", "request_id", requestID, "method", r.Method, "url", r.URL, "remote", r.Header.Get("x-forwarded-for"), "user-agent", r.UserAgent(), "referer", r.Referer(), "requestor", requestor)

		next.ServeHTTP(w, r)
	})
}

// MakeRequest - Make a request to the API, useful for testing
func MakeRequest(method, url string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	config.Init()
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
