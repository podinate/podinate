package apiclient

import (
	"context"

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

func (p *Project) Delete() error {
	_, err := C.ProjectApi.ProjectIdDelete(context.Background(), p.ID).Account(viper.GetString("account")).Execute()
	return err
}
