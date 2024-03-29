/*
 * Podinate API
 *
 * The API for the simple containerisation solution Podinate.
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type UserLoginPost200Response struct {

	// Whether the user is logged in or not
	LoggedIn bool `json:"logged_in,omitempty"`

	// The user's API key, if they are logged in
	ApiKey string `json:"api_key,omitempty"`
}

// AssertUserLoginPost200ResponseRequired checks if the required fields are not zero-ed
func AssertUserLoginPost200ResponseRequired(obj UserLoginPost200Response) error {
	return nil
}

// AssertRecurseUserLoginPost200ResponseRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of UserLoginPost200Response (e.g. [][]UserLoginPost200Response), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseUserLoginPost200ResponseRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aUserLoginPost200Response, ok := obj.(UserLoginPost200Response)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertUserLoginPost200ResponseRequired(aUserLoginPost200Response)
	})
}
