package sdk

import "github.com/johncave/podinate/lib/api_client"

type Service struct {
	Name       string  `json:"name"`
	Port       int     `json:"port"`
	TargetPort *int    `json:"target_port"`
	Protocol   string  `json:"protocol"`
	DomainName *string `json:"domain_name"`
}

type ServiceSlice []Service

// servicesFromAPI converts an API Service array to a pod ServiceSlice
func servicesFromAPI(apiServices []api_client.Service) ServiceSlice {
	services := make(ServiceSlice, len(apiServices))
	for i, apiService := range apiServices {
		services[i] = serviceFromAPI(apiService)
	}
	return services
}

// serviceFromAPI converts an API Service to a pod Service
func serviceFromAPI(apiService api_client.Service) Service {
	out := Service{
		Name:     apiService.Name,
		Port:     int(apiService.Port),
		Protocol: apiService.Protocol,
	}

	if apiService.TargetPort != nil {
		tp := int(*apiService.TargetPort)
		out.TargetPort = &tp
	}
	if apiService.DomainName != nil {
		out.DomainName = apiService.DomainName
	}
	return out
}

// servicesToAPI converts a pod ServiceSlice to an API Service array
func servicesToAPI(services ServiceSlice) []api_client.Service {
	apiServices := make([]api_client.Service, len(services))
	for i, service := range services {
		apiServices[i] = serviceToAPI(service)
	}
	return apiServices
}

// serviceToAPI converts a pod Service to an API Service
func serviceToAPI(service Service) api_client.Service {
	out := api_client.Service{
		Name:     service.Name,
		Port:     int32(service.Port),
		Protocol: service.Protocol,
	}

	if service.TargetPort != nil {
		tp := int32(*service.TargetPort)
		out.TargetPort = &tp
	}
	if service.DomainName != nil {
		out.DomainName = service.DomainName
	}
	return out
}
