# PodSharedVolumesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**VolumeId** | **string** | The short name (slug/url) of the shared volume | 
**Path** | **string** | The path to mount the shared volume at | 

## Methods

### NewPodSharedVolumesInner

`func NewPodSharedVolumesInner(volumeId string, path string, ) *PodSharedVolumesInner`

NewPodSharedVolumesInner instantiates a new PodSharedVolumesInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPodSharedVolumesInnerWithDefaults

`func NewPodSharedVolumesInnerWithDefaults() *PodSharedVolumesInner`

NewPodSharedVolumesInnerWithDefaults instantiates a new PodSharedVolumesInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetVolumeId

`func (o *PodSharedVolumesInner) GetVolumeId() string`

GetVolumeId returns the VolumeId field if non-nil, zero value otherwise.

### GetVolumeIdOk

`func (o *PodSharedVolumesInner) GetVolumeIdOk() (*string, bool)`

GetVolumeIdOk returns a tuple with the VolumeId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVolumeId

`func (o *PodSharedVolumesInner) SetVolumeId(v string)`

SetVolumeId sets VolumeId field to given value.


### GetPath

`func (o *PodSharedVolumesInner) GetPath() string`

GetPath returns the Path field if non-nil, zero value otherwise.

### GetPathOk

`func (o *PodSharedVolumesInner) GetPathOk() (*string, bool)`

GetPathOk returns a tuple with the Path field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPath

`func (o *PodSharedVolumesInner) SetPath(v string)`

SetPath sets Path field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


