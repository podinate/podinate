package engine

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	stderrors "errors"

	cmclient "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	hcldec "github.com/hashicorp/hcl2/hcldec"
	helpers "github.com/podinate/podinate/engine/helpers"
	"github.com/podinate/podinate/kube_client"
	"github.com/podinate/podinate/tui"
	"github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

const (
	KubernetesKindStatefulSet              = "StatefulSet"
	KubernetesAPIVersionStatefulSet        = "apps/v1"
	KubernetesKindService                  = "Service"
	KubernetesAPIVersionService            = "v1"
	KubernetesKindIngress                  = "Ingress"
	KubernetesAPIVersionIngress            = "networking.k8s.io/v1"
	PodinateDefaultClusterIssuerAnnotation = "podinate.com/default-cluster-issuer"
)

type Pod struct {
	ID string
	// Name        string   `cty:"name"`
	Namespace   *string  `cty:"namespace"`
	Image       string   `cty:"image"`
	Tag         *string  `cty:"tag"`
	Command     []string `cty:"command"`
	Arguments   []string `cty:"arguments"`
	Environment map[string]struct {
		Value  string `cty:"value"`
		Secret *bool  `cty:"secret"`
	} `cty:"environment"`
	Services map[string]struct {
		Port       int     `cty:"port"`
		TargetPort *int    `cty:"target_port"`
		Protocol   *string `cty:"protocol"`
		Ingress    *struct {
			HostName      string  `cty:"hostname"`
			Path          *string `cty:"path"`
			IngressClass  *string `cty:"ingress_class"`
			TLS           *bool   `cty:"tls"`
			ClusterIssuer *string `cty:"cluster_issuer"`
		} `cty:"ingress"`
	} `cty:"service"`

	// Removed for now, because these are separate Kube resources
	Volume map[string]struct {
		Size      string  `cty:"size"`
		MountPath string  `cty:"mount_path"`
		Class     *string `cty:"class"`
	} `cty:"volume"`
	SharedVolume map[string]struct {
		MountPath string `cty:"mount_path"`
	} `cty:"shared_volume"`
}

