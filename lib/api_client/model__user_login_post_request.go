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

// checks if the UserLoginPostRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UserLoginPostRequest{}

// UserLoginPostRequest struct for UserLoginPostRequest
type UserLoginPostRequest struct {
	// The user's email address
	Username string `json:"username"`
	// The user's password
	Password string `json:"password"`
	// The client name to use for the login
	Client *string `json:"client,omitempty"`
}

// NewUserLoginPostRequest instantiates a new UserLoginPostRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUserLoginPostRequest(username string, password string) *UserLoginPostRequest {
	this := UserLoginPostRequest{}
	this.Username = username
	this.Password = password
	return &this
}

// NewUserLoginPostRequestWithDefaults instantiates a new UserLoginPostRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUserLoginPostRequestWithDefaults() *UserLoginPostRequest {
	this := UserLoginPostRequest{}
	return &this
}

// GetUsername returns the Username field value
func (o *UserLoginPostRequest) GetUsername() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Username
}

// GetUsernameOk returns a tuple with the Username field value
// and a boolean to check if the value has been set.
func (o *UserLoginPostRequest) GetUsernameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Username, true
}

// SetUsername sets field value
func (o *UserLoginPostRequest) SetUsername(v string) {
	o.Username = v
}

// GetPassword returns the Password field value
func (o *UserLoginPostRequest) GetPassword() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Password
}

// GetPasswordOk returns a tuple with the Password field value
// and a boolean to check if the value has been set.
func (o *UserLoginPostRequest) GetPasswordOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Password, true
}

// SetPassword sets field value
func (o *UserLoginPostRequest) SetPassword(v string) {
	o.Password = v
}

// GetClient returns the Client field value if set, zero value otherwise.
func (o *UserLoginPostRequest) GetClient() string {
	if o == nil || IsNil(o.Client) {
		var ret string
		return ret
	}
	return *o.Client
}

// GetClientOk returns a tuple with the Client field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserLoginPostRequest) GetClientOk() (*string, bool) {
	if o == nil || IsNil(o.Client) {
		return nil, false
	}
	return o.Client, true
}

// HasClient returns a boolean if a field has been set.
func (o *UserLoginPostRequest) HasClient() bool {
	if o != nil && !IsNil(o.Client) {
		return true
	}

	return false
}

// SetClient gets a reference to the given string and assigns it to the Client field.
func (o *UserLoginPostRequest) SetClient(v string) {
	o.Client = &v
}

func (o UserLoginPostRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UserLoginPostRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["username"] = o.Username
	toSerialize["password"] = o.Password
	if !IsNil(o.Client) {
		toSerialize["client"] = o.Client
	}
	return toSerialize, nil
}

type NullableUserLoginPostRequest struct {
	value *UserLoginPostRequest
	isSet bool
}

func (v NullableUserLoginPostRequest) Get() *UserLoginPostRequest {
	return v.value
}

func (v *NullableUserLoginPostRequest) Set(val *UserLoginPostRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableUserLoginPostRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableUserLoginPostRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUserLoginPostRequest(val *UserLoginPostRequest) *NullableUserLoginPostRequest {
	return &NullableUserLoginPostRequest{value: val, isSet: true}
}

func (v NullableUserLoginPostRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUserLoginPostRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
