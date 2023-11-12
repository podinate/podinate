# UserLoginCompleteGet200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**LoggedIn** | Pointer to **bool** | Whether the user is logged in or not | [optional] 
**ApiKey** | Pointer to **string** | The user&#39;s API key, if they are logged in | [optional] 

## Methods

### NewUserLoginCompleteGet200Response

`func NewUserLoginCompleteGet200Response() *UserLoginCompleteGet200Response`

NewUserLoginCompleteGet200Response instantiates a new UserLoginCompleteGet200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUserLoginCompleteGet200ResponseWithDefaults

`func NewUserLoginCompleteGet200ResponseWithDefaults() *UserLoginCompleteGet200Response`

NewUserLoginCompleteGet200ResponseWithDefaults instantiates a new UserLoginCompleteGet200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLoggedIn

`func (o *UserLoginCompleteGet200Response) GetLoggedIn() bool`

GetLoggedIn returns the LoggedIn field if non-nil, zero value otherwise.

### GetLoggedInOk

`func (o *UserLoginCompleteGet200Response) GetLoggedInOk() (*bool, bool)`

GetLoggedInOk returns a tuple with the LoggedIn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLoggedIn

`func (o *UserLoginCompleteGet200Response) SetLoggedIn(v bool)`

SetLoggedIn sets LoggedIn field to given value.

### HasLoggedIn

`func (o *UserLoginCompleteGet200Response) HasLoggedIn() bool`

HasLoggedIn returns a boolean if a field has been set.

### GetApiKey

`func (o *UserLoginCompleteGet200Response) GetApiKey() string`

GetApiKey returns the ApiKey field if non-nil, zero value otherwise.

### GetApiKeyOk

`func (o *UserLoginCompleteGet200Response) GetApiKeyOk() (*string, bool)`

GetApiKeyOk returns a tuple with the ApiKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApiKey

`func (o *UserLoginCompleteGet200Response) SetApiKey(v string)`

SetApiKey sets ApiKey field to given value.

### HasApiKey

`func (o *UserLoginCompleteGet200Response) HasApiKey() bool`

HasApiKey returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


