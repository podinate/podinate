# SharedVolume

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | The short name (slug/url) of the shared volume | 
**Name** | Pointer to **string** | The name of the shared volume | [optional] 
**Size** | **int32** | The size of the shared volume in GB | 
**Class** | Pointer to **string** | The class of the shared volume, for example \&quot;standard\&quot; or \&quot;premium\&quot; | [optional] 
**ResourceId** | Pointer to **string** | The global Resource ID of the shared volume | [optional] 

## Methods

### NewSharedVolume

`func NewSharedVolume(id string, size int32, ) *SharedVolume`

NewSharedVolume instantiates a new SharedVolume object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSharedVolumeWithDefaults

`func NewSharedVolumeWithDefaults() *SharedVolume`

NewSharedVolumeWithDefaults instantiates a new SharedVolume object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *SharedVolume) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *SharedVolume) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *SharedVolume) SetId(v string)`

SetId sets Id field to given value.


### GetName

`func (o *SharedVolume) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SharedVolume) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SharedVolume) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *SharedVolume) HasName() bool`

HasName returns a boolean if a field has been set.

### GetSize

`func (o *SharedVolume) GetSize() int32`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *SharedVolume) GetSizeOk() (*int32, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *SharedVolume) SetSize(v int32)`

SetSize sets Size field to given value.


### GetClass

`func (o *SharedVolume) GetClass() string`

GetClass returns the Class field if non-nil, zero value otherwise.

### GetClassOk

`func (o *SharedVolume) GetClassOk() (*string, bool)`

GetClassOk returns a tuple with the Class field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClass

`func (o *SharedVolume) SetClass(v string)`

SetClass sets Class field to given value.

### HasClass

`func (o *SharedVolume) HasClass() bool`

HasClass returns a boolean if a field has been set.

### GetResourceId

`func (o *SharedVolume) GetResourceId() string`

GetResourceId returns the ResourceId field if non-nil, zero value otherwise.

### GetResourceIdOk

`func (o *SharedVolume) GetResourceIdOk() (*string, bool)`

GetResourceIdOk returns a tuple with the ResourceId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceId

`func (o *SharedVolume) SetResourceId(v string)`

SetResourceId sets ResourceId field to given value.

### HasResourceId

`func (o *SharedVolume) HasResourceId() bool`

HasResourceId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


