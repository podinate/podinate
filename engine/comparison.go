/*
Package comparison_engine provides the functionality to take any Kubernetes resource and
determine what needs to be done to make it match the desired state.
*/
package engine

import (
	"context"
	"fmt"
	"reflect"

	"github.com/podinate/podinate/kube_client"
	"github.com/podinate/podinate/tui"
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
)

// GetResourceChangeForResource takes a runtime.Object and determines what needs to be done to make it match the desired state.
// If nothing needs to be done, returns nil, nil
func GetResourceChangeForResource(ctx context.Context, object runtime.Object) (*ResourceChange, error) {
	helper, err := getRestHelperForObject(object)
	if err != nil {
		return nil, err
	}

	unstructuredObject, err := resourceToUnstructured(object)
	if err != nil {
		return nil, err
	}

	// Get the current state of the object
	currentObject, err := helper.Get(unstructuredObject.GetNamespace(), unstructuredObject.GetName())
	if errors.IsNotFound(err) { // Resource does not exist and needs to be created
		logrus.WithFields(logrus.Fields{
			"object":        object,
			"error":         err,
			"currentObject": currentObject,
		}).Trace("Object not found, needs to be created")

		rc := ResourceChange{
			ChangeType:      ChangeTypeCreate,
			CurrentResource: nil,
			DesiredResource: object,
		}

		return &rc, nil
	} else if err != nil { // Handle any other error
		logrus.WithFields(logrus.Fields{
			"object":        object,
			"error":         err,
			"currentObject": currentObject,
		}).Error("Error getting the resource ")
		return nil, err
	}

	// At this point, the resource exists and we need to determine if it needs to be updated
	dryRunResult, err := helper.DryRun(true).Replace(unstructuredObject.GetNamespace(), unstructuredObject.GetName(), false, object)
	if errors.IsInvalid(err) { // The object is invalid, user will need to change something
		logrus.WithFields(logrus.Fields{
			"object":        object,
			"error":         err,
			"currentObject": currentObject,
		}).Error(HelpInvalidObject)

		fmt.Println(tui.StyleError.Render("The following change was rejected by Kubernetes:"))

		err := YamlDiffResources(currentObject, object)
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"error": err,
			}).Error("A further error occurred when trying to display the change that was rejected by the Kubernetes API.")
		}
		return nil, err
	} else if err != nil { // Any other error, may need to add handlers for more custom errors in future
		logrus.WithFields(logrus.Fields{
			"object":        object,
			"error":         err,
			"currentObject": currentObject,
		}).Error("Error ")
		return nil, err
	}

	// If the object has not changed, don't do anything
	if reflect.DeepEqual(currentObject, dryRunResult) {
		logrus.WithFields(logrus.Fields{
			"object":        object,
			"error":         err,
			"currentObject": currentObject,
		}).Trace("Object is the same, no change needed")
		return nil, nil
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

	return restHelper, nil
}
