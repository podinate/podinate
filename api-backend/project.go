package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	api "github.com/johncave/podinate/api-backend/go"
)

// ProjectGet - Returns a list of projects.
func (s *APIService) ProjectGet(ctx context.Context, account string) (api.ImplResponse, error) {
	// TODO - update ProjectGet with the required logic for this service method.
	out, _ := json.Marshal(api.Project{Id: "podinate-blog", Name: "Podinate Blog", Image: "wordpress", Tag: "latest"})
	return api.Response(http.StatusNotImplemented, string(out)), nil
}

func (s *APIService) ProjectIdGet(ctx context.Context, id string, account string) (api.ImplResponse, error) {
	// TODO - update ProjectIdGet with the required logic for this service method.

	out, _ := json.Marshal(api.Project{Id: id, Name: "test", Image: "test", Tag: "latest"})
	return api.Response(http.StatusNotImplemented, out), errors.New("ProjectIdGet method not implemented")
}

func (s *APIService) ProjectIdPatch(ctx context.Context, id string, account string) (api.ImplResponse, error) {
	// TODO - update ProjectIdPatch with the required logic for this service method.

	out, _ := json.Marshal(api.Project{Id: id, Name: "test", Image: "test", Tag: "latest"})
	return api.Response(http.StatusNotImplemented, out), errors.New("ProjectIdPatch method not implemented")
}

func (s *APIService) ProjectPost(ctx context.Context, account string, project api.Project) (api.ImplResponse, error) {
	// TODO - update ProjectPost with the required logic for this service method.

	log.Printf("ProjectPost: %v", project)
	err := createKubesNamespace("project-" + project.Id)
	if err != nil {
		out := api.Model500Error{Code: http.StatusInternalServerError, Message: err.Error()}
		return api.Response(http.StatusInternalServerError, out), err
	}

	err = createKubesDeployment("project-"+project.Id, project.Image, project.Tag)
	if err != nil {
		out := api.Model500Error{Code: http.StatusInternalServerError, Message: err.Error()}
		return api.Response(http.StatusInternalServerError, out), err
	}

	return api.Response(http.StatusCreated, project), nil
}
