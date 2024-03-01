# Volume

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | The name of the volume | 
**Size** | **int32** | The size of the volume in GB | 
**MountPath** | **string** | The path to mount the volume at | 
**Class** | Pointer to **string** | The class of the volume, for example \&quot;standard\&quot; or \&quot;premium\&quot; | [optional] 

## Methods

### NewVolume

`func NewVolume(name string, size int32, mountPath string, ) *Volume`

NewVolume instantiates a new Volume object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVolumeWithDefaults

`func NewVolumeWithDefaults() *Volume`

NewVolumeWithDefaults instantiates a new Volume object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *Volume) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Volume) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Volume) SetName(v string)`

SetName sets Name field to given value.


### GetSize

`func (o *Volume) GetSize() int32`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *Volume) GetSizeOk() (*int32, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *Volume) SetSize(v int32)`

SetSize sets Size field to given value.


### GetMountPath

`func (o *Volume) GetMountPath() string`

GetMountPath returns the MountPath field if non-nil, zero value otherwise.

### GetMountPathOk

`func (o *Volume) GetMountPathOk() (*string, bool)`

GetMountPathOk returns a tuple with the MountPath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMountPath

`func (o *Volume) SetMountPath(v string)`

SetMountPath sets MountPath field to given value.


### GetClass

`func (o *Volume) GetClass() string`

GetClass returns the Class field if non-nil, zero value otherwise.

### GetClassOk

`func (o *Volume) GetClassOk() (*string, bool)`

GetClassOk returns a tuple with the Class field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClass

`func (o *Volume) SetClass(v string)`

SetClass sets Class field to given value.

### HasClass

`func (o *Volume) HasClass() bool`

HasClass returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


