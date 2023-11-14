/*
Podinate API

The API for the simple containerisation solution Podinate. Login should be performed over oauth from [auth.podinate.com](https://auth.podinate.com)

API version: 0.0.1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
	"fmt"
)

// ProjectProjectIdPodGet200ResponseInner - struct for ProjectProjectIdPodGet200ResponseInner
type ProjectProjectIdPodGet200ResponseInner struct {
	Pod *Pod
}

// PodAsProjectProjectIdPodGet200ResponseInner is a convenience function that returns Pod wrapped in ProjectProjectIdPodGet200ResponseInner
func PodAsProjectProjectIdPodGet200ResponseInner(v *Pod) ProjectProjectIdPodGet200ResponseInner {
	return ProjectProjectIdPodGet200ResponseInner{
		Pod: v,
	}
}


// Unmarshal JSON data into one of the pointers in the struct
func (dst *ProjectProjectIdPodGet200ResponseInner) UnmarshalJSON(data []byte) error {
	var err error
	match := 0
	// try to unmarshal data into Pod
	err = newStrictDecoder(data).Decode(&dst.Pod)
	if err == nil {
		jsonPod, _ := json.Marshal(dst.Pod)
		if string(jsonPod) == "{}" { // empty struct
			dst.Pod = nil
		} else {
			match++
		}
	} else {
		dst.Pod = nil
	}

	if match > 1 { // more than 1 match
		// reset to nil
		dst.Pod = nil

		return fmt.Errorf("data matches more than one schema in oneOf(ProjectProjectIdPodGet200ResponseInner)")
	} else if match == 1 {
		return nil // exactly one match
	} else { // no match
		return fmt.Errorf("data failed to match schemas in oneOf(ProjectProjectIdPodGet200ResponseInner)")
	}
}

// Marshal data from the first non-nil pointers in the struct to JSON
func (src ProjectProjectIdPodGet200ResponseInner) MarshalJSON() ([]byte, error) {
	if src.Pod != nil {
		return json.Marshal(&src.Pod)
	}

	return nil, nil // no data in oneOf schemas
}

// Get the actual instance
func (obj *ProjectProjectIdPodGet200ResponseInner) GetActualInstance() (interface{}) {
	if obj == nil {
		return nil
	}
	if obj.Pod != nil {
		return obj.Pod
	}

	// all schemas are nil
	return nil
}

type NullableProjectProjectIdPodGet200ResponseInner struct {
	value *ProjectProjectIdPodGet200ResponseInner
	isSet bool
}

func (v NullableProjectProjectIdPodGet200ResponseInner) Get() *ProjectProjectIdPodGet200ResponseInner {
	return v.value
}

func (v *NullableProjectProjectIdPodGet200ResponseInner) Set(val *ProjectProjectIdPodGet200ResponseInner) {
	v.value = val
	v.isSet = true
}

func (v NullableProjectProjectIdPodGet200ResponseInner) IsSet() bool {
	return v.isSet
}

func (v *NullableProjectProjectIdPodGet200ResponseInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableProjectProjectIdPodGet200ResponseInner(val *ProjectProjectIdPodGet200ResponseInner) *NullableProjectProjectIdPodGet200ResponseInner {
	return &NullableProjectProjectIdPodGet200ResponseInner{value: val, isSet: true}
}

func (v NullableProjectProjectIdPodGet200ResponseInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableProjectProjectIdPodGet200ResponseInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

