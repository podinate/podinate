package sdk

import (
	"context"
	"errors"

	"github.com/johncave/podinate/lib/api_client"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Project struct {
	ID         string
	Name       string
	ResourceID string
}

// GetProjectByID returns a project by its ID
func GetProjectByID(id string) (*Project, error) {
	// Get the project from the API

	if id == "" {
		return nil, errors.New("No project selected")
	}

	resp, r, err := C.ProjectApi.ProjectIdGet(context.Background(), id).Account(viper.GetString("account")).Execute()
	if err := handleAPIError(r, err); err != nil {
		return nil, err
	}
	return &Project{
		ID:         *resp.Id,
		Name:       *resp.Name,
		ResourceID: *resp.ResourceId,
	}, nil
}

func CreateProject(id string, name string) (*Project, error) {
	// Create the project in the API
	resp, _, err := C.ProjectApi.ProjectPost(context.Background()).Account(viper.GetString("account")).Project(api_client.Project{
		Id:   &id,
		Name: &name,
	}).Execute()
	if err != nil {
		return nil, err
	}
	return &Project{
		ID:         *resp.Id,
		Name:       *resp.Name,
		ResourceID: *resp.ResourceId,
	}, nil
}

// func ProjectGetByID(id string) (*Project, error) {
// 	// Get the project from the API
// 	resp, r, err := C.ProjectApi.ProjectIdGet(context.Background(), id).Account(viper.GetString("account")).Execute()
// 	fmt.Println("ProjectGetByID", id)
// 	if err := handleAPIError(r, err); err != nil {
// 		return nil, err
// 	}
// 	return &Project{
// 		ID:         *resp.Id,
// 		Name:       *resp.Name,
// 		ResourceID: *resp.ResourceId,
// 	}, nil
// }

func GetAllProjects(accountID string) ([]Project, error) {
	// Get the project from the API
	resp, r, err := C.ProjectApi.ProjectGet(context.Background()).Account(accountID).Execute()
	if err := handleAPIError(r, err); err != nil {
		return nil, err
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

// Update updates a project
func (p *Project) Update(in *Project) (*Project, error) {

	req := in.ToAPI()

	resp, r, err := C.ProjectApi.ProjectIdPut(context.Background(), p.ID).Account(viper.GetString("account")).Project(req).Execute()
	if err := handleAPIError(r, err); err != nil {
		return nil, err
	}

	return &Project{
		ID:         *resp.Id,
		Name:       *resp.Name,
		ResourceID: *resp.ResourceId,
	}, nil

}

// Describe returns a string representation of the project
func (p *Project) Describe() (string, error) {
	out, err := yaml.Marshal(p)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (p *Project) Delete() error {
	r, err := C.ProjectApi.ProjectIdDelete(context.Background(), p.ID).Account(viper.GetString("account")).Execute()
	if err := handleAPIError(r, err); err != nil {
		return err
	}
	return err
}

func (p *Project) ToAPI() api_client.Project {
	return api_client.Project{
		Id:   &p.ID,
		Name: &p.Name,
	}
}

// FromAPI returns the API client representation of the project
func (p *Project) FromAPI(in *api_client.Project) {
	p.ID = *in.Id
	p.Name = *in.Name
	p.ResourceID = *in.ResourceId
}
