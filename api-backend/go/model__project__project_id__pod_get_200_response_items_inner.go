/*
 * Podinate API
 *
 * The API for the simple containerisation solution Podinate.
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type ProjectProjectIdPodGet200ResponseItemsInner struct {

	// The short name (slug/url) of the pod
	Id string `json:"id,omitempty"`

	// The name of the pod
	Name string `json:"name"`

	// The container image to run for this pod
	Image string `json:"image"`

	// The image tag to run for this pod
	Tag string `json:"tag"`

	// The storage volumes attached to this pod
	Volumes []Volume `json:"volumes,omitempty"`

	// The environment variables to pass to the pod
	Environment []EnvironmentVariable `json:"environment,omitempty"`

	// The services to expose for this pod
	Services []Service `json:"services,omitempty"`

	// The current status of the pod
	Status string `json:"status,omitempty"`

	// The date and time the pod was created
	CreatedAt string `json:"created_at,omitempty"`

	// The global Resource ID of the pod
	ResourceId string `json:"resource_id,omitempty"`
}

// AssertProjectProjectIdPodGet200ResponseItemsInnerRequired checks if the required fields are not zero-ed
func AssertProjectProjectIdPodGet200ResponseItemsInnerRequired(obj ProjectProjectIdPodGet200ResponseItemsInner) error {
	elements := map[string]interface{}{
		"name":  obj.Name,
		"image": obj.Image,
		"tag":   obj.Tag,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.Volumes {
		if err := AssertVolumeRequired(el); err != nil {
			return err
		}
	}
	for _, el := range obj.Environment {
		if err := AssertEnvironmentVariableRequired(el); err != nil {
			return err
		}
	}
	for _, el := range obj.Services {
		if err := AssertServiceRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertRecurseProjectProjectIdPodGet200ResponseItemsInnerRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ProjectProjectIdPodGet200ResponseItemsInner (e.g. [][]ProjectProjectIdPodGet200ResponseItemsInner), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseProjectProjectIdPodGet200ResponseItemsInnerRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aProjectProjectIdPodGet200ResponseItemsInner, ok := obj.(ProjectProjectIdPodGet200ResponseItemsInner)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertProjectProjectIdPodGet200ResponseItemsInnerRequired(aProjectProjectIdPodGet200ResponseItemsInner)
	})
}
