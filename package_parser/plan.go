// Plan.go controls creating a Plan for changes needed to reach a new desired state on the Kubernetes cluster.

package package_parser

import (
	"context"
	"fmt"
	"os"

	"github.com/podinate/podinate/kube_client"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes"
)

type Plan struct {
	// Valid means the plan makes sense and can be applied
	Valid bool
	// Applied means the plan has finished applying
	Applied bool
	// An array of changes to make in order
	Changes []Change
}

type ChangeType string

const (
	ChangeTypeCreate ChangeType = "create"
	ChangeTypeUpdate ChangeType = "update"
	ChangeTypeDelete ChangeType = "delete"
	ChangeTypeNoop   ChangeType = "noop"
)

type ResourceType string

const (
	ResourceTypePod          ResourceType = "pod"
	ResourceTypeNamespace    ResourceType = "namespace"
	ResourceTypeService      ResourceType = "service"
	ResourceTypeVolume       ResourceType = "volume"
	ResourceTypeSharedVolume ResourceType = "shared_volume"
)

// Change represents a change to an overall resource,
// such as a Pod, Namespace, or SharedVolume
type Change struct {
	ResourceType ResourceType
	ResourceID   string
	ChangeType   ChangeType
	// Changes are applied in the order they appear in the slice
	Changes *[]ResourceChange
}

// ResourceChange represents a change to a Kubernetes resource
// For example, a change to a Pod's image or a Service's port
type ResourceChange struct {
	ChangeType      ChangeType
	CurrentResource runtime.Object
	DesiredResource runtime.Object
}

// CreatePlan creates a plan for the desired state of a package
func (pkg *Package) Plan(ctx context.Context) (*Plan, error) {

	// First of all, try to connect to the Kube
	client, err := kube_client.Client()
	if err != nil {
		return nil, err
	}

	plan := Plan{
		Valid:   false,
		Applied: false,
	}

	// Create a plan for the Namespace
	namespaceChanges, err := planNamespaceChanges(ctx, client, pkg.Namespace)
	if err != nil {
		return nil, err
	}
	plan.Changes = append(plan.Changes, *namespaceChanges)

	// Create a plan for each Pod
	for _, pod := range pkg.Pods {
		podPlan, err := planPodChanges(ctx, client, pkg, pod)
		if err != nil {
			return nil, err
		}
		plan.Changes = append(plan.Changes, *podPlan)
	}

	// We got this far, the plan must be valid
	plan.Valid = true

	return &plan, nil
}

// Display shows the plan to the user
func (plan *Plan) Display() {
	//logrus.Infof("Plan: %+v", plan)

	var created, updated, deleted, noop int

	y := printers.YAMLPrinter{}

	for _, change := range plan.Changes {
		switch change.ChangeType {
		case ChangeTypeCreate:
			fmt.Printf("%s %s will be created:\n", change.ResourceType, change.ResourceID)
			for _, c := range *change.Changes {
				//fmt.Printf("%s\n\n", c.DesiredResource)
				err := y.PrintObj(c.DesiredResource, os.Stdout)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"error": err,
					}).Error("Error printing object")
				}
				fmt.Println()
			}
			created++
		case ChangeTypeUpdate:
			fmt.Printf("%s %s will be updated\n", change.ResourceType, change.ResourceID)
			updated++
		case ChangeTypeDelete:
			fmt.Printf("%s %s will be deleted\n", change.ResourceType, change.ResourceID)
			deleted++
		case ChangeTypeNoop:
			fmt.Printf("%s %s is up to date\n", change.ResourceType, change.ResourceID)
			noop++
		}
	}

	fmt.Printf("Summary: %d created, %d updated, %d deleted, %d unchanged\n", created, updated, deleted, noop)
}

