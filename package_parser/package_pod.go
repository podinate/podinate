package package_parser

import (
	"context"
	"reflect"

	hcldec "github.com/hashicorp/hcl2/hcldec"
	"github.com/podinate/podinate/kube_client"
	"github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Pod struct {
	ID string
	// Name        string   `cty:"name"`
	Image       string   `cty:"image"`
	Tag         *string  `cty:"tag"`
	Command     []string `cty:"command"`
	Arguments   []string `cty:"arguments"`
	Environment map[string]struct {
		Value  string `cty:"value"`
		Secret *bool  `cty:"secret"`
	} `cty:"environment"`

	// Removed for now, because these are separate Kube resources
	// Service map[string]struct {
	// 	Port       int     `cty:"port"`
	// 	TargetPort *int    `cty:"target_port"`
	// 	Protocol   *string `cty:"protocol"`
	// 	DomainName *string `cty:"domain_name"`
	// } `cty:"service"`
	// Volume map[string]struct {
	// 	Size  int     `cty:"size"`
	// 	Path  string  `cty:"path"`
	// 	Class *string `cty:"class"`
	// } `cty:"volume"`
	// SharedVolume *[]struct {
	// 	VolumeID string `cty:"volume_id"`
	// 	Path     string `cty:"path"`
	// } `cty:"shared_volume"`
}

// GetHCLSpect returns the HCL spec of the pod type
var podHCLSpec = &hcldec.BlockMapSpec{
	TypeName:   "pod",
	LabelNames: []string{"id"},
	Nested: &hcldec.ObjectSpec{
		// "name": &hcldec.AttrSpec{
		// 	Name:     "name",
		// 	Type:     cty.String,
		// 	Required: true,
		// },
		"image": &hcldec.AttrSpec{
			Name:     "image",
			Type:     cty.String,
			Required: true,
		},
		"tag": &hcldec.AttrSpec{
			Name:     "tag",
			Type:     cty.String,
			Required: false,
		},
		"command": &hcldec.AttrSpec{
			Name:     "command",
			Type:     cty.List(cty.String),
			Required: false,
		},
		"arguments": &hcldec.AttrSpec{
			Name:     "arguments",
			Type:     cty.List(cty.String),
			Required: false,
		},
		"environment": &hcldec.BlockMapSpec{
			TypeName:   "environment",
			LabelNames: []string{"key"},
			Nested: &hcldec.ObjectSpec{
				"value": &hcldec.AttrSpec{
					Name:     "value",
					Type:     cty.String,
					Required: true,
				},
				"secret": &hcldec.AttrSpec{
					Name:     "secret",
					Type:     cty.Bool,
					Required: false,
				},
			},
		},
		// "service": &hcldec.BlockMapSpec{
		// 	TypeName:   "service",
		// 	LabelNames: []string{"name"},
		// 	Nested: &hcldec.ObjectSpec{
		// 		"port": &hcldec.AttrSpec{
		// 			Name:     "port",
		// 			Type:     cty.Number,
		// 			Required: true,
		// 		},
		// 		"target_port": &hcldec.AttrSpec{
		// 			Name:     "target_port",
		// 			Type:     cty.Number,
		// 			Required: false,
		// 		},
		// 		"protocol": &hcldec.AttrSpec{
		// 			Name:     "protocol",
		// 			Type:     cty.String,
		// 			Required: false,
		// 		},
		// 		"domain_name": &hcldec.AttrSpec{
		// 			Name:     "domain_name",
		// 			Type:     cty.String,
		// 			Required: false,
		// 		},
		// 	},
		// },
		// "volume": &hcldec.BlockMapSpec{
		// 	TypeName:   "volume",
		// 	LabelNames: []string{"name"},
		// 	Nested: &hcldec.ObjectSpec{
		// 		"size": &hcldec.AttrSpec{
		// 			Name:     "size",
		// 			Type:     cty.Number,
		// 			Required: true,
		// 		},
		// 		"path": &hcldec.AttrSpec{
		// 			Name:     "path",
		// 			Type:     cty.String,
		// 			Required: true,
		// 		},
		// 		"class": &hcldec.AttrSpec{
		// 			Name:     "class",
		// 			Type:     cty.String,
		// 			Required: false,
		// 		},
		// 	},
		// },
		// "shared_volume": &hcldec.BlockListSpec{
		// 	TypeName: "shared_volume",
		// 	MinItems: 0,
		// 	Nested: &hcldec.ObjectSpec{
		// 		"volume_id": &hcldec.AttrSpec{
		// 			Name:     "volume_id",
		// 			Type:     cty.String,
		// 			Required: true,
		// 		},
		// 		"path": &hcldec.AttrSpec{
		// 			Name:     "path",
		// 			Type:     cty.String,
		// 			Required: true,
		// 		},
		// 	},
		// },
	},
}

