package helpers

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// TransferFields takes two runtime.Objects and transfers some fields that have to be brought forward from the current object to the new object
// For example, if you are updating a PersistentVolumeClaim, you have to transfer the UID from the current object to the new object
func TransferFields(currentObject runtime.Object, newObject runtime.Object) (runtime.Object, error) {

	unstructuredCurrent, err := ObjectToUnstructured(currentObject)
	if err != nil {
		return nil, err
	}
	unstructuredNew, err := ObjectToUnstructured(newObject)
	if err != nil {
		return nil, err
	}
	switch currentObject.GetObjectKind().GroupVersionKind().Kind {
	case "PersistentVolumeClaim":
		var currentPVC *corev1.PersistentVolumeClaim
		var newPVC *corev1.PersistentVolumeClaim
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredCurrent.Object, &currentPVC)
		if err != nil {
			return nil, err
		}

		err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredNew.Object, &newPVC)
		if err != nil {
			return nil, err
		}
		newPVC.Spec.VolumeName = currentPVC.Spec.VolumeName
		return newPVC, nil

	}

	// Handle the default case
	return newObject, nil
}
