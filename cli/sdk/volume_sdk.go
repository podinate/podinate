package sdk

import (
	"github.com/johncave/podinate/lib/api_client"
)

type Volume struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Size  int    `json:"size"`
	Class string `json:"class"`
}

type VolumeSlice []Volume

// volumesFromAPI converts an API Volume array to a pod VolumeSlice
func volumesFromAPI(apiVolumes []api_client.Volume) VolumeSlice {
	volumes := make(VolumeSlice, len(apiVolumes))
	for i, apiVolume := range apiVolumes {
		volumes[i] = volumeFromAPI(apiVolume)
	}
	return volumes
}

// volumeFromAPI converts an API Volume to a pod Volume
func volumeFromAPI(apiVolume api_client.Volume) Volume {
	out := Volume{
		Name: apiVolume.Name,
		Path: apiVolume.Path,
		Size: int(apiVolume.Size),
	}

	if apiVolume.Class != nil {
		out.Class = *apiVolume.Class
	}

	return out
}

// volumesToAPI converts a pod VolumeSlice to an API Volume array
func volumesToAPI(volumes VolumeSlice) []api_client.Volume {
	apiVolumes := make([]api_client.Volume, len(volumes))
	for i, volume := range volumes {
		apiVolumes[i] = volumeToAPI(volume)
	}
	return apiVolumes
}

// volumeToAPI converts a pod Volume to an API Volume
func volumeToAPI(volume Volume) api_client.Volume {
	out := api_client.Volume{
		Name: volume.Name,
		Path: volume.Path,
		Size: int32(volume.Size),
	}

	if volume.Class != "" {
		out.Class = &volume.Class
	}

	return out
}
