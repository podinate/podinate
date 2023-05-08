package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	account "github.com/johncave/podinate/api-backend/account"
	api "github.com/johncave/podinate/api-backend/go"
	"github.com/johncave/podinate/api-backend/project"
)

// ProjectGet - Returns a list of projects.
func (s *APIService) ProjectGet(ctx context.Context, requestedAccount string, page int32, limit int32) (api.ImplResponse, error) {
	// Check the account exists

	theAccount := account.Account{}
	err := theAccount.GetByID(requestedAccount)
	if err != nil {
		return api.Response(http.StatusBadRequest, err.Error()), nil
	}

	projects, apiErr := theAccount.GetProjects(page, limit)
	if apiErr.Code != http.StatusOK {
		return api.Response(http.StatusInternalServerError, err.Error()), nil
	}
	// Assemble the output
	var out []api.Project
	for _, project := range projects {
		out = append(out, project.ToAPI())
	}
	return api.Response(http.StatusOK, out), nil

}

func (s *APIService) ProjectIdGet(ctx context.Context, id string, requestedAccount string) (api.ImplResponse, error) {
	theProject := project.Project{}
	theAccount := account.Account{}
	theAccount.GetByID(requestedAccount)
	err := theProject.GetByID(theAccount, id)
	if err != nil {
		// We can pass this error directly to the API response
		return api.Response(http.StatusNotFound, err.Error()), nil
	}
	return api.Response(http.StatusOK, theProject.ToAPI()), nil
}

func (s *APIService) ProjectIdPatch(ctx context.Context, id string, account string) (api.ImplResponse, error) {
	// TODO - update ProjectIdPatch with the required logic for this service method.

	out, _ := json.Marshal(api.Project{Id: id, Name: "test", Image: "test", Tag: "latest"})
	return api.Response(http.StatusNotImplemented, out), errors.New("ProjectIdPatch method not implemented")
}

func (s *APIService) ProjectPost(ctx context.Context, requestedAccount string, newProject api.Project) (api.ImplResponse, error) {
	// Check the account exists
	theAccount := account.Account{}
	err := theAccount.GetByID(requestedAccount)
	if err != nil {
		return api.Response(http.StatusBadRequest, err.Error()), nil
	}

	log.Printf("ProjectPost: %v", newProject)
	var created project.Project
	apiErr := created.Create(newProject, theAccount)
	if err != nil {
		return api.Response(apiErr.Code, apiErr.Error()), nil
	}

	/// Ignore the Kubes logic for now
	// err := createKubesNamespace("project-" + project.Id)
	// if err != nil {
	// 	out := api.Model500Error{Code: http.StatusInternalServerError, Message: err.Error()}
	// 	return api.Response(http.StatusInternalServerError, out), err
	// }

	// err = createKubesDeployment("project-"+project.Id, project.Image, project.Tag)
	// if err != nil {
	// 	out := api.Model500Error{Code: http.StatusInternalServerError, Message: err.Error()}
	// 	return api.Response(http.StatusInternalServerError, out), err
	// }

	return api.Response(http.StatusCreated, created.ToAPI()), nil
}