// Apply applies the plan to the Kubernetes cluster
func (plan *Plan) Apply(ctx context.Context) error {

	logrus.Info("Applying plan")
	// Connect to the Kube
	client, err := kube_client.Client()
	if err != nil {
		return err
	}

	restConfig, err := kube_client.RestConfig()
	if err != nil {
		return err
	}

	for _, change := range plan.Changes {

		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"resourceType": change.ResourceType,
			"resourceID":   change.ResourceID,
			"changeType":   change.ChangeType,
		}).Info("Applying change")

		switch change.ChangeType {
		case ChangeTypeCreate:

			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"resourceType": change.ResourceType,
				"resourceID":   change.ResourceID,
			}).Info("Creating resource")

			for _, c := range *change.Changes {
				_, err := createObject(client, *restConfig, c.DesiredResource, false)
				if err != nil {
					return err
				}

				logrus.WithContext(ctx).WithFields(logrus.Fields{
					"resorceType": change.ResourceType,
					"ResourceID":  change.ResourceID,
					"resource":    c.DesiredResource,
				}).Info("Created resource")
			}
		case ChangeTypeUpdate:

			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"resourceType": change.ResourceType,
				"resourceID":   change.ResourceID,
			}).Info("Updating resource")

			for _, c := range *change.Changes {
				_, err := createObject(client, *restConfig, c.DesiredResource, true)
				if err != nil {
					return err
				}

				logrus.WithContext(ctx).WithFields(logrus.Fields{
					"resorceType": change.ResourceType,
					"ResourceID":  change.ResourceID,
					"resource":    c.DesiredResource,
				}).Info("Updated resource")
			}

			fmt.Println("Update not implemented")
		case ChangeTypeDelete:
			fmt.Print("Delete not implemented\n")
		case ChangeTypeNoop:
			// Do nothing
		}
	}

	return nil
}

func planNamespaceChanges(ctx context.Context, client *kubernetes.Clientset, namespace string) (*Change, error) {
	// Conn3ct to the cluster and check if the namespace exists
	// If it doesn't, return a ChangeTypeCreate
	// If it does, return a ChangeTypeNoop
	_, err := client.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if errors.IsNotFound(err) {

		// Create the namespace runtime.Object
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		ns.Kind = "Namespace"
		ns.APIVersion = "v1"

		return &Change{
			ResourceType: ResourceTypeNamespace,
			ResourceID:   namespace,
			ChangeType:   ChangeTypeCreate,
			Changes: &[]ResourceChange{
				{
					ChangeType:      ChangeTypeCreate,
					CurrentResource: nil,
					DesiredResource: ns,
				},
			},
		}, nil
	}
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"namespace": namespace,
		}).Error("Error getting namespace")
		return nil, err
	}

	// Cannot meaningfully update a namespace
	return &Change{
		ResourceType: ResourceTypeNamespace,
		ResourceID:   namespace,
		ChangeType:   ChangeTypeNoop,
	}, nil
}

func planPodChanges(ctx context.Context, client *kubernetes.Clientset, pkg *Package, pod Pod) (*Change, error) {
	// Connect to the cluster and check if the pod's statefulset exists
	// If it doesn't, return a ChangeTypeCreate
	// If it does, return a ChangeTypeUpdate
	ss, err := client.AppsV1().StatefulSets(pkg.Namespace).Get(ctx, pod.ID, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		// Resource doesn't exist, so create it

		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"error":     err,
			"pod":       pod.ID,
			"namespace": pkg.Namespace,
		}).Debug("Pod (StatefulSet) doesn't exist, creating it")

		var change = &Change{
			ResourceType: ResourceTypePod,
			ResourceID:   pod.ID,
			ChangeType:   ChangeTypeCreate,
		}
		// Get the resources for the pod
		resources, err := pod.GetResources(ctx, pkg)
		if err != nil {
			return nil, err
		}
		change.Changes = new([]ResourceChange)
		for _, resource := range resources {
			*change.Changes = append(*change.Changes, ResourceChange{
				ChangeType:      ChangeTypeCreate,
				CurrentResource: nil,
				DesiredResource: resource,
			})
		}

		return change, nil
	}
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"pod":   pod.Name,
		}).Error("Error getting pod")
		return nil, err
	}

	// Resource exists, so update it
	var change = &Change{
		ResourceType: ResourceTypePod,
		ResourceID:   pod.ID,
		ChangeType:   ChangeTypeUpdate,
	}
	// Get the resources for the pod
	resources, err := pod.GetResources(ctx, pkg)
	if err != nil {
		return nil, err
	}

	change.Changes = new([]ResourceChange)
	for _, resource := range resources {
		*change.Changes = append(*change.Changes, ResourceChange{
			ChangeType:      ChangeTypeUpdate,
			CurrentResource: ss,
			DesiredResource: resource,
		})
	}

	return change, nil

}
