# UserLoginInitGet200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Url** | Pointer to **string** | The URL to redirect the user to | [optional] 
**Token** | Pointer to **string** | The token to use to get the user&#39;s actual API key once they have completed the oauth flow | [optional] 

## Methods

### NewUserLoginInitGet200Response

`func NewUserLoginInitGet200Response() *UserLoginInitGet200Response`

NewUserLoginInitGet200Response instantiates a new UserLoginInitGet200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUserLoginInitGet200ResponseWithDefaults

`func NewUserLoginInitGet200ResponseWithDefaults() *UserLoginInitGet200Response`

NewUserLoginInitGet200ResponseWithDefaults instantiates a new UserLoginInitGet200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUrl

`func (o *UserLoginInitGet200Response) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *UserLoginInitGet200Response) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *UserLoginInitGet200Response) SetUrl(v string)`

SetUrl sets Url field to given value.

### HasUrl

`func (o *UserLoginInitGet200Response) HasUrl() bool`

HasUrl returns a boolean if a field has been set.

### GetToken

`func (o *UserLoginInitGet200Response) GetToken() string`

GetToken returns the Token field if non-nil, zero value otherwise.

### GetTokenOk

`func (o *UserLoginInitGet200Response) GetTokenOk() (*string, bool)`

GetTokenOk returns a tuple with the Token field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetToken

`func (o *UserLoginInitGet200Response) SetToken(v string)`

SetToken sets Token field to given value.

### HasToken

`func (o *UserLoginInitGet200Response) HasToken() bool`

HasToken returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


