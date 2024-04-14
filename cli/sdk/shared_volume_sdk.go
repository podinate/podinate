package sdk

import (
	"context"

	"github.com/johncave/podinate/lib/api_client"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type SharedVolume struct {
	ID      string
	Name    *string
	Project *Project
	Size    int
	Class   *string
}

type SharedVolumeAttachment struct {
	ID   string
	Path string
}

type SharedVolumeAttachmentSlice []SharedVolumeAttachment

// GetSharedVolumeByID returns a shared volume by ID from the given project
func (p *Project) GetSharedVolumeByID(id string) (*SharedVolume, *SDKError) {
	resp, r, err := C.SharedVolumeApi.ProjectProjectIdSharedVolumesVolumeIdGet(context.Background(), p.ID, id).Account(viper.GetString("account")).Execute()
	if err := handleAPIError(r, err); err != nil {
		return nil, err
	}
	return getSharedVolumeFromApi(p, resp), nil
}

// GetSharedVolumes returns all shared volumes from the given project
func (p *Project) GetSharedVolumes() ([]*SharedVolume, *SDKError) {
	resp, r, err := C.SharedVolumeApi.ProjectProjectIdSharedVolumesGet(context.Background(), p.ID).Account(viper.GetString("account")).Execute()
	if err := handleAPIError(r, err); err != nil {
		return nil, err
	}

	var volumes []*SharedVolume

	for _, v := range resp.Items {
		volumes = append(volumes, getSharedVolumeFromApi(p, &v))
	}

	return volumes, nil
}

// CreateSharedVolume creates a new shared volume in the given project
func (p *Project) CreateSharedVolume(in SharedVolume) (*SharedVolume, error) {
	zap.S().Debugw("Creating shared volume", "project", p.ID, "shared_volume", in)
	resp, r, err := C.SharedVolumeApi.ProjectProjectIdSharedVolumesPost(context.Background(), p.ID).Account(viper.GetString("account")).SharedVolume(*sharedVolumeToAPI(&in)).Execute()
	if err := handleAPIError(r, err); err != nil {
		return nil, err
	}
	return getSharedVolumeFromApi(p, resp), nil
}

// Update updates a shared volume
func (v *SharedVolume) Update(in *SharedVolume) error {
	resp, r, err := C.SharedVolumeApi.ProjectProjectIdSharedVolumesVolumeIdPut(context.Background(), v.Project.ID, v.ID).Account(viper.GetString("account")).SharedVolume(*sharedVolumeToAPI(v)).Execute()
	if err := handleAPIError(r, err); err != nil {
		return err
	}

	v = getSharedVolumeFromApi(v.Project, resp)
	return nil
}

// Delete deletes a shared volume
func (v *SharedVolume) Delete() error {
	r, err := C.SharedVolumeApi.ProjectProjectIdSharedVolumesVolumeIdDelete(context.Background(), v.Project.ID, v.ID).Account(viper.GetString("account")).Execute()
	return handleAPIError(r, err)
}

func getSharedVolumeFromApi(p *Project, in *api_client.SharedVolume) *SharedVolume {
	return &SharedVolume{
		ID:      in.Id,
		Name:    in.Name,
		Project: p,
		Size:    int(in.Size),
		Class:   in.Class,
	}
}

func sharedVolumeToAPI(in *SharedVolume) *api_client.SharedVolume {
	return &api_client.SharedVolume{
		Id:    in.ID,
		Name:  in.Name,
		Size:  int32(in.Size),
		Class: in.Class,
	}
}

// sharedVolumeAttachmentsFromAPI converts an API SharedVolumeAttachment array to a SDK SharedVolumeAttachmentSlice
func sharedVolumeAttachmentsFromAPI(apiAttachments []api_client.PodSharedVolumesInner) SharedVolumeAttachmentSlice {
	attachments := make(SharedVolumeAttachmentSlice, len(apiAttachments))
	for i, apiAttachment := range apiAttachments {
		attachments[i] = sharedVolumeAttachmentFromAPI(apiAttachment)
	}
	return attachments
}

// sharedVolumeAttachmentFromAPI converts an API SharedVolumeAttachment to a SDK SharedVolumeAttachment
func sharedVolumeAttachmentFromAPI(apiAttachment api_client.PodSharedVolumesInner) SharedVolumeAttachment {
	return SharedVolumeAttachment{
		ID:   apiAttachment.VolumeId,
		Path: apiAttachment.Path,
	}
}

// sharedVolumeAttachmentsToAPI converts a SDK SharedVolumeAttachmentSlice to an API SharedVolumeAttachment array
func sharedVolumeAttachmentsToAPI(attachments SharedVolumeAttachmentSlice) []api_client.PodSharedVolumesInner {
	apiAttachments := make([]api_client.PodSharedVolumesInner, len(attachments))
	for i, attachment := range attachments {
		apiAttachments[i] = sharedVolumeAttachmentToAPI(attachment)
	}
	return apiAttachments
}

// sharedVolumeAttachmentToAPI converts a SDK SharedVolumeAttachment to an API SharedVolumeAttachment
func sharedVolumeAttachmentToAPI(attachment SharedVolumeAttachment) api_client.PodSharedVolumesInner {
	return api_client.PodSharedVolumesInner{
		VolumeId: attachment.ID,
		Path:     attachment.Path,
	}
}
