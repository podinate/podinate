/*
 * Podinate API
 *
 * The API for the simple containerisation solution Podinate. Login should be performed over oauth from [auth.podinate.com](https://auth.podinate.com)
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// PodApiController binds http requests to an api service and writes the service results to the http response
type PodApiController struct {
	service PodApiServicer
	errorHandler ErrorHandler
}

// PodApiOption for how the controller is set up.
type PodApiOption func(*PodApiController)

// WithPodApiErrorHandler inject ErrorHandler into controller
func WithPodApiErrorHandler(h ErrorHandler) PodApiOption {
	return func(c *PodApiController) {
		c.errorHandler = h
	}
}

// NewPodApiController creates a default api controller
func NewPodApiController(s PodApiServicer, opts ...PodApiOption) Router {
	controller := &PodApiController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the PodApiController
func (c *PodApiController) Routes() Routes {
	return Routes{ 
		{
			"ProjectProjectIdPodGet",
			strings.ToUpper("Get"),
			"/v0/project/{project_id}/pod",
			c.ProjectProjectIdPodGet,
		},
		{
			"ProjectProjectIdPodPodIdDelete",
			strings.ToUpper("Delete"),
			"/v0/project/{project_id}/pod/{pod_id}",
			c.ProjectProjectIdPodPodIdDelete,
		},
		{
			"ProjectProjectIdPodPodIdGet",
			strings.ToUpper("Get"),
			"/v0/project/{project_id}/pod/{pod_id}",
			c.ProjectProjectIdPodPodIdGet,
		},
		{
			"ProjectProjectIdPodPodIdPatch",
			strings.ToUpper("Patch"),
			"/v0/project/{project_id}/pod/{pod_id}",
			c.ProjectProjectIdPodPodIdPatch,
		},
		{
			"ProjectProjectIdPodPost",
			strings.ToUpper("Post"),
			"/v0/project/{project_id}/pod",
			c.ProjectProjectIdPodPost,
		},
	}
}

// ProjectProjectIdPodGet - Get a list of pods for a given project
func (c *PodApiController) ProjectProjectIdPodGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query := r.URL.Query()
	projectIdParam := params["project_id"]
	accountParam := r.Header.Get("account")
	pageParam, err := parseInt32Parameter(query.Get("page"), false)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	limitParam, err := parseInt32Parameter(query.Get("limit"), false)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	result, err := c.service.ProjectProjectIdPodGet(r.Context(), projectIdParam, accountParam, pageParam, limitParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// ProjectProjectIdPodPodIdDelete - Delete a pod
func (c *PodApiController) ProjectProjectIdPodPodIdDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectIdParam := params["project_id"]
	podIdParam := params["pod_id"]
	accountParam := r.Header.Get("account")
	result, err := c.service.ProjectProjectIdPodPodIdDelete(r.Context(), projectIdParam, podIdParam, accountParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// ProjectProjectIdPodPodIdGet - Get a pod by ID
func (c *PodApiController) ProjectProjectIdPodPodIdGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectIdParam := params["project_id"]
	podIdParam := params["pod_id"]
	accountParam := r.Header.Get("account")
	result, err := c.service.ProjectProjectIdPodPodIdGet(r.Context(), projectIdParam, podIdParam, accountParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// ProjectProjectIdPodPodIdPatch - Update a pod
func (c *PodApiController) ProjectProjectIdPodPodIdPatch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectIdParam := params["project_id"]
	podIdParam := params["pod_id"]
	accountParam := r.Header.Get("account")
	podParam := Pod{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&podParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertPodRequired(podParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.ProjectProjectIdPodPodIdPatch(r.Context(), projectIdParam, podIdParam, accountParam, podParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// ProjectProjectIdPodPost - Create a new pod
func (c *PodApiController) ProjectProjectIdPodPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectIdParam := params["project_id"]
	accountParam := r.Header.Get("account")
	podParam := Pod{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&podParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertPodRequired(podParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.ProjectProjectIdPodPost(r.Context(), projectIdParam, accountParam, podParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}
