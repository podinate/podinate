package apiclient

import (
	"github.com/johncave/podinate/lib/api_client"
)

type Volume struct {
	Name      string `json:"name"`
	MountPath string `json:"mount_path"`
	Size      int    `json:"size"`
	Class     string `json:"class"`
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
		Name:      apiVolume.Name,
		MountPath: apiVolume.MountPath,
		Size:      int(apiVolume.Size),
	}

	if apiVolume.Class != nil {
		out.Class = *apiVolume.Class
	}

	return out
}
