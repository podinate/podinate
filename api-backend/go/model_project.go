/*
 * Podinate API
 *
 * The API for the simple containerisation solution Podinate. Login should be performed over oauth from [auth.podinate.com](https://auth.podinate.com)
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type Project struct {

	// The short name (slug/url) of the project
	Id string `json:"id,omitempty"`

	// The name of the app
	Name string `json:"name,omitempty"`

	// The container image to run for this app
	Image string `json:"image,omitempty"`

	// The image tag to run for this app
	Tag string `json:"tag,omitempty"`
}

// AssertProjectRequired checks if the required fields are not zero-ed
func AssertProjectRequired(obj Project) error {
	return nil
}

// AssertRecurseProjectRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of Project (e.g. [][]Project), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseProjectRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aProject, ok := obj.(Project)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertProjectRequired(aProject)
	})
}
