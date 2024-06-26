package engine

import (
	"context"

	"github.com/hashicorp/hcl2/hcldec"
	"github.com/podinate/podinate/kube_client"
	"github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

const (
	KubernetesKindPersistentVolumeClaim       = "PersistentVolumeClaim"
	KubernetesAPIVersionPersistentVolumeClaim = "v1"

	// Other
	KubernetesAnnotationDefaultStorageClass = "storageclass.kubernetes.io/is-default-class" // Ref: https://kubernetes.io/docs/tasks/administer-cluster/change-default-storage-class/
)

type SharedVolume struct {
	ID        string
	Size      string  `cty:"size"`
	Class     *string `cty:"class"`
	Namespace *string `cty:"namespace"`
}

var SharedVolumeHCLSpec = &hcldec.BlockMapSpec{
	TypeName:   "shared_volume",
	LabelNames: []string{"id"},
	Nested: &hcldec.ObjectSpec{
		"size": &hcldec.AttrSpec{
			Name:     "size",
			Type:     cty.String,
			Required: true,
		},
		"class": &hcldec.AttrSpec{
			Name:     "class",
			Type:     cty.String,
			Required: false,
		},
		"namespace": &hcldec.AttrSpec{
			Name:     "namespace",
			Type:     cty.String,
			Required: false,
		},
	},
}

// PlanSharedVolumeChanges takes a SharedVolume from a plan and determines what needs to be done to make it match the desired state.
// If nothing needs to be done, returns nil, nil
func (sv *SharedVolume) PlanChanges(ctx context.Context) (*Change, error) {
	pvcSpec, err := sv.ToPVC(ctx)
	if err != nil {
		return nil, err
	}

	rc, err := GetResourceChangeForResource(ctx, pvcSpec)
	if errors.IsInvalid(err) {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"resource_type": ResourceTypeSharedVolume,
			"resource_id":   sv.ID,
			"error":         err,
		}).Trace("Cannot update a SharedVolume after creation. To resize the volume, edit the PersistentVolume that backs it.")
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	logrus.WithContext(ctx).WithFields(logrus.Fields{
		"resource_type": ResourceTypeSharedVolume,
		"resource_id":   sv.ID,
	}).Debug("Planned change for SharedVolume")

	if rc == nil {
		return &Change{
			ChangeType:   ChangeTypeNoop,
			ResourceType: ResourceTypeSharedVolume,
			ResourceID:   sv.ID,
		}, nil
	}

	return &Change{
		ChangeType:   rc.ChangeType,
		ResourceType: ResourceTypeSharedVolume,
		ResourceID:   sv.ID,
		Changes:      &[]ResourceChange{*rc},
	}, nil

}

// ToPVC converts a SharedVolume to a PVC
func (sv *SharedVolume) ToPVC(ctx context.Context) (*corev1.PersistentVolumeClaim, error) {
	size, err := resource.ParseQuantity(sv.Size)
	if err != nil {
		return nil, err
	}

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: sv.ID,
			Finalizers: []string{
				"kubernetes.io/pvc-protection",
			},
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: size,
				},
			},

			VolumeMode: func(val corev1.PersistentVolumeMode) *corev1.PersistentVolumeMode { return &val }(corev1.PersistentVolumeFilesystem),
		},
	}
	pvc.Kind = KubernetesKindPersistentVolumeClaim
	pvc.APIVersion = KubernetesAPIVersionPersistentVolumeClaim

	if sv.Class == nil {
		defaultClass, err := GetDefaultStorageClass(ctx)
		if err != nil {
			return nil, err
		}
		if defaultClass == nil {
			return nil, errors.NewInvalid(schema.ParseGroupKind(KubernetesKindPersistentVolumeClaim), sv.ID, field.ErrorList{field.Invalid(field.NewPath("class"), *sv.Class, "No default StorageClass found create one using this documentation https://kubernetes.io/docs/tasks/administer-cluster/change-default-storage-class/")})
		}
		pvc.Spec.StorageClassName = defaultClass
	} else {
		ok, err := sv.CheckStorageClassExists(ctx)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.NewInvalid(schema.ParseGroupKind(KubernetesKindPersistentVolumeClaim), sv.ID, field.ErrorList{field.Invalid(field.NewPath("class"), *sv.Class, "StorageClass does not exist")})
		}
		pvc.Spec.StorageClassName = sv.Class
	}

	if sv.Namespace != nil {
		pvc.ObjectMeta.Namespace = *sv.Namespace
	}

	return pvc, nil
}

// CheckStorageClassExists checks if the storage class exists
func (sv *SharedVolume) CheckStorageClassExists(ctx context.Context) (bool, error) {
	if sv.Class == nil {
		return false, nil
	}

	clientSet, err := kube_client.Client()
	if err != nil {
		return false, err
	}

	// Get the storage class
	sc, err := clientSet.StorageV1().StorageClasses().Get(ctx, *sv.Class, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if sc == nil {
		return false, nil
	}

	return true, nil
}

// GetDefaultStorageClass gets the default storage class from Kubernetes
func GetDefaultStorageClass(ctx context.Context) (*string, error) {
	clientSet, err := kube_client.Client()
	if err != nil {
		return nil, err
	}

	scList, err := clientSet.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, sc := range scList.Items {
		if sc.Annotations[KubernetesAnnotationDefaultStorageClass] == "true" {
			return &sc.Name, nil
		}
	}

	return nil, nil
}
