package pod

import (
	"context"
	"log"
	"strings"

	"github.com/johncave/podinate/controller/apierror"
	"github.com/johncave/podinate/controller/config"
	api "github.com/johncave/podinate/controller/go"
	lh "github.com/johncave/podinate/controller/loghandler"
	corev1 "k8s.io/api/core/v1"
	v1networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Service struct {
	Name       string
	Port       int
	TargetPort int
	Protocol   string
	DomainName string
}

type ServiceSlice []Service

// loadServices returns the services for a pod
func (p *Pod) loadServices() error {
	rows, err := config.DB.Query("SELECT name, port, target_port, protocol, domain_name FROM pod_services WHERE pod_uuid = $1", p.Uuid)
	if err != nil {
		log.Println("Error getting pod services", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var s Service
		err := rows.Scan(&s.Name, &s.Port, &s.TargetPort, &s.Protocol, &s.DomainName)
		if err != nil {
			log.Println("Error scanning pod services", err)
			return err
		}
		p.Services = append(p.Services, s)
	}
	return nil
}

// servicesFromAPI converts an API Service array to a pod ServiceSlice
func servicesFromAPI(apiServices []api.Service) ServiceSlice {
	services := make(ServiceSlice, len(apiServices))
	for i, apiService := range apiServices {
		services[i] = serviceFromAPI(apiService)
	}
	return services
}

// serviceFromAPI converts an API Service to a pod Service
func serviceFromAPI(apiService api.Service) Service {
	if apiService.Protocol == "" {
		apiService.Protocol = "tcp"
	}
	out := Service{
		Name:       apiService.Name,
		Port:       int(apiService.Port),
		TargetPort: int(apiService.TargetPort),
		Protocol:   apiService.Protocol,
		DomainName: apiService.DomainName,
	}
	return out
}

// servicesToAPI converts a pod ServiceSlice to an API Service array
func ServicesToAPI(services ServiceSlice) []api.Service {
	apiServices := make([]api.Service, len(services))
	for i, service := range services {
		apiServices[i] = serviceToAPI(service)
	}
	return apiServices
}

// serviceToAPI converts a pod Service to an API Service
func serviceToAPI(service Service) api.Service {
	return api.Service{
		Name:       service.Name,
		Port:       int32(service.Port),
		TargetPort: int32(service.TargetPort),
		Protocol:   service.Protocol,
		DomainName: service.DomainName,
	}
}

// getServiceSpec returns the kubernetes service spec for a p
func (p *Pod) getServiceSpec() *[]corev1.Service {
	if len(p.Services) == 0 {
		return nil
	}
	services := make([]corev1.Service, len(p.Services))
	for i, service := range p.Services {
		if service.Protocol == "http" {
			service.Protocol = "TCP"
		}
		services[i] = corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      service.Name,
				Namespace: p.getNamespaceName(),
				Labels: map[string]string{
					"podinate.com/project":     p.Project.ID,
					"podinate.com/pod":         p.ID,
					"podinate.com/domain-name": service.DomainName,
				},
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{
					"podinate.com/pod":     p.ID,
					"podinate.com/project": p.Project.ID,
				},
				Ports: []corev1.ServicePort{
					{
						Name:       service.Name,
						Port:       int32(service.Port),
						TargetPort: intstr.FromInt(service.TargetPort),
						Protocol:   corev1.Protocol(strings.ToUpper(service.Protocol)),
					},
				},
			},
		}
	}
	return &services
}

