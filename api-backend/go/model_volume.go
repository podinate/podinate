/*
 * Podinate API
 *
 * The API for the simple containerisation solution Podinate.
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// Volume - A storage volume that is attached to each instance of a pod.
type Volume struct {

	// The name of the volume
	Name string `json:"name"`

	// The size of the volume in GB
	Size int32 `json:"size"`

	// The path to mount the volume at
	MountPath string `json:"mount_path"`

	// The class of the volume, for example \"standard\" or \"premium\"
	Class string `json:"class,omitempty"`
}

// AssertVolumeRequired checks if the required fields are not zero-ed
func AssertVolumeRequired(obj Volume) error {
	elements := map[string]interface{}{
		"name":       obj.Name,
		"size":       obj.Size,
		"mount_path": obj.MountPath,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertRecurseVolumeRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of Volume (e.g. [][]Volume), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseVolumeRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aVolume, ok := obj.(Volume)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertVolumeRequired(aVolume)
	})
}
