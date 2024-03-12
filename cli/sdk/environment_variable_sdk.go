package sdk

import "github.com/johncave/podinate/lib/api_client"

type EnvironmentVariable struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Secret bool   `json:"secret,omitempty"`
}

type EnvironmentSlice []EnvironmentVariable

// EnvVarSliceFromAPI converts an API EnvironmentVariable array to a pod EnvironmentSlive
func environmentVariablesFromAPI(apiEnvVars []api_client.EnvironmentVariable) EnvironmentSlice {
	envVars := make(EnvironmentSlice, len(apiEnvVars))
	for i, apiEnvVar := range apiEnvVars {
		envVars[i] = environmentVariableFromAPI(apiEnvVar)
	}
	return envVars
}

// EnvVarFromAPI converts an API EnvironmentVariable to a pod EnvironmentVariable
func environmentVariableFromAPI(apiEnvVar api_client.EnvironmentVariable) EnvironmentVariable {
	out := EnvironmentVariable{
		Key:   apiEnvVar.Key,
		Value: apiEnvVar.Value,
	}

	if apiEnvVar.Secret != nil {
		out.Secret = *apiEnvVar.Secret
	} else {
		out.Secret = false
	}

	// I want to do this but it made errors when trying to communicate with the API
	// if out.Secret {
	// 	out.Value = ""
	// }

	return out
}

// environmentVariablesToAPI converts a pod EnvironmentSlice to an API EnvironmentVariable array
func environmentVariablesToAPI(envVars EnvironmentSlice) []api_client.EnvironmentVariable {
	apiEnvVars := make([]api_client.EnvironmentVariable, len(envVars))
	for i, envVar := range envVars {
		apiEnvVars[i] = environmentVariableToAPI(envVar)
	}
	return apiEnvVars
}

// environmentVariableToAPI converts a pod EnvironmentVariable to an API EnvironmentVariable
func environmentVariableToAPI(envVar EnvironmentVariable) api_client.EnvironmentVariable {
	out := api_client.EnvironmentVariable{
		Key:   envVar.Key,
		Value: envVar.Value,
	}

	if envVar.Secret {
		out.Secret = &envVar.Secret
	}

	return out
}