// getIngressSpec returns the kubernetes ingress spec for a pod
func (p *Pod) getIngressSpec() *[]v1networking.Ingress {
	if len(p.Services) == 0 {
		return nil
	}
	ingresses := make([]v1networking.Ingress, 0)
	for _, service := range p.Services {
		//lh.Log.Debugw("getIngressSpec", "service", service)
		if service.DomainName == "" {
			continue
		}
		if service.Protocol != "http" {
			continue
		}
		new := v1networking.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:        service.Name,
				Namespace:   p.getNamespaceName(),
				Annotations: p.getAnnotations(),
			},
			Spec: v1networking.IngressSpec{
				Rules: []v1networking.IngressRule{
					{
						Host: service.DomainName,
						IngressRuleValue: v1networking.IngressRuleValue{
							HTTP: &v1networking.HTTPIngressRuleValue{
								Paths: []v1networking.HTTPIngressPath{
									{
										Path:     "/",
										PathType: func() *v1networking.PathType { p := v1networking.PathTypePrefix; return &p }(),
										Backend: v1networking.IngressBackend{
											Service: &v1networking.IngressServiceBackend{
												Name: service.Name,
												Port: v1networking.ServiceBackendPort{
													Number: int32(service.Port),
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
		ingresses = append(ingresses, new)
	}

	lh.Log.Infow("getIngressSpec", "ingresses", ingresses)

	return &ingresses
}

// ensureIngresses ensures that the ingresses for a pod exist
func (p *Pod) ensureIngresses(ctx context.Context) error {

	clientset, err := getKubesClient()
	if err != nil {
		log.Printf("error getting kubernetes client: %v\n", err)
		return err
	}

	// TODO - Figure out how to get the service port from the pod
	ingressSpec := p.getIngressSpec()

	if ingressSpec == nil {
		return nil
	}

	// Loop over the ingresses and create them if they don't exist, or update them if they do
	for _, ingress := range *ingressSpec {
		_, err := clientset.NetworkingV1().
			Ingresses(p.getNamespaceName()).
			Get(context.Background(), ingress.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			lh.Log.Infow("error getting ingress: %v\n", err)
			_, err := clientset.NetworkingV1().
				Ingresses(p.getNamespaceName()).
				Create(context.Background(), &ingress, metav1.CreateOptions{})
			if err != nil {
				lh.Log.Errorw("error creating ingress", "err", err, "ingres_object", ingress)
				return err
			}
			lh.Info(ctx, "Created ingress", "ingress", ingress)
		} else {
			_, err := clientset.NetworkingV1().
				Ingresses(p.getNamespaceName()).
				Update(context.Background(), &ingress, metav1.UpdateOptions{})
			if err != nil {
				lh.Log.Errorw("error updating ingress: %v\n", err)
				return err
			}
			lh.Info(ctx, "Updated ingress", "ingress", ingress)
		}

		// If the service has a domain name, add an ingress for it

	}

	return nil
}

// ensureServices ensures that the services for a pod exist
func (p *Pod) ensureServices(ctx context.Context) *apierror.ApiError {

	clientset, err := getKubesClient()
	if err != nil {
		log.Printf("error getting kubernetes client: %v\n", err)
		return apierror.New(500, "error getting kubernetes client: "+err.Error())
	}

	// TODO - Figure out how to get the service port from the pod
	serviceSpec := p.getServiceSpec()

	if serviceSpec == nil {
		lh.Info(ctx, "No services for pod", "pod", p)
		return nil
	}

	// Loop over the services and create them if they don't exist, or update them if they do
	for _, service := range *serviceSpec {
		_, err := clientset.CoreV1().
			Services(p.getNamespaceName()).
			Get(context.Background(), service.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			lh.Debug(ctx, "Couldn't get service, creating", "error", err)
			_, err := clientset.CoreV1().
				Services(p.getNamespaceName()).
				Create(context.Background(), &service, metav1.CreateOptions{})
			if err != nil {
				lh.Error(ctx, "error creating service, aborting", "error", err)
				return apierror.New(500, "error creating service: "+err.Error())
			}
			lh.Info(ctx, "Created service", "service", service.ObjectMeta.Name)
		} else {
			_, err := clientset.CoreV1().
				Services(p.getNamespaceName()).
				Update(context.Background(), &service, metav1.UpdateOptions{})
			if err != nil {
				lh.Error(ctx, "error updating service, aborting", "error", err)
				return apierror.New(500, "error updating service: "+err.Error())
			}
			lh.Info(ctx, "Updated service", "service", service.ObjectMeta.Name)
		}

		// If the service has a domain name, add an ingress for it

	}

	err = p.ensureIngresses(ctx)
	if err != nil {
		return apierror.New(500, "error ensuring ingresses: "+err.Error())
	}

	lh.Info(ctx, "Created services for pod", "pod", p)
	return nil
}

func (p *Pod) serviceExists(serviceName string) (bool, error) {
	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM pod_services WHERE pod_uuid = $1 AND name = $2", p.Uuid, serviceName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
