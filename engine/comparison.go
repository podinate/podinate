/*
Package comparison_engine provides the functionality to take any Kubernetes resource and
determine what needs to be done to make it match the desired state.
*/
package engine

import (
	"context"
	"fmt"
	"reflect"
	"slices"

	"github.com/podinate/podinate/engine/helpers"
	"github.com/podinate/podinate/kube_client"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/restmapper"
)

const (
	HelpInvalidObject      = "The object is invalid. Check the error for more information and update values accordingly."
	HelpUnknownUpdateError = "Error validating the resource. Check the error and update values accordingly."

	FieldValidationStrict = "Strict"
)

// GetResourceChangeForResource takes a runtime.Object and determines what needs to be done to make it match the desired state.
// If nothing needs to be done, returns nil, nil
func GetResourceChangeForResource(ctx context.Context, object runtime.Object) (*ResourceChange, error) {

	// Alright this method does most of the work

	// Get a REST helper for the object. This allows us to work with the object whatever it is
	helper, err := getRestHelperForObject(object)
	if err != nil {
		return nil, err
	}

	// Make the Object into an unstructured object
	// This allows us to tell some information about the object
	unstructuredObject, err := resourceToUnstructured(object)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"error":  err,
			"object": object,
		}).Error("Error converting object to unstructured")
		return nil, err
	}

	// Get the current state of the object
	currentObject, err := helper.Get(unstructuredObject.GetNamespace(), unstructuredObject.GetName())
	logrus.WithContext(ctx).WithFields(logrus.Fields{
		"object":        object,
		"currentObject": currentObject,
		"error":         err,
		"is_not_found":  errors.IsNotFound(err),
	}).Trace("Getting the current state of the object")
	if errors.IsNotFound(err) { // Resource does not exist and needs to be created
		logrus.WithFields(logrus.Fields{
			"object":        object,
			"error":         err,
			"currentObject": currentObject,
		}).Trace("Object not found, needs to be created")

		// Do a dry run of creating the object. This should tell us if it's valid or not, but not actually create it
		dryRunCreated, dryRunErr := helper.DryRun(true).Create(unstructuredObject.GetNamespace(), false, object)
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"object": object,
			"error":  dryRunErr,
			"result": dryRunCreated,
		}).Trace("Dry run to create object")
		if errors.IsInvalid(dryRunErr) || errors.IsBadRequest(dryRunErr) { // The object is invalid, user will need to change something
			logrus.WithFields(logrus.Fields{
				"desired_object": object,
				"error":          err,
			}).Debug(HelpInvalidObject)

			err := helpers.PrintKubernetesValidationError(ctx, nil, object, dryRunErr)
			if err != nil {
				logrus.WithContext(ctx).WithFields(logrus.Fields{
					"error": err,
				}).Error("A further error occurred when trying to display the change that was rejected by the Kubernetes API.")
				return nil, err
			}

			return nil, dryRunErr

		} else if errors.IsNotFound(dryRunErr) || dryRunErr == nil {
			rc := ResourceChange{
				ChangeType:      ChangeTypeCreate,
				CurrentResource: nil,
				DesiredResource: object,
			}

			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"resourceChange": rc,
				"kind":           object.GetObjectKind().GroupVersionKind().Kind,
			}).Trace("Resource change of type create was created")

			return &rc, nil
		} else if dryRunErr != nil { // Any other error, may need to add handlers for more errors in future
			logrus.WithFields(logrus.Fields{
				"object":        object,
				"error":         dryRunErr,
				"currentObject": currentObject,
			}).Error("Error planning to create object")
			return nil, err
		}

	} else if err != nil { // Handle any other error
		logrus.WithFields(logrus.Fields{
			"object":        object,
			"error":         err,
			"currentObject": currentObject,
		}).Error("Error getting the resource ")
		return nil, err
	}

	// Transfer labels, annotations etc from the current object to the desired object
	object, err = transferMetadata(currentObject, object)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"error":          err,
			"current_object": currentObject,
			"desired_object": object,
		}).Error("Error transferring metadata, labels etc from current object to desired object")
		return nil, err
	}

	object, err = helpers.TransferFields(currentObject, object)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"error":          err,
			"current_object": currentObject,
			"desired_object": object,
		}).Error("Error transferring fields from current object to desired object")
		return nil, err
	}

	// At this point, the resource exists and we need to determine if it needs to be updated
	fmt.Println("Doing dry run")
	dryRunResult, dryRunErr := helper.DryRun(true).Replace(unstructuredObject.GetNamespace(), unstructuredObject.GetName(), true, object)
	fmt.Println("Dry run done")
	if errors.IsInvalid(dryRunErr) || errors.IsBadRequest(dryRunErr) { // The object is invalid, user will need to change something
		logrus.WithFields(logrus.Fields{
			"desired_object": object,
			"error":          dryRunErr,
			"current_object": currentObject,
		}).Debug(HelpInvalidObject)

		err := helpers.PrintKubernetesValidationError(ctx, currentObject, object, dryRunErr)
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"error": err,
			}).Error("A further error occurred when trying to display the change that was rejected by the Kubernetes API.")
		}

		return nil, dryRunErr
	} else if dryRunErr != nil { // Any other error, may need to add handlers for more custom errors in future
		logrus.WithFields(logrus.Fields{
			"object":        object,
			"error":         dryRunErr,
			"currentObject": currentObject,
			"object_error":  errors.IsBadRequest(dryRunErr),
		}).Error("Error ")
		return nil, dryRunErr
	}

	// If the object has not changed, don't do anything
	if reflect.DeepEqual(currentObject, dryRunResult) {
		logrus.WithFields(logrus.Fields{
			"object":        object,
			"error":         err,
			"currentObject": currentObject,
		}).Trace("Object is the same, no change needed")

		rc := ResourceChange{
			ChangeType:      ChangeTypeNoop,
			CurrentResource: currentObject,
			DesiredResource: object,
		}

		return &rc, nil
	}

	// If we reach here, the object is different, we need to update it
	logrus.WithFields(logrus.Fields{
		"object":         object,
		"current_object": currentObject,
		"dry_run_result": dryRunResult,
	}).Trace("Object is different, needs to be updated")
	var resourceChange = ResourceChange{
		ChangeType:      ChangeTypeUpdate,
		CurrentResource: currentObject,
		DesiredResource: dryRunResult,
	}

	return &resourceChange, nil
}

