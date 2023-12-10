package pod

import (
	"database/sql/driver"
	"encoding/json"

	api "github.com/johncave/podinate/api-backend/go"
)

// Type EnvironmentVairable is used to store environment variables for a pod
type EnvironmentVariable struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Secret bool   `json:"secret,omitempty"`
}

type EnvironmentSlice []EnvironmentVariable

// Scan implements the sql.Scanner interface
func (e *EnvironmentSlice) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), e)
}

func (e EnvironmentSlice) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e EnvironmentSlice) MarshalJSON() ([]byte, error) {
	return json.Marshal([]EnvironmentVariable(e))
}

func (e *EnvironmentSlice) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*[]EnvironmentVariable)(e))
}

// FromAPI converts an API EnvironmentVariable array to a pod EnvironmentSlive
func EnvVarSliceFromAPI(apiEnvVars []api.EnvironmentVariable) EnvironmentSlice {
	envVars := make(EnvironmentSlice, len(apiEnvVars))
	for i, apiEnvVar := range apiEnvVars {
		envVars[i] = EnvVarFromAPI(apiEnvVar)
	}
	return envVars
}

// func FromAPI converts an API EnvironmentVariable to a pod EnvironmentVariable
func EnvVarFromAPI(apiEnvVar api.EnvironmentVariable) EnvironmentVariable {
	return EnvironmentVariable{
		Key:    apiEnvVar.Key,
		Value:  apiEnvVar.Value,
		Secret: apiEnvVar.Secret,
	}
}

// FromAPIMany converts an array of API EnvironmentVariables to an array of pod EnvironmentVariables
func EnvVarFromAPIMany(apiEnvVars []api.EnvironmentVariable) []EnvironmentVariable {
	envVars := make([]EnvironmentVariable, len(apiEnvVars))
	for i, apiEnvVar := range apiEnvVars {
		envVars[i] = EnvVarFromAPI(apiEnvVar)
	}
	return envVars
}

// func ToAPI converts a pod EnvironmentVariable to an API EnvironmentVariable
func EnvVarToAPI(envVar EnvironmentVariable) api.EnvironmentVariable {
	return api.EnvironmentVariable{
		Key:    envVar.Key,
		Value:  envVar.Value,
		Secret: envVar.Secret,
	}
}

// ToAPIMany converts an array of pod EnvironmentVariables to an array of API EnvironmentVariables
func EnvVarToAPIMany(envVars []EnvironmentVariable) []api.EnvironmentVariable {
	apiEnvVars := make([]api.EnvironmentVariable, len(envVars))
	for i, envVar := range envVars {
		apiEnvVars[i] = EnvVarToAPI(envVar)
	}
	return apiEnvVars
}
