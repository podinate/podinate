package pod

import (
	"github.com/johncave/podinate/api-backend/config"
	api "github.com/johncave/podinate/api-backend/go"
	lh "github.com/johncave/podinate/api-backend/loghandler"
)

type Volume struct {
	Name      string
	MountPath string
	Size      int
	Class     string
}

type VolumeSlice []Volume

// volumesFromAPI converts an API Volume array to a pod VolumeSlice
func volumesFromAPI(apiVolumes []api.Volume) VolumeSlice {
	volumes := make(VolumeSlice, len(apiVolumes))
	for i, apiVolume := range apiVolumes {
		volumes[i] = volumeFromAPI(apiVolume)
	}
	return volumes
}

// volumeFromAPI converts an API Volume to a pod Volume
func volumeFromAPI(apiVolume api.Volume) Volume {
	out := Volume{
		Name:      apiVolume.Name,
		MountPath: apiVolume.MountPath,
		Size:      int(apiVolume.Size),
		Class:     apiVolume.Class,
	}
	return out
}

// ToAPI converts a pod VolumeSlice to an API Volume array
func (volumes VolumeSlice) ToAPI() []api.Volume {
	apiVolumes := make([]api.Volume, len(volumes))
	for i, volume := range volumes {
		apiVolumes[i] = volume.ToAPI()
	}
	//lh.Log.Debug("Converted volumes to API", "volumes", volumes, "apiVolumes", apiVolumes)
	return apiVolumes
}

// ToAPI converts a pod Volume to an API Volume
func (v *Volume) ToAPI() api.Volume {
	out := api.Volume{
		Name:      v.Name,
		MountPath: v.MountPath,
		Size:      int32(v.Size),
		Class:     v.Class,
	}
	return out
}

// loadVolumes loads the volumes from the database
func (p *Pod) loadVolumes() error {
	rows, err := config.DB.Query("SELECT name, mount_path, size FROM pod_volumes WHERE pod_uuid = $1", p.Uuid)
	if err != nil {
		lh.Log.Errorw("Error getting pod volumes", "error", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var v Volume
		err := rows.Scan(&v.Name, &v.MountPath, &v.Size)
		if err != nil {
			return err
		}
		lh.Log.Debugw("Loaded volume", "volume", v)
		p.Volumes = append(p.Volumes, v)
	}
	return nil
}

// volumeExists checks if a volume exists in the pod
func (p *Pod) volumeExists(name string) (bool, error) {
	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM pod_volumes WHERE pod_uuid = $1 AND name = $2", p.Uuid, name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
