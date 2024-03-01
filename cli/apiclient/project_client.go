package apiclient

import (
	"context"
	"errors"

	"encoding/json"

	"github.com/johncave/podinate/lib/api_client"
	"github.com/spf13/viper"
)

type Project struct {
	ID         string
	Name       string
	ResourceID string
}

func ProjectGetByID(id string) (*Project, error) {
	// Get the project from the API
	resp, _, err := C.ProjectApi.ProjectIdGet(context.Background(), id).Account(viper.GetString("account")).Execute()
	if err != nil {
		return nil, err
	}
	return &Project{
		ID:         *resp.Id,
		Name:       *resp.Name,
		ResourceID: *resp.ResourceId,
	}, nil
}

func GetAllProjects(accountID string) ([]Project, error) {
	// Get the project from the API
	resp, r, err := C.ProjectApi.ProjectGet(context.Background()).Account(accountID).Execute()

	if r == nil {
		return nil, errors.New("Could not connect to Podinate API: " + err.Error())
	}
	if err != nil {
		var apierr api_client.Error
		err = json.NewDecoder(r.Body).Decode(&apierr)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(*apierr.Message)
	}

	var projects []Project
	for _, p := range resp.Items {
		projects = append(projects, Project{
			ID:         *p.Project.Id,
			Name:       *p.Project.Name,
			ResourceID: *p.Project.ResourceId,
		})
	}
	return projects, nil
}

func (p *Project) Delete() error {
	_, err := C.ProjectApi.ProjectIdDelete(context.Background(), p.ID).Account(viper.GetString("account")).Execute()
	return err
}
