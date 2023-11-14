/*
Podinate API

The API for the simple containerisation solution Podinate. Login should be performed over oauth from [auth.podinate.com](https://auth.podinate.com)

API version: 0.0.1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// checks if the UserLoginInitGet200Response type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UserLoginInitGet200Response{}

// UserLoginInitGet200Response struct for UserLoginInitGet200Response
type UserLoginInitGet200Response struct {
	// The URL to redirect the user to
	Url *string `json:"url,omitempty"`
	// The token to use to get the user's actual API key once they have completed the oauth flow
	Token *string `json:"token,omitempty"`
}

// NewUserLoginInitGet200Response instantiates a new UserLoginInitGet200Response object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUserLoginInitGet200Response() *UserLoginInitGet200Response {
	this := UserLoginInitGet200Response{}
	return &this
}

// NewUserLoginInitGet200ResponseWithDefaults instantiates a new UserLoginInitGet200Response object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUserLoginInitGet200ResponseWithDefaults() *UserLoginInitGet200Response {
	this := UserLoginInitGet200Response{}
	return &this
}

// GetUrl returns the Url field value if set, zero value otherwise.
func (o *UserLoginInitGet200Response) GetUrl() string {
	if o == nil || IsNil(o.Url) {
		var ret string
		return ret
	}
	return *o.Url
}

// GetUrlOk returns a tuple with the Url field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserLoginInitGet200Response) GetUrlOk() (*string, bool) {
	if o == nil || IsNil(o.Url) {
		return nil, false
	}
	return o.Url, true
}

// HasUrl returns a boolean if a field has been set.
func (o *UserLoginInitGet200Response) HasUrl() bool {
	if o != nil && !IsNil(o.Url) {
		return true
	}

	return false
}

// SetUrl gets a reference to the given string and assigns it to the Url field.
func (o *UserLoginInitGet200Response) SetUrl(v string) {
	o.Url = &v
}

// GetToken returns the Token field value if set, zero value otherwise.
func (o *UserLoginInitGet200Response) GetToken() string {
	if o == nil || IsNil(o.Token) {
		var ret string
		return ret
	}
	return *o.Token
}

// GetTokenOk returns a tuple with the Token field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserLoginInitGet200Response) GetTokenOk() (*string, bool) {
	if o == nil || IsNil(o.Token) {
		return nil, false
	}
	return o.Token, true
}

// HasToken returns a boolean if a field has been set.
func (o *UserLoginInitGet200Response) HasToken() bool {
	if o != nil && !IsNil(o.Token) {
		return true
	}

	return false
}

// SetToken gets a reference to the given string and assigns it to the Token field.
func (o *UserLoginInitGet200Response) SetToken(v string) {
	o.Token = &v
}

func (o UserLoginInitGet200Response) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UserLoginInitGet200Response) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Url) {
		toSerialize["url"] = o.Url
	}
	if !IsNil(o.Token) {
		toSerialize["token"] = o.Token
	}
	return toSerialize, nil
}

type NullableUserLoginInitGet200Response struct {
	value *UserLoginInitGet200Response
	isSet bool
}

func (v NullableUserLoginInitGet200Response) Get() *UserLoginInitGet200Response {
	return v.value
}

func (v *NullableUserLoginInitGet200Response) Set(val *UserLoginInitGet200Response) {
	v.value = val
	v.isSet = true
}

func (v NullableUserLoginInitGet200Response) IsSet() bool {
	return v.isSet
}

func (v *NullableUserLoginInitGet200Response) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUserLoginInitGet200Response(val *UserLoginInitGet200Response) *NullableUserLoginInitGet200Response {
	return &NullableUserLoginInitGet200Response{value: val, isSet: true}
}

func (v NullableUserLoginInitGet200Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUserLoginInitGet200Response) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

