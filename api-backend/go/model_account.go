/*
 * Podinate API
 *
 * The API for the simple containerisation solution Podinate.
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type Account struct {

	// The short name (slug/url) of the account. The account ID is globally unique and cannot be changed.
	Id string `json:"id,omitempty"`

	// The human readable name of the account, used for display purposes.
	Name string `json:"name,omitempty"`

	// The global Resource ID of the account
	ResourceId string `json:"resource_id,omitempty"`
}

// AssertAccountRequired checks if the required fields are not zero-ed
func AssertAccountRequired(obj Account) error {
	return nil
}

// AssertRecurseAccountRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of Account (e.g. [][]Account), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseAccountRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aAccount, ok := obj.(Account)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertAccountRequired(aAccount)
	})
}
