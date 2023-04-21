/*
 * Podinate API
 *
 * The API for the simple containerisation solution Podinate
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type Model500Error struct {

	// The error code for this error
	Code float32 `json:"code,omitempty"`

	// Friendly human readable description of the error
	Message string `json:"message,omitempty"`
}

// AssertModel500ErrorRequired checks if the required fields are not zero-ed
func AssertModel500ErrorRequired(obj Model500Error) error {
	return nil
}

// AssertRecurseModel500ErrorRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of Model500Error (e.g. [][]Model500Error), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseModel500ErrorRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aModel500Error, ok := obj.(Model500Error)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertModel500ErrorRequired(aModel500Error)
	})
}
