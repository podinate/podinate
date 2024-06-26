// Apply contains helpers for applying resource changes

package engine

import (
	"errors"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

// Ref: https://stackoverflow.com/questions/58783939 kubectl apply in go

func createObject(kubeClientset kubernetes.Interface, restConfig rest.Config, obj runtime.Object, update bool) (runtime.Object, error) {
	// Create a REST mapper that tracks information about the available resources in the cluster.
	groupResources, err := restmapper.GetAPIGroupResources(kubeClientset.Discovery())
	if err != nil {
		return nil, err
	}
	rm := restmapper.NewDiscoveryRESTMapper(groupResources)

	// Get some metadata needed to make the REST request.
	gvk := obj.GetObjectKind().GroupVersionKind()
	gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}
	mapping, err := rm.RESTMapping(gk, gvk.Version)
	if err != nil {
		return nil, err
	}

	namespace, err := meta.NewAccessor().Namespace(obj)
	if err != nil {
		return nil, err
	}

	// Create a client specifically for creating the object.
	restClient, err := newRestClient(restConfig, mapping.GroupVersionKind.GroupVersion())
	if err != nil {
		return nil, err
	}

	// Use the REST helper to create the object in the "default" namespace.
	restHelper := resource.NewHelper(restClient, mapping)

	logrus.WithFields(logrus.Fields{
		"namespace": namespace,
		"update":    update,
		"obj":       obj,
	}).Info("Creating object")

	if update {
		// Get the name out of the runtime.Object
		u, err := resourceToUnstructured(obj)
		if err != nil {
			return nil, err
		}

		return restHelper.Replace(namespace, u.GetName(), update, obj)
	}

	return restHelper.Create(namespace, update, obj)
}

func newRestClient(restConfig rest.Config, gv schema.GroupVersion) (rest.Interface, error) {
	restConfig.ContentConfig = resource.UnstructuredPlusDefaultContentConfig()
	restConfig.GroupVersion = &gv
	if len(gv.Group) == 0 {
		restConfig.APIPath = "/api"
	} else {
		restConfig.APIPath = "/apis"
	}

	return rest.RESTClientFor(&restConfig)
}

// resourceToUnstructured converts a runtime.Object to an unstructured.Unstructured
func resourceToUnstructured(obj runtime.Object) (*unstructured.Unstructured, error) {
	if obj == nil {
		return nil, errors.New("object is nil")
	}
	innerObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	return &unstructured.Unstructured{Object: innerObj}, nil
}

// stripManagedFields removes the "managedFields" field from the object
func stripManagedFields(resource runtime.Object) error {
	// Strip ManagedFields from the old resource
	o, err := resourceToUnstructured(resource)
	if err != nil {
		return err
	}
	o.SetManagedFields(nil)
	return runtime.DefaultUnstructuredConverter.FromUnstructured(o.Object, resource)
}
