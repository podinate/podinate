# ProjectProjectIdSharedVolumesGet200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Items** | Pointer to [**[]SharedVolume**](SharedVolume.md) |  | [optional] 
**Total** | **int32** | The total number of shared volumes | 
**Page** | **int32** | The current page number | 
**Limit** | **int32** | The number of items per page | 

## Methods

### NewProjectProjectIdSharedVolumesGet200Response

`func NewProjectProjectIdSharedVolumesGet200Response(total int32, page int32, limit int32, ) *ProjectProjectIdSharedVolumesGet200Response`

NewProjectProjectIdSharedVolumesGet200Response instantiates a new ProjectProjectIdSharedVolumesGet200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProjectProjectIdSharedVolumesGet200ResponseWithDefaults

`func NewProjectProjectIdSharedVolumesGet200ResponseWithDefaults() *ProjectProjectIdSharedVolumesGet200Response`

NewProjectProjectIdSharedVolumesGet200ResponseWithDefaults instantiates a new ProjectProjectIdSharedVolumesGet200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetItems

`func (o *ProjectProjectIdSharedVolumesGet200Response) GetItems() []SharedVolume`

GetItems returns the Items field if non-nil, zero value otherwise.

### GetItemsOk

`func (o *ProjectProjectIdSharedVolumesGet200Response) GetItemsOk() (*[]SharedVolume, bool)`

GetItemsOk returns a tuple with the Items field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetItems

`func (o *ProjectProjectIdSharedVolumesGet200Response) SetItems(v []SharedVolume)`

SetItems sets Items field to given value.

### HasItems

`func (o *ProjectProjectIdSharedVolumesGet200Response) HasItems() bool`

HasItems returns a boolean if a field has been set.

### GetTotal

`func (o *ProjectProjectIdSharedVolumesGet200Response) GetTotal() int32`

GetTotal returns the Total field if non-nil, zero value otherwise.

### GetTotalOk

`func (o *ProjectProjectIdSharedVolumesGet200Response) GetTotalOk() (*int32, bool)`

GetTotalOk returns a tuple with the Total field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotal

`func (o *ProjectProjectIdSharedVolumesGet200Response) SetTotal(v int32)`

SetTotal sets Total field to given value.


### GetPage

`func (o *ProjectProjectIdSharedVolumesGet200Response) GetPage() int32`

GetPage returns the Page field if non-nil, zero value otherwise.

### GetPageOk

`func (o *ProjectProjectIdSharedVolumesGet200Response) GetPageOk() (*int32, bool)`

GetPageOk returns a tuple with the Page field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPage

`func (o *ProjectProjectIdSharedVolumesGet200Response) SetPage(v int32)`

SetPage sets Page field to given value.


### GetLimit

`func (o *ProjectProjectIdSharedVolumesGet200Response) GetLimit() int32`

GetLimit returns the Limit field if non-nil, zero value otherwise.

### GetLimitOk

`func (o *ProjectProjectIdSharedVolumesGet200Response) GetLimitOk() (*int32, bool)`

GetLimitOk returns a tuple with the Limit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLimit

`func (o *ProjectProjectIdSharedVolumesGet200Response) SetLimit(v int32)`

SetLimit sets Limit field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


