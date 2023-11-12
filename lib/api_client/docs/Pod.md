# Pod

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | The short name (slug/url) of the pod | [optional] 
**Name** | Pointer to **string** | The name of the pod | [optional] 
**Image** | Pointer to **string** | The container image to run for this pod | [optional] 
**Tag** | Pointer to **string** | The image tag to run for this pod | [optional] 
**Status** | Pointer to **string** | The current status of the pod | [optional] 
**CreatedAt** | Pointer to **string** | The date and time the pod was created | [optional] 

## Methods

### NewPod

`func NewPod() *Pod`

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

### HasId

`func (o *Pod) HasId() bool`

HasId returns a boolean if a field has been set.

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

### HasName

`func (o *Pod) HasName() bool`

HasName returns a boolean if a field has been set.

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

### HasImage

`func (o *Pod) HasImage() bool`

HasImage returns a boolean if a field has been set.

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

### HasTag

`func (o *Pod) HasTag() bool`

HasTag returns a boolean if a field has been set.

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


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


