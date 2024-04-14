package router

import (
	"context"
	"net/http"

	api "github.com/johncave/podinate/controller/go"
	"github.com/johncave/podinate/controller/responder"
	"github.com/johncave/podinate/controller/shared_volume"
)

type SharedVolumeApiService struct {
	api.SharedVolumeApiService
}

func NewSharedVolumeAPIService() api.SharedVolumeApiServicer {
	return &SharedVolumeApiService{}
}

// ProjectProjectIdSharedVolumesGet - Get a list of shared volumes for a given project
func (s *SharedVolumeApiService) ProjectProjectIdSharedVolumesGet(ctx context.Context, projectId string, accountId string, page int32, limit int32) (api.ImplResponse, error) {
	project, apiErr := getProject(ctx, accountId, projectId)
	if apiErr != nil {
		return apiErr.ToResponse(), nil
	}

	sharedVolumes, apiErr := shared_volume.GetByProject(ctx, project)
	if apiErr != nil {
		return apiErr.ToResponse(), nil
	}
	var items []api.SharedVolume
	for _, sharedVolume := range sharedVolumes {
		items = append(items, sharedVolume.ToAPI())
	}

	out := api.ProjectProjectIdSharedVolumesGet200Response{
		Items: items,
		Total: int32(len(items)),
		Page:  page,
	}

	return responder.Response(http.StatusOK, out), nil
}

// ProjectProjectIdSharedVolumesPost - Create a new shared volume
func (s *SharedVolumeApiService) ProjectProjectIdSharedVolumesPost(ctx context.Context, projectId string, account string, in api.SharedVolume) (api.ImplResponse, error) {
	project, err := getProject(ctx, account, projectId)
	if err != nil {
		return err.ToResponse(), nil
	}

	// IAM calls are inside the Create function
	sharedVolume, apiErr := shared_volume.Create(ctx, project, in)
	if apiErr != nil {
		return apiErr.ToResponse(), nil
	}

	return responder.Response(http.StatusCreated, sharedVolume.ToAPI()), nil
}

// ProjectProjectIdSharedVolumesVolumeIdDelete - Delete a shared volume
func (s *SharedVolumeApiService) ProjectProjectIdSharedVolumesVolumeIdDelete(ctx context.Context, projectId string, volumeId string, account string) (api.ImplResponse, error) {
	project, err := getProject(ctx, account, projectId)
	if err != nil {
		return err.ToResponse(), nil
	}

	vol, err := shared_volume.GetByID(ctx, project, volumeId)
	if err != nil {
		return err.ToResponse(), nil
	}

	err = vol.Delete(ctx)
	if err != nil {
		return err.ToResponse(), nil
	}

	return responder.Response(http.StatusNoContent, nil), nil
}

// ProjectProjectIdSharedVolumesVolumeIdGet - Get a shared volume by ID
func (s *SharedVolumeApiService) ProjectProjectIdSharedVolumesVolumeIdGet(ctx context.Context, projectId string, volumeId string, account string) (api.ImplResponse, error) {

	// Get the project
	project, err := getProject(ctx, account, projectId)
	if err != nil {
		return err.ToResponse(), nil
	}

	// Get the shared volume
	sharedVolume, err := shared_volume.GetByID(ctx, project, volumeId)
	if err != nil {
		return err.ToResponse(), nil
	}

	return responder.Response(http.StatusOK, sharedVolume.ToAPI()), nil
}

// ProjectProjectIdSharedVolumesVolumeIdPut - Update a shared volume&#39;s spec
func (s *SharedVolumeApiService) ProjectProjectIdSharedVolumesVolumeIdPut(ctx context.Context, projectId string, volumeId string, account string, in api.SharedVolume) (api.ImplResponse, error) {

	// Get the project
	project, err := getProject(ctx, account, projectId)
	if err != nil {
		return err.ToResponse(), nil
	}

	// Get the shared volume
	sharedVolume, err := shared_volume.GetByID(ctx, project, volumeId)
	if err != nil {
		return err.ToResponse(), nil
	}

	// Update the shared volume
	err = sharedVolume.Update(ctx, in)
	if err != nil {
		return err.ToResponse(), nil
	}

	return responder.Response(http.StatusOK, sharedVolume.ToAPI()), nil
}