// GetResources returns the resources needed for the Pod
func (p *Pod) GetResources(ctx context.Context, pkg *Package) (*ChangeType, []runtime.Object, error) {
	var out []runtime.Object

	var imageID = p.Image
	if p.Tag != nil {
		imageID = imageID + ":" + *p.Tag
	}

	ss := &appsv1.StatefulSet{

		ObjectMeta: metav1.ObjectMeta{
			Name:      p.ID,
			Namespace: pkg.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: func(val int32) *int32 { return &val }(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"podinate.com/pod": p.ID,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"podinate.com/pod": p.ID,
					},
					Name: p.ID,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    p.ID,
							Image:   imageID,
							Command: p.Command,
							Args:    p.Arguments,
						},
					},
				},
			},
		},
	}

	ss.Kind = "StatefulSet"
	ss.APIVersion = "apps/v1"

	logrus.WithContext(ctx).WithFields(logrus.Fields{
		"statefulset": ss,
		"kind":        ss.GetObjectKind(),
	}).Debug("StatefulSet")

	// Check if the StatefulSet already exists
	kube, err := kube_client.Client()
	if err != nil {
		return nil, nil, err
	}

	ss, err = kube.AppsV1().StatefulSets(pkg.Namespace).Update(ctx, ss, metav1.UpdateOptions{DryRun: []string{metav1.DryRunAll}})
	if err != nil {
		return nil, nil, err
	}

	existing, err := kube.AppsV1().StatefulSets(pkg.Namespace).Get(ctx, p.ID, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		// Resource doesn't exist, so create it
		out = append(out, ss)
	} else if err != nil {
		return nil, nil, err
	} else {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"statefulset": ss,
		}).Debug("StatefulSet exists, deep comparing")
		// Deep compare the object to see if it needs updating
		if reflect.DeepEqual(ss.Spec, existing.Spec) {
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"statefulset": ss.Spec,
				"existing":    existing.Spec,
			}).Debug("StatefulSet is up to date")
			return func(s ChangeType) *ChangeType { return &s }(ChangeTypeNoop), nil, nil
		} else {
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"statefulset": ss.Spec,
				"existing":    existing.Spec,
			}).Debug("StatefulSet needs updating")

			// Resource exists, but needs updating
			out = append(out, ss)
			return func(s ChangeType) *ChangeType { return &s }(ChangeTypeUpdate), out, nil
		}
	}

	//out = append(out, ss)

	// TODO: Add services and volumes

	return nil, out, nil
}

// Tosdk returns the API client representation of the pod

// TODO: Update with a function to generate the kube config for the pod

// func (p *Pod) ToSDK() (*sdk.Pod, error) {
// 	theProject, err := sdk.GetProjectByID(p.ProjectID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Get all services
// 	var services sdk.ServiceSlice

// 	for k, v := range p.Service {
// 		new := sdk.Service{
// 			Name: k,
// 			Port: v.Port,
// 		}
// 		if v.TargetPort != nil {
// 			new.TargetPort = v.TargetPort
// 		}
// 		if v.Protocol != nil {
// 			new.Protocol = *v.Protocol
// 		}
// 		if v.DomainName != nil {
// 			new.DomainName = v.DomainName
// 		}
// 		services = append(services, new)
// 	}

// 	// Get all volumes
// 	var volumes sdk.VolumeSlice
// 	for k, v := range p.Volume {
// 		new := sdk.Volume{
// 			Name: k,
// 			Size: v.Size,
// 			Path: v.Path,
// 		}
// 		if v.Class != nil {
// 			new.Class = *v.Class
// 		}
// 		volumes = append(volumes, new)
// 	}

// 	var sharedVolumes sdk.SharedVolumeAttachmentSlice
// 	if p.SharedVolume != nil {
// 		for _, v := range *p.SharedVolume {
// 			new := sdk.SharedVolumeAttachment{
// 				ID:   v.VolumeID,
// 				Path: v.Path,
// 			}
// 			sharedVolumes = append(sharedVolumes, new)
// 		}
// 	}

// 	out := &sdk.Pod{
// 		Project:       theProject,
// 		ID:            p.ID,
// 		Name:          p.Name,
// 		Image:         p.Image,
// 		Command:       p.Command,
// 		Arguments:     p.Arguments,
// 		Services:      services,
// 		Volumes:       volumes,
// 		SharedVolumes: sharedVolumes,
// 	}

// 	if p.Tag != nil {
// 		out.Tag = p.Tag
// 	}

// 	for k, v := range p.Environment {
// 		new := sdk.EnvironmentVariable{
// 			Key:   k,
// 			Value: v.Value,
// 		}
// 		if v.Secret != nil {
// 			new.Secret = *v.Secret
// 		} else {
// 			new.Secret = false
// 		}
// 		out.Environment = append(out.Environment, new)
// 	}

// 	return out, nil
// }
