package router

import (
	"context"
	"log"
	"net/http"

	api "github.com/johncave/podinate/api-backend/go"
	pod "github.com/johncave/podinate/api-backend/pod"
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
func (s *PodAPIService) ProjectProjectIdPodGet(ctx context.Context, id string, acc string, page int32, limit int32) (api.ImplResponse, error) {
	// TODO - Implement ProjectProjectIdPodGet
	return responder.Response(http.StatusNotImplemented, "ProjectProjectIdPodGet needs to be implemented!"), nil
}

// ProjectProjectIdPodPodIdGet - Returns a pod for a project.
func (s *PodAPIService) ProjectProjectIdPodPodIdGet(ctx context.Context, id string, acc string, podId string) (api.ImplResponse, error) {
	// TODO - Implement ProjectProjectIdPodPodIdGet
	return responder.Response(http.StatusNotImplemented, "ProjectProjectIdPodPodIdGet needs to be implemented!"), nil
}

// ProjectProjectIdPodPodIdPatch - Updates a pod for a project.
func (s *PodAPIService) ProjectProjectIdPodPodIdPatch(ctx context.Context, id string, acc string, podId string, pod api.Pod) (api.ImplResponse, error) {
	// TODO - Implement ProjectProjectIdPodPodIdPatch
	return responder.Response(http.StatusNotImplemented, "ProjectProjectIdPodPodIdPatch needs to be implemented!"), nil
}

// ProjectProjectIdPodPost - Creates a pod for a project.
func (s *PodAPIService) ProjectProjectIdPodPost(ctx context.Context, projectId string, acc string, requestedPod api.Pod) (api.ImplResponse, error) {

	thepod, err := pod.Create(requestedPod)
	if err != nil {
		return responder.Response(http.StatusInternalServerError, err.Error()), nil
	}
	log.Printf("%+v", thepod)

	// TODO - Implement ProjectProjectIdPodPost
	return responder.Response(http.StatusNotImplemented, "ProjectProjectIdPodPost needs to be implemented!"), nil
}

// ProjectProjectIdPodPodIdDelete - Deletes a pod for a project.
func (s *PodAPIService) ProjectProjectIdPodPodIdDelete(ctx context.Context, id string, acc string, podId string) (api.ImplResponse, error) {
	// TODO - Implement ProjectProjectIdPodPodIdDelete
	return responder.Response(http.StatusNotImplemented, "ProjectProjectIdPodPodIdDelete needs to be implemented!"), nil
}
