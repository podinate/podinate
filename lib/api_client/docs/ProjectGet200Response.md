# ProjectGet200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Items** | Pointer to [**[]ProjectGet200ResponseItemsInner**](ProjectGet200ResponseItemsInner.md) |  | [optional] 
**Total** | Pointer to **int32** | The total number of projects | [optional] 
**Page** | Pointer to **int32** | The current page number | [optional] 
**Limit** | Pointer to **int32** | The number of items per page | [optional] 

## Methods

### NewProjectGet200Response

`func NewProjectGet200Response() *ProjectGet200Response`

NewProjectGet200Response instantiates a new ProjectGet200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProjectGet200ResponseWithDefaults

`func NewProjectGet200ResponseWithDefaults() *ProjectGet200Response`

NewProjectGet200ResponseWithDefaults instantiates a new ProjectGet200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetItems

`func (o *ProjectGet200Response) GetItems() []ProjectGet200ResponseItemsInner`

GetItems returns the Items field if non-nil, zero value otherwise.

### GetItemsOk

`func (o *ProjectGet200Response) GetItemsOk() (*[]ProjectGet200ResponseItemsInner, bool)`

GetItemsOk returns a tuple with the Items field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetItems

`func (o *ProjectGet200Response) SetItems(v []ProjectGet200ResponseItemsInner)`

SetItems sets Items field to given value.

### HasItems

`func (o *ProjectGet200Response) HasItems() bool`

HasItems returns a boolean if a field has been set.

### GetTotal

`func (o *ProjectGet200Response) GetTotal() int32`

GetTotal returns the Total field if non-nil, zero value otherwise.

### GetTotalOk

`func (o *ProjectGet200Response) GetTotalOk() (*int32, bool)`

GetTotalOk returns a tuple with the Total field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotal

`func (o *ProjectGet200Response) SetTotal(v int32)`

SetTotal sets Total field to given value.

### HasTotal

`func (o *ProjectGet200Response) HasTotal() bool`

HasTotal returns a boolean if a field has been set.

### GetPage

`func (o *ProjectGet200Response) GetPage() int32`

GetPage returns the Page field if non-nil, zero value otherwise.

### GetPageOk

`func (o *ProjectGet200Response) GetPageOk() (*int32, bool)`

GetPageOk returns a tuple with the Page field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPage

`func (o *ProjectGet200Response) SetPage(v int32)`

SetPage sets Page field to given value.

### HasPage

`func (o *ProjectGet200Response) HasPage() bool`

HasPage returns a boolean if a field has been set.

### GetLimit

`func (o *ProjectGet200Response) GetLimit() int32`

GetLimit returns the Limit field if non-nil, zero value otherwise.

### GetLimitOk

`func (o *ProjectGet200Response) GetLimitOk() (*int32, bool)`

GetLimitOk returns a tuple with the Limit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLimit

`func (o *ProjectGet200Response) SetLimit(v int32)`

SetLimit sets Limit field to given value.

### HasLimit

`func (o *ProjectGet200Response) HasLimit() bool`

HasLimit returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


