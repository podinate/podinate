package apiclient

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

	if out.Secret {
		out.Value = ""
	}

	return out
}