// GetHCLSpect returns the HCL spec of the pod type
var podHCLSpec = &hcldec.BlockMapSpec{
	TypeName:   "pod",
	LabelNames: []string{"id"},
	Nested: &hcldec.ObjectSpec{
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
		"namespace": &hcldec.AttrSpec{
			Name:     "namespace",
			Type:     cty.String,
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
		"service": &hcldec.BlockMapSpec{
			TypeName:   "service",
			LabelNames: []string{"name"},
			Nested: &hcldec.ObjectSpec{
				// The listening port of the service
				// For example an app might listen on port 3000 for http traffic
				// In that case we set port to 80 and target_port to 3000
				"port": &hcldec.AttrSpec{
					Name:     "port",
					Type:     cty.Number,
					Required: true,
				},
				// The service on the container that the port is forwarded to
				"target_port": &hcldec.AttrSpec{
					Name:     "target_port",
					Type:     cty.Number,
					Required: false,
				},
				// Protocol of the service. Defaults to TCP
				// Can be set to UDP
				"protocol": &hcldec.AttrSpec{
					Name:     "protocol",
					Type:     cty.String,
					Required: false,
				},
				// make the service externally available
				// If protocol is set to "http" or "https", an Ingress will be created
				// If a domain name followed by a path, the service will be available at that path on the ingress
				"ingress": &hcldec.BlockSpec{
					TypeName: "ingress",
					Required: false,
					Nested: &hcldec.ObjectSpec{
						"hostname": &hcldec.AttrSpec{
							Name:     "hostname",
							Type:     cty.String,
							Required: true,
						},
						"path": &hcldec.AttrSpec{
							Name:     "path",
							Type:     cty.String,
							Required: false,
						},
						"ingress_class": &hcldec.AttrSpec{
							Name:     "ingress_class",
							Type:     cty.String,
							Required: false, // If not set, default class is used
						},
						"tls": &hcldec.AttrSpec{
							Name:     "tls",
							Type:     cty.Bool,
							Required: false,
						},
						"cluster_issuer": &hcldec.AttrSpec{
							Name:     "cluster_issuer",
							Type:     cty.String,
							Required: false,
						},
					},
				},
			},
		},
		"volume": &hcldec.BlockMapSpec{
			TypeName:   "volume",
			LabelNames: []string{"name"},
			Nested: &hcldec.ObjectSpec{
				"size": &hcldec.AttrSpec{
					Name:     "size",
					Type:     cty.String,
					Required: true,
				},
				"mount_path": &hcldec.AttrSpec{
					Name:     "mount_path",
					Type:     cty.String,
					Required: true,
				},
				"class": &hcldec.AttrSpec{
					Name:     "class",
					Type:     cty.String,
					Required: false,
				},
			},
		},
		"shared_volume": &hcldec.BlockMapSpec{
			TypeName:   "shared_volume",
			LabelNames: []string{"volume_id"},
			Nested: &hcldec.ObjectSpec{
				"mount_path": &hcldec.AttrSpec{
					Name:     "mount_path",
					Type:     cty.String,
					Required: true,
				},
			},
		},
	},
}

// GetResources returns the resources needed for the Pod
// This function deals with the StatefulSet, then calls other functions to add services, shared volumes, etc.
func (p *Pod) GetResources(ctx context.Context, pkg *Package) (*ChangeType, []ResourceChange, error) {
	var out []ResourceChange
	podChangeType := ChangeTypeNoop

	var imageID = p.Image
	if p.Tag != nil {
		imageID = imageID + ":" + *p.Tag
	}

	ssSpec := &appsv1.StatefulSet{

		ObjectMeta: metav1.ObjectMeta{
			Name:      p.ID,
			Namespace: pkg.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			PersistentVolumeClaimRetentionPolicy: func(val appsv1.StatefulSetPersistentVolumeClaimRetentionPolicy) *appsv1.StatefulSetPersistentVolumeClaimRetentionPolicy {
				return &val
			}(appsv1.StatefulSetPersistentVolumeClaimRetentionPolicy{WhenDeleted: appsv1.RetainPersistentVolumeClaimRetentionPolicyType, WhenScaled: appsv1.RetainPersistentVolumeClaimRetentionPolicyType}),
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

	ssSpec.Kind = KubernetesKindStatefulSet
	ssSpec.APIVersion = KubernetesAPIVersionStatefulSet

	// Add environment variables
	for k, v := range p.Environment {
		ssSpec.Spec.Template.Spec.Containers[0].Env = append(ssSpec.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{
			Name:  k,
			Value: v.Value,
		})
	}

	// Add service ports to the StatefulSet
	// This port is just an annotation to help admins know what the container does
	for k, v := range p.Services {
		//ssSpec.Spec.ServiceName = k
		port := corev1.ContainerPort{
			Name:          k,
			ContainerPort: int32(v.Port),
		}

		if v.TargetPort != nil {
			port.ContainerPort = int32(*v.TargetPort)
		}
		ssSpec.Spec.Template.Spec.Containers[0].Ports = append(ssSpec.Spec.Template.Spec.Containers[0].Ports, port)

	}

	// Add volumes to the StatefulSet
	for k, volume := range p.Volume {
		size, err := resource.ParseQuantity(volume.Size)
		if err != nil {
			return nil, nil, err
		}
		ssSpec.Spec.Template.Spec.Containers[0].VolumeMounts = append(ssSpec.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
			Name:      k,
			MountPath: volume.MountPath,
		})
		// Add volume claim templates
		newPVC := corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      k,
				Namespace: *p.Namespace,
				Annotations: map[string]string{
					"volumeType": "local",
				},
				Labels: map[string]string{
					"podinate.com/pod": p.ID,
				},
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					"ReadWriteOnce",
				},
				//StorageClassName: func(val string) *string { return &val }("local-path"),
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						"storage": size,
					},
				},
				VolumeMode: func(val corev1.PersistentVolumeMode) *corev1.PersistentVolumeMode { return &val }(corev1.PersistentVolumeFilesystem),
			},
		}
		newPVC.Kind = "PersistentVolumeClaim"
		newPVC.APIVersion = "v1"

		if volume.Class != nil {
			// TODO: Check if the SC exists
			newPVC.Spec.StorageClassName = volume.Class
		}

		ssSpec.Spec.VolumeClaimTemplates = append(ssSpec.Spec.VolumeClaimTemplates, newPVC)
	}

	// Add shared volumes to the StatefulSet
	for k, sharedVolume := range p.SharedVolume {
		// Specify which volume we're talking about
		ssSpec.Spec.Template.Spec.Volumes = append(ssSpec.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: k,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: k,
				},
			},
		})

		// Specify where to mount it in the container
		ssSpec.Spec.Template.Spec.Containers[0].VolumeMounts = append(ssSpec.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
			Name:      k,
			MountPath: sharedVolume.MountPath,
		})
	}

	// Add service changes
	serviceChangeType, svcChanges, err := p.GetServiceResourceChanges(ctx)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"error":           err,
			"service_changes": svcChanges,
			"service_ct":      serviceChangeType,
		}).Error("Error getting service changes")
		return nil, nil, err
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{
		"service_changes": svcChanges,
		"service_ct":      serviceChangeType,
	}).Trace("Service changes")

	//changeType = *ct
	out = append(out, svcChanges...)

	// Check if the StatefulSet already exists
	kube, err := kube_client.Client()
	if err != nil {
		return nil, nil, err
	}

	ss, err := kube.AppsV1().StatefulSets(pkg.Namespace).Update(ctx, ssSpec, metav1.UpdateOptions{DryRun: []string{metav1.DryRunAll}})
	ss.Kind = KubernetesKindStatefulSet
	ss.APIVersion = KubernetesAPIVersionStatefulSet

	logrus.WithContext(ctx).WithFields(logrus.Fields{
		"statefulset_from_dry_run": ss,
		"kind":                     ss.GetObjectKind(),
		"statefulset_generated":    ssSpec,
		"error":                    err,
	}).Debug("StatefulSet")

	// If the StatefulSet didn't exist when we called update, create it
	if errors.IsNotFound(err) {

		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"statefulset": ss,
		}).Debug("StatefulSet needs to be created")

		// Resource doesn't exist, so create it
		podChangeType = ChangeTypeCreate

		// Dry Run creation to fill out default fields
		// Disabled for now, fails if the Namespace isn't already created
		// ss, err := kube.AppsV1().StatefulSets(pkg.Namespace).Create(ctx, ssSpec, metav1.CreateOptions{DryRun: []string{metav1.DryRunAll}})
		// if err != nil {
		// 	return nil, nil, err
		// }

		out = append(out, ResourceChange{
			ChangeType:      ChangeTypeCreate,
			CurrentResource: nil,
			DesiredResource: ssSpec,
		})
	} else if err != nil {
		if errors.IsInvalid(err) { // Kubernetes rejected our change
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"error": err,
			}).Error("Podinate is unable to make this change. Sorry but this is a limitation of the Kubernetes API. Most likely you're trying to update the volumes attached directly to a Podinate Pod. If you just want to change the size, run 'kubectl -n <namespace> get pvc', then 'kubectl edit pvc <pvc-name>' and change the size there. Then, go into your pod definition and change the size there too.")

			// Grab the existing spec so we can show the user what we're trying to change
			existing, geterr := kube.AppsV1().StatefulSets(pkg.Namespace).Get(ctx, p.ID, metav1.GetOptions{})
			if geterr != nil {
				logrus.WithContext(ctx).WithFields(logrus.Fields{
					"error": geterr,
				}).Error("Error getting existing StatefulSet for comparison")
				return nil, nil, err
			}
			existing.Kind = KubernetesKindStatefulSet
			existing.APIVersion = KubernetesAPIVersionStatefulSet

			fmt.Println(tui.StyleError.Render("The following change was rejected by Kubernetes:"))

			err := helpers.YamlDiffObjects(ctx, existing, ssSpec)
			if err != nil {
				logrus.WithContext(ctx).WithFields(logrus.Fields{
					"error": err,
				}).Error("A further error occurred when trying to display the change that was rejected by the Kubernetes API.")
			}
		} else {
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"error": err,
			}).Error("Unknown error checking StatefulSet specification")
		}
		return nil, nil, err
	} else {
		existing, err := kube.AppsV1().StatefulSets(pkg.Namespace).Get(ctx, p.ID, metav1.GetOptions{})
		if err != nil {
			return nil, nil, err
		}
		existing.Kind = KubernetesKindStatefulSet
		existing.APIVersion = KubernetesAPIVersionStatefulSet

		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"statefulset": ss,
		}).Debug("StatefulSet exists, deep comparing")
		// Deep compare the object to see if it needs updating
		if reflect.DeepEqual(ss.Spec, existing.Spec) {
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"statefulset": ss.Spec,
				"existing":    existing.Spec,
			}).Debug("StatefulSet is up to date")
			//changeType = ChangeTypeNoop
		} else {
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"statefulset": ss,
				"new-kind":    ss.GetObjectKind(),
				"old-kind":    existing.GetObjectKind(),
				"existing":    existing,
			}).Debug("StatefulSet needs updating")

			// Resource exists, but needs updating
			out = append(out, ResourceChange{
				ChangeType:      ChangeTypeUpdate,
				CurrentResource: existing,
				DesiredResource: ss,
			})
			podChangeType = ChangeTypeUpdate
		}
	}

	// Decide the overall change type
	changeType := ChangeTypeNoop
	if podChangeType == ChangeTypeCreate || podChangeType == ChangeTypeUpdate {
		changeType = podChangeType
	} else if (*serviceChangeType == ChangeTypeCreate || *serviceChangeType == ChangeTypeUpdate) && podChangeType == ChangeTypeNoop {
		changeType = ChangeTypeUpdate
	}

	return &changeType, out, nil
}

