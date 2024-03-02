package apiclient

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
