package engine

import (
	"context"
	"strings"

	stderrors "errors"

	cmclient "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	hcldec "github.com/hashicorp/hcl/v2/hcldec"
	"github.com/podinate/podinate/kube_client"
	"github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
	Namespace   *string           `cty:"namespace"`
	Image       string            `cty:"image"`
	Tag         *string           `cty:"tag"`
	Command     []string          `cty:"command"`
	Arguments   []string          `cty:"arguments"`
	Environment map[string]string `cty:"environment"`
	Resource    map[string]struct {
		Limits   *string `cty:"limits"`
		Requests *string `cty:"requests"`
	} `cty:"resource"`
	Services map[string]struct {
		Port       int     `cty:"port"`
		TargetPort *int    `cty:"target_port"`
		Protocol   *string `cty:"protocol"`
		Type       *string `cty:"type"`
		Ingress    *struct {
			Annotation map[string]struct {
				Value string `cty:"value"`
			} `cty:"annotation"`
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
		"environment": &hcldec.BlockAttrsSpec{
			TypeName:    "environment",
			ElementType: cty.String,
			Required:    false,
		},
		"resource": &hcldec.BlockMapSpec{
			TypeName:   "resource",
			LabelNames: []string{"key"},
			Nested: &hcldec.ObjectSpec{
				"limits": &hcldec.AttrSpec{
					Name:     "limits",
					Type:     cty.String,
					Required: false,
				},
				"requests": &hcldec.AttrSpec{
					Name:     "requests",
					Type:     cty.String,
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
				// The Service type. Defaults to ClusterIP
				// Can be set to NodePort, LoadBalancer, or ExternalName
				"type": &hcldec.AttrSpec{
					Name:     "type",
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
						"annotation": &hcldec.BlockMapSpec{
							TypeName:   "annotation",
							LabelNames: []string{"key"},
							Nested: &hcldec.ObjectSpec{
								"value": &hcldec.AttrSpec{
									Name:     "value",
									Type:     cty.String,
									Required: false,
								},
							},
						},
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

// Implementing the Resource interface
func (p Pod) GetName() string {
	return p.ID
}

func (p Pod) GetType() ResourceType {
	return ResourceTypePod
}

// GetObjects returns the Kubernetes objects needed for the Pod
func (p Pod) GetObjects(ctx context.Context) ([]runtime.Object, error) {
	var out []runtime.Object

	var imageID = p.Image
	if p.Tag != nil {
		imageID = imageID + ":" + *p.Tag
	}

	ssSpec := &appsv1.StatefulSet{

		ObjectMeta: metav1.ObjectMeta{
			Name:      p.ID,
			Namespace: *p.Namespace,
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
							Resources: corev1.ResourceRequirements{
								Limits:   corev1.ResourceList{},
								Requests: corev1.ResourceList{},
							},
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
			Value: v,
		})
	}

	for k, v := range p.Resource {
		if v.Limits != nil {
			limitsQuantity, err := resource.ParseQuantity(*v.Limits)
			if err != nil {
				return nil, err
			}
			ssSpec.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceName(k)] = limitsQuantity
		}
		if v.Requests != nil {
			requestsQuantity, err := resource.ParseQuantity(*v.Requests)
			if err != nil {
				return nil, err
			}
			ssSpec.Spec.Template.Spec.Containers[0].Resources.Requests[corev1.ResourceName(k)] = requestsQuantity
		}

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
			return nil, err
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

	out = append(out, ssSpec)

	// Add services
	svcObjects, err := p.GetServiceObjects(ctx)
	if err != nil {
		return nil, err
	}
	out = append(out, svcObjects...)

	return out, nil
}

// GetServiceObjects returns the Kubernetes objects needed for the services
func (p *Pod) GetServiceObjects(ctx context.Context) ([]runtime.Object, error) {
	var out []runtime.Object
	for name, service := range p.Services {
		if len(name) > 15 {
			return nil, stderrors.New("service names must be no more than 15 characters")
		}
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

		// Set the service type
		if service.Type != nil {
			svcSpec.Spec.Type = corev1.ServiceType(*service.Type)
		} else {
			svcSpec.Spec.Type = corev1.ServiceTypeClusterIP
		}

		if svcSpec.Spec.Type == corev1.ServiceTypeNodePort {
			svcSpec.Spec.Ports[0].NodePort = int32(service.Port)
		}

		if service.Protocol != nil && strings.ToLower(*service.Protocol) == "udp" {
			svcSpec.Spec.Ports[0].Protocol = corev1.ProtocolUDP
		}

		if service.TargetPort != nil {
			svcSpec.Spec.Ports[0].TargetPort = intstr.FromInt(*service.TargetPort)
		}

		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"service": svcSpec,
		}).Trace("generated service spec")

		// Check if the service needs an Ingress
		// Disable ingress for now
		ingressObjects, err := p.GetServiceIngressObjects(ctx, name)
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"error": err,
			}).Error("Error getting Ingress changes")
			return nil, err
		}
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"ingress_changes": ingressObjects,
		}).Trace("Ingress changes")

		out = append(out, svcSpec)
		out = append(out, ingressObjects...)

	}
	return out, nil
}

// GetServiceIngressObjects returns the Kubernetes objects needed for the Ingresses of the services
func (pod *Pod) GetServiceIngressObjects(ctx context.Context, serviceName string) ([]runtime.Object, error) {
	if pod.Services[serviceName].Ingress == nil {
		// No ingress needed
		return nil, nil
	}

	client, err := kube_client.Client()
	if err != nil {
		return nil, err
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
			return nil, err
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
				return nil, err
			}

			ingressSpec.Annotations["cert-manager.io/cluster-issuer"] = *ci

		} else {
			ok, err := checkClusterIssuer(ctx, client, *ingressRequest.ClusterIssuer)
			if err != nil {
				return nil, err
			}
			if !ok {
				logrus.WithContext(ctx).WithFields(logrus.Fields{
					"cluster_issuer": *ingressRequest.ClusterIssuer,
				}).Error("ClusterIssuer not found. To list all available ClusterIssuers, run 'kubectl get clusterissuers'. If you can't list ClusterIssuers, you may need to install Cert-Manager. See https://cert-manager.io/docs/concepts/issuer/ and https://docs.podinate.com/kubernetes/certificates/ for more information.")
				return nil, stderrors.New("ClusterIssuer not found")
			}
			ingressSpec.Annotations["cert-manager.io/cluster-issuer"] = *ingressRequest.ClusterIssuer
		}

	}

	for k, v := range ingressRequest.Annotation {
		ingressSpec.Annotations[k] = v.Value
	}
	return []runtime.Object{ingressSpec}, nil
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
