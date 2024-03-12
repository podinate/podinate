# Pod

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | The short name (slug/url) of the pod | 
**Name** | **string** | The name of the pod | 
**Image** | **string** | The container image to run for this pod | 
**Tag** | **string** | The image tag to run for this pod | 
**Volumes** | Pointer to [**[]Volume**](Volume.md) | The storage volumes attached to this pod | [optional] 
**Environment** | Pointer to [**[]EnvironmentVariable**](EnvironmentVariable.md) | The environment variables to pass to the pod | [optional] 
**Services** | Pointer to [**[]Service**](Service.md) | The services to expose for this pod | [optional] 
**Status** | Pointer to **string** | The current status of the pod | [optional] 
**CreatedAt** | Pointer to **string** | The date and time the pod was created | [optional] 
**ResourceId** | Pointer to **string** | The global Resource ID of the pod | [optional] 

## Methods

### NewPod

`func NewPod(id string, name string, image string, tag string, ) *Pod`

NewPod instantiates a new Pod object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPodWithDefaults

`func NewPodWithDefaults() *Pod`

NewPodWithDefaults instantiates a new Pod object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Pod) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Pod) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Pod) SetId(v string)`

SetId sets Id field to given value.


### GetName

`func (o *Pod) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Pod) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Pod) SetName(v string)`

SetName sets Name field to given value.


### GetImage

`func (o *Pod) GetImage() string`

GetImage returns the Image field if non-nil, zero value otherwise.

### GetImageOk

`func (o *Pod) GetImageOk() (*string, bool)`

GetImageOk returns a tuple with the Image field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImage

`func (o *Pod) SetImage(v string)`

SetImage sets Image field to given value.


### GetTag

`func (o *Pod) GetTag() string`

GetTag returns the Tag field if non-nil, zero value otherwise.

### GetTagOk

`func (o *Pod) GetTagOk() (*string, bool)`

GetTagOk returns a tuple with the Tag field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTag

`func (o *Pod) SetTag(v string)`

SetTag sets Tag field to given value.


### GetVolumes

`func (o *Pod) GetVolumes() []Volume`

GetVolumes returns the Volumes field if non-nil, zero value otherwise.

### GetVolumesOk

`func (o *Pod) GetVolumesOk() (*[]Volume, bool)`

GetVolumesOk returns a tuple with the Volumes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVolumes

`func (o *Pod) SetVolumes(v []Volume)`

SetVolumes sets Volumes field to given value.

### HasVolumes

`func (o *Pod) HasVolumes() bool`

HasVolumes returns a boolean if a field has been set.

### GetEnvironment

`func (o *Pod) GetEnvironment() []EnvironmentVariable`

GetEnvironment returns the Environment field if non-nil, zero value otherwise.

### GetEnvironmentOk

`func (o *Pod) GetEnvironmentOk() (*[]EnvironmentVariable, bool)`

GetEnvironmentOk returns a tuple with the Environment field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnvironment

`func (o *Pod) SetEnvironment(v []EnvironmentVariable)`

SetEnvironment sets Environment field to given value.

### HasEnvironment

`func (o *Pod) HasEnvironment() bool`

HasEnvironment returns a boolean if a field has been set.

### GetServices

`func (o *Pod) GetServices() []Service`

GetServices returns the Services field if non-nil, zero value otherwise.

### GetServicesOk

`func (o *Pod) GetServicesOk() (*[]Service, bool)`

GetServicesOk returns a tuple with the Services field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServices

`func (o *Pod) SetServices(v []Service)`

SetServices sets Services field to given value.

### HasServices

`func (o *Pod) HasServices() bool`

HasServices returns a boolean if a field has been set.

### GetStatus

`func (o *Pod) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *Pod) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *Pod) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *Pod) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetCreatedAt

`func (o *Pod) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *Pod) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *Pod) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *Pod) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetResourceId

`func (o *Pod) GetResourceId() string`

GetResourceId returns the ResourceId field if non-nil, zero value otherwise.

### GetResourceIdOk

`func (o *Pod) GetResourceIdOk() (*string, bool)`

GetResourceIdOk returns a tuple with the ResourceId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceId

`func (o *Pod) SetResourceId(v string)`

SetResourceId sets ResourceId field to given value.

### HasResourceId

`func (o *Pod) HasResourceId() bool`

HasResourceId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