func (p *Pod) GetServiceResourceChanges(ctx context.Context) (*ChangeType, []ResourceChange, error) {
	var out []ResourceChange
	changeType := ChangeTypeNoop

	for name, service := range p.Services {
		svcSpec := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: *p.Namespace,
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{
					"podinate.com/pod": p.ID,
				},
				Ports: []corev1.ServicePort{
					{
						Port:     int32(service.Port),
						Protocol: corev1.ProtocolTCP,
					},
				},
			},
		}

		svcSpec.Kind = KubernetesKindService
		svcSpec.APIVersion = KubernetesAPIVersionService

		if service.Protocol != nil && strings.ToLower(*service.Protocol) == "udp" {
			svcSpec.Spec.Ports[0].Protocol = corev1.ProtocolUDP
		}

		if service.TargetPort != nil {
			svcSpec.Spec.Ports[0].TargetPort = intstr.FromInt(*service.TargetPort)
		}

		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"service": svcSpec,
		}).Trace("generated service spec")

		rc, err := GetResourceChangeForResource(ctx, svcSpec)
		if errors.IsInvalid(err) {
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"resource_type": ResourceTypeService,
				"resource_id":   name,
				"error":         err,
			}).Debug("Got invalid error trying to get resource change for service")
			return nil, nil, err
		} else if err != nil {
			return nil, nil, err
		}

		// Check if the service needs an Ingress
		// Disable ingress for now
		ct, ingressChanges, err := p.GetServiceIngressChanges(ctx, name)
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"error": err,
			}).Error("Error getting Ingress changes")
			return nil, nil, err
		}
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"ingress_changes": ingressChanges,
			"changeType":      changeType,
		}).Trace("Ingress changes")
		if ct != nil {
			changeType = *ct
		}

		out = append(out, ingressChanges...)

		if rc != nil {
			out = append(out, *rc)
			if rc.ChangeType == ChangeTypeCreate || rc.ChangeType == ChangeTypeUpdate {
				changeType = ChangeTypeUpdate
			}

		}
	}

	logrus.WithContext(ctx).WithFields(logrus.Fields{
		"changeType":       changeType,
		"resource_changes": out,
	}).Debug("GetServiceResourceChanges")

	return &changeType, out, nil
}

