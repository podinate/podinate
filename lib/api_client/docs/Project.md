# Project

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | The short name (slug/url) of the project | [optional] 
**Name** | Pointer to **string** | The name of the app | [optional] 
**Image** | Pointer to **string** | The container image to run for this app | [optional] 
**Tag** | Pointer to **string** | The image tag to run for this app | [optional] 

## Methods

### NewProject

`func NewProject() *Project`

NewProject instantiates a new Project object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProjectWithDefaults

`func NewProjectWithDefaults() *Project`

NewProjectWithDefaults instantiates a new Project object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Project) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Project) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Project) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *Project) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *Project) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Project) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Project) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *Project) HasName() bool`

HasName returns a boolean if a field has been set.

### GetImage

`func (o *Project) GetImage() string`

GetImage returns the Image field if non-nil, zero value otherwise.

### GetImageOk

`func (o *Project) GetImageOk() (*string, bool)`

GetImageOk returns a tuple with the Image field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImage

`func (o *Project) SetImage(v string)`

SetImage sets Image field to given value.

### HasImage

`func (o *Project) HasImage() bool`

HasImage returns a boolean if a field has been set.

### GetTag

`func (o *Project) GetTag() string`

GetTag returns the Tag field if non-nil, zero value otherwise.

### GetTagOk

`func (o *Project) GetTagOk() (*string, bool)`

GetTagOk returns a tuple with the Tag field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTag

`func (o *Project) SetTag(v string)`

SetTag sets Tag field to given value.

### HasTag

`func (o *Project) HasTag() bool`

HasTag returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