// transferMetadata transfers labels, annotations etc from the current object to the desired object
func transferMetadata(currentObject runtime.Object, desiredObject runtime.Object) (runtime.Object, error) {

	uCurrent, err := resourceToUnstructured(currentObject)
	if err != nil {
		logrus.WithContext(context.Background()).WithFields(logrus.Fields{
			"error":          err,
			"current_object": currentObject,
			"desired_object": desiredObject,
		}).Error("Error converting current object to unstructured")
		return nil, err
	}

	uDesired, err := resourceToUnstructured(desiredObject)
	if err != nil {
		logrus.WithContext(context.Background()).WithFields(logrus.Fields{
			"error":          err,
			"current_object": currentObject,
			"desired_object": desiredObject,
		}).Error("Error converting desired object to unstructured")
		return nil, err
	}

	// Transfer labels, annotations etc from the current object to the desired object
	currentLabels := uCurrent.GetLabels()
	if currentLabels == nil {
		currentLabels = make(map[string]string)
	}
	currentAnnotations := uCurrent.GetAnnotations()
	if currentAnnotations == nil {
		currentAnnotations = make(map[string]string)
	}

	newLabels := uDesired.GetLabels()
	newAnnotations := uDesired.GetAnnotations()

	for k, v := range newLabels {
		currentLabels[k] = v
	}
	for k, v := range newAnnotations {
		currentAnnotations[k] = v
	}

	uDesired.SetLabels(currentLabels)
	uDesired.SetAnnotations(currentAnnotations)

	finalizers := uCurrent.GetFinalizers()
	newFinalizers := uDesired.GetFinalizers()

	logrus.WithFields(logrus.Fields{
		"currentFinalizers": finalizers,
		"newFinalizers":     newFinalizers,
	}).Trace("Finalizers")

	for _, v := range newFinalizers {
		if !slices.Contains(finalizers, v) {
			finalizers = append(finalizers, v)
		}
	}
	uDesired.SetFinalizers(finalizers)

	return uDesired.DeepCopyObject(), nil
}

// getRestHelperForObject returns a Kubernetes REST helper for a given object
// This RestHelper will understand the object and be able to make REST requests
func getRestHelperForObject(object runtime.Object) (*resource.Helper, error) {
	client, err := kube_client.Client()
	if err != nil {
		return nil, err
	}

	restConfig, err := kube_client.GetRestConfig()
	if err != nil {
		return nil, err
	}

	// Create a REST mapper that tracks information about the available resources in the cluster.
	groupResources, err := restmapper.GetAPIGroupResources(client.Discovery())
	if err != nil {
		return nil, err
	}
	rm := restmapper.NewDiscoveryRESTMapper(groupResources)

	// Get some metadata needed to make the REST request.
	gvk := object.GetObjectKind().GroupVersionKind()
	gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}
	mapping, err := rm.RESTMapping(gk, gvk.Version)
	if err != nil {
		return nil, err
	}

	// Create a restClient which can understand any resource type
	restClient, err := newRestClient(*restConfig, mapping.GroupVersionKind.GroupVersion())
	if err != nil {
		return nil, err
	}

	restHelper := resource.NewHelper(restClient, mapping)

	// Make the Helper do strict validation
	restHelper = restHelper.WithFieldValidation(FieldValidationStrict)

	return restHelper, nil
}