func (pod *Pod) GetServiceIngressChanges(ctx context.Context, serviceName string) (*ChangeType, []ResourceChange, error) {
	// var out []ResourceChange
	// changeType := ChangeTypeNoop

	if pod.Services[serviceName].Ingress == nil {
		// No ingress needed
		return nil, nil, nil
	}

	client, err := kube_client.Client()
	if err != nil {
		return nil, nil, err
	}

	ingressRequest := pod.Services[serviceName].Ingress

	// Validate and add defaults
	if ingressRequest.Path == nil {
		ingressRequest.Path = func(val string) *string { return &val }("/")
	}

	if ingressRequest.IngressClass == nil {
		ic, err := getDefaultIngressClass(ctx, client)
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"error": err,
			}).Error("Error getting default IngressClass")
			return nil, nil, err
		}
		ingressRequest.IngressClass = ic
	}

	ingressSpec := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: *pod.Namespace,
			Labels: map[string]string{
				"podinate.com/pod": pod.ID,
			},
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: ingressRequest.IngressClass,
			Rules: []networkingv1.IngressRule{
				{
					Host: ingressRequest.HostName,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									PathType: func(val networkingv1.PathType) *networkingv1.PathType { return &val }(networkingv1.PathTypePrefix),
									Path:     *ingressRequest.Path,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: serviceName,
											Port: networkingv1.ServiceBackendPort{
												Number: int32(pod.Services[serviceName].Port),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	ingressSpec.Kind = KubernetesKindIngress
	ingressSpec.APIVersion = KubernetesAPIVersionIngress

	// Check if the url protocol is https, and add tls annotations if so
	if ingressRequest.TLS != nil && *ingressRequest.TLS {
		if ingressSpec.Annotations == nil {
			ingressSpec.Annotations = make(map[string]string)
		}

		ingressSpec.Annotations["nginx.ingress.kubernetes.io/ssl-redirect"] = "true"
		ingressSpec.Spec.TLS = []networkingv1.IngressTLS{
			{
				Hosts: []string{
					ingressRequest.HostName,
				},
				SecretName: serviceName + "-tls",
			},
		}

		if ingressRequest.ClusterIssuer == nil {
			ci, err := getDefaultClusterIssuer(ctx, client)
			if err != nil {
				logrus.WithContext(ctx).WithFields(logrus.Fields{
					"error": err,
				}).Error("Could not find a default ClusterIssuer. Either specify a ClusterIssuer in the Pod > Service > Ingress definition or create a default ClusterIssuer in the cluster by adding the annotation '" + PodinateDefaultClusterIssuerAnnotation + ": true' to a ClusterIssuer. See https://cert-manager.io/docs/concepts/issuer/ and https://docs.podinate.com/kubernetes/certificates/ for more information.")
				return nil, nil, err
			}

			ingressSpec.Annotations["cert-manager.io/cluster-issuer"] = *ci

		} else {
			ok, err := checkClusterIssuer(ctx, client, *ingressRequest.ClusterIssuer)
			if err != nil {
				return nil, nil, err
			}
			if !ok {
				logrus.WithContext(ctx).WithFields(logrus.Fields{
					"cluster_issuer": *ingressRequest.ClusterIssuer,
				}).Error("ClusterIssuer not found. To list all available ClusterIssuers, run 'kubectl get clusterissuers'. If you can't list ClusterIssuers, you may need to install Cert-Manager. See https://cert-manager.io/docs/concepts/issuer/ and https://docs.podinate.com/kubernetes/certificates/ for more information.")
				return nil, nil, stderrors.New("ClusterIssuer not found")
			}
			ingressSpec.Annotations["cert-manager.io/cluster-issuer"] = *ingressRequest.ClusterIssuer
		}

	}

	rc, err := GetResourceChangeForResource(ctx, ingressSpec)
	if err != nil {
		return nil, nil, err
	}
	if rc == nil {
		return nil, nil, nil
	}

	return &rc.ChangeType, []ResourceChange{*rc}, nil
}

// getDefaultIngressClass gets the default IngressClass from the cluster
func getDefaultIngressClass(ctx context.Context, client *kubernetes.Clientset) (*string, error) {
	// Get the IngressClass from the cluster

	ingressClasses, err := client.NetworkingV1().IngressClasses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	logrus.WithContext(ctx).WithFields(logrus.Fields{
		"ingressClasses": ingressClasses,
	}).Debug("Got IngressClasses")

	// Find the default IngressClass
	for _, ic := range ingressClasses.Items {
		if ic.ObjectMeta.Annotations["ingressclass.kubernetes.io/is-default-class"] == "true" {
			return &ic.Name, nil
		}
	}

	return nil, stderrors.New("No default IngressClass found")
}

// getDefaultClusterIssuer gets the default ClusterIssuer from the cluster
func getDefaultClusterIssuer(ctx context.Context, client *kubernetes.Clientset) (*string, error) {
	// Get the ClusterIssuers from the cluster

	//cm := cmclient.New(client.RESTClient())
	rc, err := kube_client.GetRestConfig()
	if err != nil {
		return nil, err
	}
	cm, err := cmclient.NewForConfig(rc)
	if err != nil {
		return nil, err
	}

	clusterIssuers, err := cm.CertmanagerV1().ClusterIssuers().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	logrus.WithContext(ctx).WithFields(logrus.Fields{
		"clusterIssuers": clusterIssuers,
	}).Debug("Got ClusterIssuers")

	// Find the default ClusterIssuer
	for _, ci := range clusterIssuers.Items {
		if ci.ObjectMeta.Annotations[PodinateDefaultClusterIssuerAnnotation] == "true" {
			return &ci.Name, nil
		}
	}

	return nil, stderrors.New("No default ClusterIssuer found")
}

// checkClusterIssuer checks if the ClusterIssuer exists in the cluster
func checkClusterIssuer(ctx context.Context, client *kubernetes.Clientset, ci string) (bool, error) {
	rc, err := kube_client.GetRestConfig()
	if err != nil {
		return false, err
	}
	cm, err := cmclient.NewForConfig(rc)
	if err != nil {
		return false, err
	}

	_, err = cm.CertmanagerV1().ClusterIssuers().Get(ctx, ci, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
