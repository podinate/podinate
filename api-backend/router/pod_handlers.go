package router

import (
	"context"
	"log"
	"net/http"

	"github.com/johncave/podinate/api-backend/account"
	api "github.com/johncave/podinate/api-backend/go"
	pod "github.com/johncave/podinate/api-backend/pod"
	"github.com/johncave/podinate/api-backend/project"
	"github.com/johncave/podinate/api-backend/responder"
)

// PodAPIService holds all the handlers for the Pod API
type PodAPIService struct {
	api.PodApiService
}

// NewPodAPIService creates a new service for handling requests in the Pod API
func NewPodAPIService() api.PodApiServicer {
	return &PodAPIService{}
}

// ProjectProjectIdPodGet - Returns a list of pods for a project.
func (s *PodAPIService) ProjectProjectIdPodGet(ctx context.Context, projectID string, accountID string, page int32, limit int32) (api.ImplResponse, error) {

	log.Printf("%s %s %d %d", projectID, accountID, page, limit)

	// Get the account
	theAccount, apiErr := account.GetByID(accountID)
	if apiErr != nil {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}

	// Get the project to list the pods of
	theProject, apiErr := project.GetByID(theAccount, projectID)
	if apiErr != nil {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}

	// Get the pods for the project
	pods, apiErr := pod.GetByProject(theProject, page, limit)
	if apiErr != nil {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}

	// Convert the pods to the API format
	apiPods := make([]api.Pod, 0)
	for _, p := range pods {
		apiPods = append(apiPods, p.ToAPI())
	}

	return responder.Response(http.StatusOK, apiPods), nil
}

// ProjectProjectIdPodPodIdGet - Returns a pod for a project.
func (s *PodAPIService) ProjectProjectIdPodPodIdGet(ctx context.Context, projectID string, podID string, accountID string) (api.ImplResponse, error) {

	log.Printf("%s %s %s", projectID, podID, accountID)

	// Get the account by ID
	theAccount, apiErr := account.GetByID(accountID)
	if apiErr != nil {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}

	// Get the project this pod lives in by ID
	theProject, apiErr := project.GetByID(theAccount, projectID)
	if apiErr != nil {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}

	pod, apiErr := pod.GetByID(theProject, podID)
	if apiErr != nil {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}
	return responder.Response(http.StatusOK, pod.ToAPI()), nil
}

// ProjectProjectIdPodPodIdPatch - Updates a pod for a project.
func (s *PodAPIService) ProjectProjectIdPodPodIdPatch(ctx context.Context, id string, acc string, podId string, pod api.Pod) (api.ImplResponse, error) {
	// TODO - Implement ProjectProjectIdPodPodIdPatch
	return responder.Response(http.StatusNotImplemented, "ProjectProjectIdPodPodIdPatch needs to be implemented!"), nil
}

// ProjectProjectIdPodPost - Creates a pod for a project.
func (s *PodAPIService) ProjectProjectIdPodPost(ctx context.Context, projectId string, accountID string, requestedPod api.Pod) (api.ImplResponse, error) {
	// Get the account by ID
	theAccount, apiErr := account.GetByID(accountID)
	if apiErr != nil {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}

	// Get the project this pod lives in by ID
	theProject, apiErr := project.GetByID(theAccount, projectId)
	if apiErr != nil {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}

	thepod, err := pod.Create(theAccount, theProject, requestedPod)
	if err != nil {
		return responder.Response(err.Code, err.Message), nil
	}
	log.Printf("%+v", thepod)

	return responder.Response(http.StatusOK, thepod.ToAPI()), nil
}

// ProjectProjectIdPodPodIdDelete - Deletes a pod for a project.
func (s *PodAPIService) ProjectProjectIdPodPodIdDelete(ctx context.Context, projectID string, accountID string, podID string) (api.ImplResponse, error) {

	log.Printf("%s %s %s", projectID, podID, accountID)

	// Get the account by ID
	theAccount, apiErr := account.GetByID(accountID)
	if apiErr != nil {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}

	// Get the project this pod lives in by ID
	theProject, apiErr := project.GetByID(theAccount, projectID)
	if apiErr != nil {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}

	// Get the pod by ID
	thePod, apiErr := pod.GetByID(theProject, podID)
	if apiErr != nil {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}

	// Delete the pod
	apiErr = pod.Delete(theProject, thePod)
	if apiErr != nil {
		return responder.Response(apiErr.Code, apiErr.Message), nil
	}

	return responder.Response(http.StatusOK, nil), nil
}
