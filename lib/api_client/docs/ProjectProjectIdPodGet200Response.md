# ProjectProjectIdPodGet200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Items** | Pointer to [**[]ProjectProjectIdPodGet200ResponseItemsInner**](ProjectProjectIdPodGet200ResponseItemsInner.md) |  | [optional] 
**Total** | **int32** | The total number of pods | 
**Page** | **int32** | The current page number | 
**Limit** | **int32** | The number of items per page | 

## Methods

### NewProjectProjectIdPodGet200Response

`func NewProjectProjectIdPodGet200Response(total int32, page int32, limit int32, ) *ProjectProjectIdPodGet200Response`

NewProjectProjectIdPodGet200Response instantiates a new ProjectProjectIdPodGet200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProjectProjectIdPodGet200ResponseWithDefaults

`func NewProjectProjectIdPodGet200ResponseWithDefaults() *ProjectProjectIdPodGet200Response`

NewProjectProjectIdPodGet200ResponseWithDefaults instantiates a new ProjectProjectIdPodGet200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetItems

`func (o *ProjectProjectIdPodGet200Response) GetItems() []ProjectProjectIdPodGet200ResponseItemsInner`

GetItems returns the Items field if non-nil, zero value otherwise.

### GetItemsOk

`func (o *ProjectProjectIdPodGet200Response) GetItemsOk() (*[]ProjectProjectIdPodGet200ResponseItemsInner, bool)`

GetItemsOk returns a tuple with the Items field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetItems

`func (o *ProjectProjectIdPodGet200Response) SetItems(v []ProjectProjectIdPodGet200ResponseItemsInner)`

SetItems sets Items field to given value.

### HasItems

`func (o *ProjectProjectIdPodGet200Response) HasItems() bool`

HasItems returns a boolean if a field has been set.

### GetTotal

`func (o *ProjectProjectIdPodGet200Response) GetTotal() int32`

GetTotal returns the Total field if non-nil, zero value otherwise.

### GetTotalOk

`func (o *ProjectProjectIdPodGet200Response) GetTotalOk() (*int32, bool)`

GetTotalOk returns a tuple with the Total field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotal

`func (o *ProjectProjectIdPodGet200Response) SetTotal(v int32)`

SetTotal sets Total field to given value.


### GetPage

`func (o *ProjectProjectIdPodGet200Response) GetPage() int32`

GetPage returns the Page field if non-nil, zero value otherwise.

### GetPageOk

`func (o *ProjectProjectIdPodGet200Response) GetPageOk() (*int32, bool)`

GetPageOk returns a tuple with the Page field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPage

`func (o *ProjectProjectIdPodGet200Response) SetPage(v int32)`

SetPage sets Page field to given value.


### GetLimit

`func (o *ProjectProjectIdPodGet200Response) GetLimit() int32`

GetLimit returns the Limit field if non-nil, zero value otherwise.

### GetLimitOk

`func (o *ProjectProjectIdPodGet200Response) GetLimitOk() (*int32, bool)`

GetLimitOk returns a tuple with the Limit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLimit

`func (o *ProjectProjectIdPodGet200Response) SetLimit(v int32)`

SetLimit sets Limit field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


