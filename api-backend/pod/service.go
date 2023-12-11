package pod

import (
	"context"
	"fmt"
	"log"

	"github.com/johncave/podinate/api-backend/config"
	api "github.com/johncave/podinate/api-backend/go"
	corev1 "k8s.io/api/core/v1"
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
	return Service{
		Name:       apiService.Name,
		Port:       int(apiService.Port),
		TargetPort: int(apiService.TargetPort),
		Protocol:   apiService.Protocol,
		DomainName: apiService.DomainName,
	}
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
				Name:      p.ID + "-" + service.Name,
				Namespace: p.getNamespaceName(),
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
						Protocol:   corev1.Protocol(service.Protocol),
					},
				},
			},
		}
	}
	return &services
}

func (p *Pod) ensureServices() error {

	clientset, err := getKubesClient()
	if err != nil {
		log.Printf("error getting kubernetes client: %v\n", err)
		return err
	}

	// TODO - Figure out how to get the service port from the pod
	serviceSpec := p.getServiceSpec()

	if serviceSpec == nil {
		return nil
	}

	// Loop over the services and create them if they don't exist, or update them if they do
	for _, service := range *serviceSpec {
		_, err := clientset.CoreV1().
			Services(p.getNamespaceName()).
			Get(context.Background(), service.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("error getting service: %v\n", err)
			_, err := clientset.CoreV1().
				Services(p.getNamespaceName()).
				Create(context.Background(), &service, metav1.CreateOptions{})
			if err != nil {
				fmt.Printf("error creating service: %v\n", err)
				return err
			}
		} else {
			_, err := clientset.CoreV1().
				Services(p.getNamespaceName()).
				Update(context.Background(), &service, metav1.UpdateOptions{})
			if err != nil {
				fmt.Printf("error updating service: %v\n", err)
				return err
			}
		}
	}

	return nil
}
