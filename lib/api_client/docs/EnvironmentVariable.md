# EnvironmentVariable

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Key** | **string** | The key of the environment variable | 
**Value** | **string** | The value of the environment variable | 
**Secret** | Pointer to **bool** | Whether the value of the environment variable is a secret or not. If it is a secret, it will not be returned in the API response. | [optional] 

## Methods

### NewEnvironmentVariable

`func NewEnvironmentVariable(key string, value string, ) *EnvironmentVariable`

NewEnvironmentVariable instantiates a new EnvironmentVariable object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEnvironmentVariableWithDefaults

`func NewEnvironmentVariableWithDefaults() *EnvironmentVariable`

NewEnvironmentVariableWithDefaults instantiates a new EnvironmentVariable object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKey

`func (o *EnvironmentVariable) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *EnvironmentVariable) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *EnvironmentVariable) SetKey(v string)`

SetKey sets Key field to given value.


### GetValue

`func (o *EnvironmentVariable) GetValue() string`

GetValue returns the Value field if non-nil, zero value otherwise.

### GetValueOk

`func (o *EnvironmentVariable) GetValueOk() (*string, bool)`

GetValueOk returns a tuple with the Value field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValue

`func (o *EnvironmentVariable) SetValue(v string)`

SetValue sets Value field to given value.


### GetSecret

`func (o *EnvironmentVariable) GetSecret() bool`

GetSecret returns the Secret field if non-nil, zero value otherwise.

### GetSecretOk

`func (o *EnvironmentVariable) GetSecretOk() (*bool, bool)`

GetSecretOk returns a tuple with the Secret field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecret

`func (o *EnvironmentVariable) SetSecret(v bool)`

SetSecret sets Secret field to given value.

### HasSecret

`func (o *EnvironmentVariable) HasSecret() bool`

HasSecret returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


