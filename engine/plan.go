// Plan.go controls creating a Plan for changes needed to reach a new desired state on the Kubernetes cluster.

package engine

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/podinate/podinate/kube_client"
	"github.com/podinate/podinate/tui"
	"github.com/sirupsen/logrus"
	"github.com/sters/yaml-diff/yamldiff"
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
		Applied: true,
	}

	// Create a plan for the Namespace
	namespaceChanges, err := planNamespaceChanges(ctx, client, pkg.Namespace)
	logrus.WithContext(ctx).WithFields(logrus.Fields{
		"namespace": pkg.Namespace,
		"changes":   namespaceChanges,
		"error":     err,
	}).Trace("Planned namespace changes")
	if err != nil {
		return nil, err
	}
	plan.Changes = append(plan.Changes, *namespaceChanges)

	logrus.WithContext(ctx).WithFields(logrus.Fields{
		"shared_volumes":      pkg.SharedVolumes,
		"shared_volume_count": len(pkg.SharedVolumes),
	}).Debug("Planning shared volumes")

	for _, sv := range pkg.SharedVolumes {

		svplan, err := sv.PlanChanges(ctx)
		if err != nil {
			return nil, err
		}

		plan.Changes = append(plan.Changes, *svplan)

	}

	// Create a plan for each Pod
	for _, pod := range pkg.Pods {
		podPlan, err := planPodChanges(ctx, client, pkg, pod)
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"pod":     pod.ID,
			"changes": podPlan,
			"error":   err,
		}).Trace("Planned pod changes")
		if err != nil {
			return nil, err
		}
		plan.Changes = append(plan.Changes, *podPlan)
	}

	// We got this far, the plan must be valid
	plan.Valid = true

	for _, change := range plan.Changes {
		if change.ChangeType != ChangeTypeNoop {
			plan.Applied = false
			break
		}
	}

	return &plan, nil
}

// Display shows the plan to the user
func (plan *Plan) Display() error {
	var created, updated, deleted, noop int

	y := printers.YAMLPrinter{}

	for _, change := range plan.Changes {
		switch change.ChangeType {
		case ChangeTypeCreate:
			fmt.Printf("%s "+tui.StyleItalic.Render("%s")+" will be "+tui.StyleSuccess.Render("created")+":\n", change.ResourceType, change.ResourceID)
			created++
		case ChangeTypeUpdate:
			fmt.Printf("%s %s will be "+tui.StyleUpdated.Render("updated")+"\n", change.ResourceType, change.ResourceID)

			updated++
		case ChangeTypeDelete:
			fmt.Printf("%s %s will be deleted\n", change.ResourceType, change.ResourceID)
			deleted++
		case ChangeTypeNoop:
			fmt.Printf("%s %s is up to date\n", change.ResourceType, change.ResourceID)
			noop++
		}
		if change.Changes != nil {
			for i, c := range *change.Changes {

				if i > 0 {
					fmt.Println(tui.StyleSuccess.Render("---"))
				}

				if c.ChangeType == ChangeTypeUpdate {
					err := YamlDiffResources(c.CurrentResource, c.DesiredResource)
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"error": err,
						}).Error("Error diffing resources")
						return err
					}
				} else if c.ChangeType == ChangeTypeCreate {
					logrus.WithFields(logrus.Fields{
						"kind": c.DesiredResource.GetObjectKind().GroupVersionKind().Kind,
					}).Debug("creating object")

					err := stripManagedFields(c.DesiredResource)
					if err != nil {
						return err
					}

					err = y.PrintObj(c.DesiredResource, os.Stdout)
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"error": err,
						}).Error("Error printing object")
						return err
					}
				}
			}
		}
	}

	fmt.Printf("Summary: %d created, %d updated, %d deleted, %d unchanged\n", created, updated, deleted, noop)

	if created == 0 && updated == 0 && deleted == 0 {
		fmt.Println("\n" + tui.StyleSuccess.Render("Everything up to date. Nothing to do."))

	}

	return nil
}

