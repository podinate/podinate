/*
 * Podinate API
 *
 * The API for the simple containerisation solution Podinate
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"context"
	"net/http"
)



// DefaultApiRouter defines the required methods for binding the api requests to a responses for the DefaultApi
// The DefaultApiRouter implementation should parse necessary information from the http request,
// pass the data to a DefaultApiServicer to perform the required actions, then write the service results to the http response.
type DefaultApiRouter interface { 
	ProjectGet(http.ResponseWriter, *http.Request)
	ProjectIdGet(http.ResponseWriter, *http.Request)
	ProjectIdPatch(http.ResponseWriter, *http.Request)
	ProjectPost(http.ResponseWriter, *http.Request)
}


// DefaultApiServicer defines the api actions for the DefaultApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type DefaultApiServicer interface { 
	ProjectGet(context.Context, string) (ImplResponse, error)
	ProjectIdGet(context.Context, string, string) (ImplResponse, error)
	ProjectIdPatch(context.Context, string, string) (ImplResponse, error)
	ProjectPost(context.Context, string, App) (ImplResponse, error)
}
