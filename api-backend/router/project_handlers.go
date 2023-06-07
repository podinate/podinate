package router

import (
	"context"
	"log"
	"net/http"

	account "github.com/johncave/podinate/api-backend/account"
	api "github.com/johncave/podinate/api-backend/go"
	"github.com/johncave/podinate/api-backend/project"
	"github.com/johncave/podinate/api-backend/responder"
)

// Import the default service

type ProjectAPIService struct {
	api.ProjectApiService
}

// NewAPIService creates a new service for handling requests in the default API
func NewProjectAPIService() api.ProjectApiServicer {
	return &ProjectAPIService{}
}

// ProjectGet - Returns a list of projects.
func (s *ProjectAPIService) ProjectGet(ctx context.Context, requestedAccount string, page int32, limit int32) (api.ImplResponse, error) {
	// Check the account exists

	log.Println("Trying to get account by ID")
	theAccount, apiErr := account.GetByID(requestedAccount)
	if apiErr != nil {
		return responder.Response(http.StatusBadRequest, apiErr.Message), nil
	}

	log.Println("Trying to get projects by account")
	projects, apiErr := project.GetByAccount(theAccount, page, limit)
	if apiErr != nil {
		return responder.Response(http.StatusInternalServerError, apiErr.Message), nil
	}
	// Assemble the output
	var out []api.Project
	for _, project := range projects {
		out = append(out, project.ToAPI())
	}
	return responder.Response(http.StatusOK, out), nil

}

func (s *ProjectAPIService) ProjectIdGet(ctx context.Context, id string, requestedAccount string) (api.ImplResponse, error) {
	theAccount, apiErr := account.GetByID(requestedAccount)
	if apiErr.Code != http.StatusOK {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}
	theProject, err := project.GetByID(theAccount, id)
	if err != nil {
		// We can pass this error directly to the API response
		return responder.Response(http.StatusNotFound, err.Message), nil
	}
	return responder.Response(http.StatusOK, theProject.ToAPI()), nil
}

func (s *ProjectAPIService) ProjectIdPatch(ctx context.Context, id string, acc string, proj api.Project) (api.ImplResponse, error) {
	theAccount, apiErr := account.GetByID(acc)
	if apiErr.Code != http.StatusOK {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}
	project, apiErr := project.GetByID(theAccount, id)
	if apiErr != nil {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}
	// Update the project
	apiErr = project.Patch(proj)
	if apiErr.Code != http.StatusOK {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}

	return responder.Response(http.StatusOK, project.ToAPI()), nil
}

func (s *ProjectAPIService) ProjectPost(ctx context.Context, requestedAccount string, newProject api.Project) (api.ImplResponse, error) {
	// Check the account exists
	theAccount, apiErr := account.GetByID(requestedAccount)
	if apiErr.Code != http.StatusOK {
		return responder.Response(http.StatusBadRequest, apiErr.Message), nil
	}

	log.Printf("ProjectPost: %v", newProject)
	created, apiErr := project.Create(newProject, theAccount)
	if apiErr.Code != http.StatusOK {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}

	/// Ignore the Kubes logic for now
	// err := createKubesNamespace("project-" + project.Id)
	// if err != nil {
	// 	out := api.Model500Error{Code: http.StatusInternalServerError, Message: err.Error()}
	// 	return responder.Response(http.StatusInternalServerError, out), err
	// }

	// err = createKubesDeployment("project-"+project.Id, project.Image, project.Tag)
	// if err != nil {
	// 	out := api.Model500Error{Code: http.StatusInternalServerError, Message: err.Error()}
	// 	return responder.Response(http.StatusInternalServerError, out), err
	// }

	return responder.Response(http.StatusCreated, created.ToAPI()), nil
}

func (s *ProjectAPIService) ProjectIdDelete(ctx context.Context, id string, requestedAccount string) (api.ImplResponse, error) {
	theAccount, apiErr := account.GetByID(requestedAccount)
	if apiErr != nil {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}
	theProject, apiErr := project.GetByID(theAccount, id)
	if apiErr != nil {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}
	apiErr = theProject.Delete()
	if apiErr != nil {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}
	return responder.Response(http.StatusAccepted, nil), nil
}
