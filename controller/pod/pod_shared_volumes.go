package pod

import (
	api "github.com/johncave/podinate/controller/go"
)

type SharedVolume struct {
	ID   string
	Path string
}

type SharedVolumeSlice []SharedVolume

// SharedVolumesFromApiMany converts from API shared volumes to internal shared volumes
func SharedVolumesFromApiMany(apiSV []api.PodSharedVolumesInner) *SharedVolumeSlice {
	var out SharedVolumeSlice
	for _, sv := range apiSV {
		out = append(out, SharedVolumeFromAPI(sv))
	}
	return &out
}

func SharedVolumeFromAPI(apiSV api.PodSharedVolumesInner) SharedVolume {
	out := SharedVolume{
		ID:   apiSV.VolumeId,
		Path: apiSV.Path,
	}
	return out
}

// SharedVolumeToApiMany translates a SharedVolumeSlice to an api.PodSharedItemsInner slice
func SharedVolumesToApiMany(svs SharedVolumeSlice) []api.PodSharedVolumesInner {
	var out []api.PodSharedVolumesInner
	for _, sv := range svs {
		out = append(out, SharedVolumeToAPI(sv))
	}
	return out
}

func SharedVolumeToAPI(sv SharedVolume) api.PodSharedVolumesInner {
	out := api.PodSharedVolumesInner{
		VolumeId: sv.ID,
		Path:     sv.Path,
	}
	return out
}
