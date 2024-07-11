package helpers

import (
	"errors"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// resourceToUnstructured converts a runtime.Object to an unstructured.Unstructured
func ObjectToUnstructured(obj runtime.Object) (*unstructured.Unstructured, error) {
	if obj == nil {
		return nil, errors.New("object is nil")
	}
	innerObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	return &unstructured.Unstructured{Object: innerObj}, nil
}

// unstructuredToObject converts an unstructured.Unstructured to a runtime.Object
func UnstructuredToObject(u *unstructured.Unstructured) (runtime.Object, error) {
	if u == nil {
		return nil, errors.New("unstructured object is nil")
	}
	return u.DeepCopyObject(), nil
}
