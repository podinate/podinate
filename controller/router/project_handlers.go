package router

import (
	"context"
	"log"
	"net/http"

	account "github.com/johncave/podinate/controller/account"
	api "github.com/johncave/podinate/controller/go"
	"github.com/johncave/podinate/controller/iam"
	"github.com/johncave/podinate/controller/project"
	"github.com/johncave/podinate/controller/responder"
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
	//var out []api.Project
	var items []api.ProjectGet200ResponseItemsInner
	for _, theProject := range projects {
		if iam.RequestorCan(ctx, &theAccount, theProject, project.ActionView) {
			//out = append(out, theProject.ToAPI())
			items = append(items, api.ProjectGet200ResponseItemsInner{
				Id:         theProject.ID,
				Name:       theProject.Name,
				ResourceId: theProject.GetResourceID(),
			})
		}
	}

	return responder.Response(http.StatusOK, api.ProjectGet200Response{Items: items, Total: int32(len(items)), Page: page, Limit: limit}), nil

}

func (s *ProjectAPIService) ProjectIdGet(ctx context.Context, id string, requestedAccount string) (api.ImplResponse, error) {
	theAccount, apiErr := account.GetByID(requestedAccount)
	if apiErr != nil {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}
	theProject, err := project.GetByID(theAccount, id)
	if err != nil {
		// We can pass this error directly to the API response
		return responder.Response(http.StatusNotFound, err.Message), nil
	}
	if !iam.RequestorCan(ctx, &theAccount, theProject, project.ActionView) {
		return responder.Response(http.StatusForbidden, "You do not have permission to view this project"), nil
	}
	return responder.Response(http.StatusOK, theProject.ToAPI()), nil
}

func (s *ProjectAPIService) ProjectIdPut(ctx context.Context, id string, acc string, proj api.Project) (api.ImplResponse, error) {
	theAccount, apiErr := account.GetByID(acc)
	if apiErr != nil {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}
	theProject, apiErr := project.GetByID(theAccount, id)
	if apiErr != nil {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}

	// Check if the user has permission to do the action on the resource
	if !iam.RequestorCan(ctx, &theAccount, theProject, project.ActionUpdate) {
		return responder.Response(http.StatusForbidden, "You do not have permission to update this project"), nil
	}

	// Update the project
	apiErr = theProject.Patch(proj)
	if apiErr != nil {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}

	return responder.Response(http.StatusOK, theProject.ToAPI()), nil
}

func (s *ProjectAPIService) ProjectPost(ctx context.Context, requestedAccount string, newProject api.Project) (api.ImplResponse, error) {
	// Check the account exists
	//lh.Info(ctx, "Trying to get account by ID", "account", requestedAccount, "project", newProject)
	theAccount, apiErr := account.GetByID(requestedAccount)
	if apiErr != nil {
		return responder.Response(http.StatusBadRequest, apiErr.Message), nil
	}

	log.Printf("ProjectPost: %v", newProject)

	res := iam.NewResource(theAccount.GetResourceID() + "/project:" + newProject.Id)
	if !iam.RequestorCan(ctx, &theAccount, res, project.ActionCreate) {
		return responder.Response(http.StatusForbidden, "You do not have permission to create this project"), nil
	}
	created, apiErr := project.Create(ctx, newProject, theAccount)
	if apiErr != nil {
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

	// Check if the user has permission to do the action on the resource
	if !iam.RequestorCan(ctx, &theAccount, theProject, project.ActionDelete) {
		return responder.Response(http.StatusForbidden, "You do not have permission to delete this project"), nil
	}

	apiErr = theProject.Delete()
	if apiErr != nil {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}
	return responder.Response(http.StatusAccepted, nil), nil
}
