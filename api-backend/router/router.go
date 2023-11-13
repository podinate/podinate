package router

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/johncave/podinate/api-backend/config"
	api "github.com/johncave/podinate/api-backend/go"
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

	return api.NewRouter(ProjectApiController, AccountApiController, PodApiController, UserApiController)

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
