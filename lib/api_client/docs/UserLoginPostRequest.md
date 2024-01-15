# UserLoginPostRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Username** | **string** | The user&#39;s email address | 
**Password** | **string** | The user&#39;s password | 
**Client** | Pointer to **string** | The client name to use for the login | [optional] 

## Methods

### NewUserLoginPostRequest

`func NewUserLoginPostRequest(username string, password string, ) *UserLoginPostRequest`

NewUserLoginPostRequest instantiates a new UserLoginPostRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUserLoginPostRequestWithDefaults

`func NewUserLoginPostRequestWithDefaults() *UserLoginPostRequest`

NewUserLoginPostRequestWithDefaults instantiates a new UserLoginPostRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUsername

`func (o *UserLoginPostRequest) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *UserLoginPostRequest) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *UserLoginPostRequest) SetUsername(v string)`

SetUsername sets Username field to given value.


### GetPassword

`func (o *UserLoginPostRequest) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *UserLoginPostRequest) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *UserLoginPostRequest) SetPassword(v string)`

SetPassword sets Password field to given value.


### GetClient

`func (o *UserLoginPostRequest) GetClient() string`

GetClient returns the Client field if non-nil, zero value otherwise.

### GetClientOk

`func (o *UserLoginPostRequest) GetClientOk() (*string, bool)`

GetClientOk returns a tuple with the Client field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClient

`func (o *UserLoginPostRequest) SetClient(v string)`

SetClient sets Client field to given value.

### HasClient

`func (o *UserLoginPostRequest) HasClient() bool`

HasClient returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