func YamlDiffResources(oldResource runtime.Object, newResource runtime.Object) error {
	y := printers.YAMLPrinter{}
	b := new(bytes.Buffer)

	// Strip the ManagedFields from the old resource
	// as they're not suitable for comparison
	err := stripManagedFields(oldResource)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":        err,
			"old_resource": oldResource,
			"new_resource": newResource,
		}).Error("Error stripping managed fields from oldResource")
		return err
	}

	err = y.PrintObj(oldResource, b)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":        err,
			"old_resource": oldResource,
			"new_resource": newResource,
		}).Error("Error printing old object")
		return err
	}
	oldstr := b.String()
	b.Truncate(0)

	yamldiff.Load(oldstr)

	// Strip managed fields to make comparison easier
	err = stripManagedFields(newResource)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":        err,
			"old_resource": oldResource,
			"new_resource": newResource,
		}).Error("Error stripping managed fields from newResource")
		return err
	}

	err = y.PrintObj(newResource, b)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":        err,
			"old_resource": oldResource,
			"new_resource": newResource,
		}).Error("Error printing new object")
		return err
	}
	newstr := b.String()

	oldYaml, err := yamldiff.Load(oldstr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error loading old YAML")
		return err
	}

	newYaml, err := yamldiff.Load(newstr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error loading new YAML")
		return err
	}

	for _, diff := range yamldiff.Do(oldYaml, newYaml) {
		fmt.Println(diff.Dump())
	}

	return nil
}

// Apply applies the plan to the Kubernetes cluster
func (plan *Plan) Apply(ctx context.Context) error {

	logrus.Info("Applying plan")
	// Connect to the Kube
	client, err := kube_client.Client()
	if err != nil {
		return err
	}

	restConfig, err := kube_client.GetRestConfig()
	if err != nil {
		return err
	}

	// "Plan changes" are the changes as the user sees them
	for _, change := range plan.Changes {

		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"resourceType": change.ResourceType,
			"resourceID":   change.ResourceID,
			"changeType":   change.ChangeType,
		}).Info("Applying change")

		// "Resource changes" are the changes to the actual resources
		if change.Changes != nil {
			for _, resourceChange := range *change.Changes {
				switch resourceChange.ChangeType {
				case ChangeTypeCreate:

					_, err := createObject(client, *restConfig, resourceChange.DesiredResource, false)
					if err != nil {
						return err
					}
					logrus.WithContext(ctx).WithFields(logrus.Fields{
						"objectType": change.ResourceType,
						"objectID":   change.ResourceID,
						"object":     resourceChange.DesiredResource,
					}).Info("Created object")

				case ChangeTypeUpdate:
					_, err := createObject(client, *restConfig, resourceChange.DesiredResource, true)
					if err != nil {
						return err
					}
					logrus.WithContext(ctx).WithFields(logrus.Fields{
						"objectType": change.ResourceType,
						"objectID":   change.ResourceID,
						"object":     resourceChange.DesiredResource,
					}).Info("Updated object")

				case ChangeTypeDelete:
					fmt.Print("Delete not implemented\n")
				case ChangeTypeNoop:
					// Do nothing
					continue
				}

			}
		}

	}

	plan.Applied = true

	fmt.Println("\n", tui.StyleSuccess.Render("All changes applied successfully"))

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
	// Get the resources for the pod
	ct, resourceChanges, err := pod.GetResources(ctx, pkg)
	if err != nil {
		return nil, err
	}

	logrus.WithContext(ctx).WithFields(logrus.Fields{
		"pod": pod.ID,
	}).Debug("Appending changes")

	var change = &Change{
		ResourceType: ResourceTypePod,
		ResourceID:   pod.ID,
		ChangeType:   *ct,
	}

	change.Changes = new([]ResourceChange)
	*change.Changes = append(*change.Changes, resourceChanges...)

	return change, nil

}
