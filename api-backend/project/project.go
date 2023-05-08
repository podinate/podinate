package project

import (
	"net/http"

	"github.com/johncave/podinate/api-backend/apierror"
	"github.com/johncave/podinate/api-backend/config"
	api "github.com/johncave/podinate/api-backend/go"
)

type Project struct {
	Uuid string `json:"uuid"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Validate project checks that a user's desired project properties are allowed
func (p *Project) ValidateNew() *apierror.ApiError {
	// check the project id and name are not too long
	if len(p.ID) > 30 {
		return &apierror.ApiError{Code: 400, Message: "Project ID too long"}
	}
	if len(p.Name) > 64 {
		return &apierror.ApiError{Code: 400, Message: "Project name too long"}
	}
	return nil
}

// Create creates a new project in the database
func (p *Project) Create(new api.Project) *apierror.ApiError {
	err := p.ValidateNew()
	if err != nil {
		return err
	}
	_, dberr := config.DB.Exec("INSERT INTO project(uuid, id, name) VALUES(gen_random_uuid(), $1, $2)", p.ID, p.Name)
	// Check if insert was successful
	if dberr != nil {
		return &apierror.ApiError{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	return nil
}

// GetByID retrieves a project from the database by its id
func (p *Project) GetByID(account string, id string) *apierror.ApiError {
	err := config.DB.QueryRow("SELECT uuid, id, name FROM project WHERE account = $1 AND id = $2", account, id).Scan(&p.Uuid, &p.ID, &p.Name)
	if err != nil {
		return &apierror.ApiError{Code: http.StatusNotFound, Message: "Project not found"}
	}
	return nil
}

// ToAPI converts a project to an api.Project
func (p *Project) ToAPI() api.Project {
	return api.Project{Id: p.ID, Name: p.Name}
}
