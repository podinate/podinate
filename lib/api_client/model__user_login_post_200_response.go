/*
Podinate API

The API for the simple containerisation solution Podinate.

API version: 0.0.1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api_client

import (
	"encoding/json"
)

// checks if the UserLoginPost200Response type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UserLoginPost200Response{}

// UserLoginPost200Response struct for UserLoginPost200Response
type UserLoginPost200Response struct {
	// Whether the user is logged in or not
	LoggedIn *bool `json:"logged_in,omitempty"`
	// The user's API key, if they are logged in
	ApiKey *string `json:"api_key,omitempty"`
}

// NewUserLoginPost200Response instantiates a new UserLoginPost200Response object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUserLoginPost200Response() *UserLoginPost200Response {
	this := UserLoginPost200Response{}
	return &this
}

// NewUserLoginPost200ResponseWithDefaults instantiates a new UserLoginPost200Response object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUserLoginPost200ResponseWithDefaults() *UserLoginPost200Response {
	this := UserLoginPost200Response{}
	return &this
}

// GetLoggedIn returns the LoggedIn field value if set, zero value otherwise.
func (o *UserLoginPost200Response) GetLoggedIn() bool {
	if o == nil || IsNil(o.LoggedIn) {
		var ret bool
		return ret
	}
	return *o.LoggedIn
}

// GetLoggedInOk returns a tuple with the LoggedIn field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserLoginPost200Response) GetLoggedInOk() (*bool, bool) {
	if o == nil || IsNil(o.LoggedIn) {
		return nil, false
	}
	return o.LoggedIn, true
}

// HasLoggedIn returns a boolean if a field has been set.
func (o *UserLoginPost200Response) HasLoggedIn() bool {
	if o != nil && !IsNil(o.LoggedIn) {
		return true
	}

	return false
}

// SetLoggedIn gets a reference to the given bool and assigns it to the LoggedIn field.
func (o *UserLoginPost200Response) SetLoggedIn(v bool) {
	o.LoggedIn = &v
}

// GetApiKey returns the ApiKey field value if set, zero value otherwise.
func (o *UserLoginPost200Response) GetApiKey() string {
	if o == nil || IsNil(o.ApiKey) {
		var ret string
		return ret
	}
	return *o.ApiKey
}

// GetApiKeyOk returns a tuple with the ApiKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserLoginPost200Response) GetApiKeyOk() (*string, bool) {
	if o == nil || IsNil(o.ApiKey) {
		return nil, false
	}
	return o.ApiKey, true
}

// HasApiKey returns a boolean if a field has been set.
func (o *UserLoginPost200Response) HasApiKey() bool {
	if o != nil && !IsNil(o.ApiKey) {
		return true
	}

	return false
}

// SetApiKey gets a reference to the given string and assigns it to the ApiKey field.
func (o *UserLoginPost200Response) SetApiKey(v string) {
	o.ApiKey = &v
}

func (o UserLoginPost200Response) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UserLoginPost200Response) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.LoggedIn) {
		toSerialize["logged_in"] = o.LoggedIn
	}
	if !IsNil(o.ApiKey) {
		toSerialize["api_key"] = o.ApiKey
	}
	return toSerialize, nil
}

type NullableUserLoginPost200Response struct {
	value *UserLoginPost200Response
	isSet bool
}

func (v NullableUserLoginPost200Response) Get() *UserLoginPost200Response {
	return v.value
}

func (v *NullableUserLoginPost200Response) Set(val *UserLoginPost200Response) {
	v.value = val
	v.isSet = true
}

func (v NullableUserLoginPost200Response) IsSet() bool {
	return v.isSet
}

func (v *NullableUserLoginPost200Response) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUserLoginPost200Response(val *UserLoginPost200Response) *NullableUserLoginPost200Response {
	return &NullableUserLoginPost200Response{value: val, isSet: true}
}

func (v NullableUserLoginPost200Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUserLoginPost200Response) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
