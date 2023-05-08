package project

import (
	"net/http"

	"github.com/johncave/podinate/api-backend/account"
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
		return &apierror.ApiError{Code: http.StatusBadRequest, Message: "Project ID too long"}
	}
	if len(p.Name) > 64 {
		return &apierror.ApiError{Code: http.StatusBadRequest, Message: "Project name too long"}
	}
	return nil
}

// GetByID retrieves a project from the database by its id
func (p *Project) GetByID(theAccount account.Account, id string) *apierror.ApiError {
	err := config.DB.QueryRow("SELECT uuid, id, name FROM project WHERE account_uuid = $1 AND id = $2", theAccount.Uuid, id).Scan(&p.Uuid, &p.ID, &p.Name)
	if err != nil {
		return &apierror.ApiError{Code: http.StatusNotFound, Message: "Project not found"}
	}
	return nil
}

// GetProjects returns the projects of the account
func (p *Project) GetAccountProjects(a account.Account, page int32, limit int32) ([]project.Project, apierror.ApiError) {
	rows, err := config.DB.Query("SELECT uuid, id, name FROM project WHERE account_uuid = $1 OFFSET $2 LIMIT $3", a.Uuid, page, limit)
	if err != nil {
		return nil, apierror.New(http.StatusInternalServerError, "Could not retrieve projects")
	}
	defer rows.Close()
	// Read all the projects for the account
	projects := make([]Project, 0)
	for rows.Next() {
		var project Project
		err = rows.Scan(&project.Uuid, &project.ID, &project.Name)
		if err != nil {
			return nil, apierror.New(http.StatusInternalServerError, "Could not retrieve projects")
		}
		projects = append(projects, project)
	}
	return projects, apierror.New(http.StatusOK, "Projects retrieved successfully")

}

// Create creates a new project in the database
func (p *project.Project) CreateProject(new api.Project, inAccount Account) *apierror.ApiError {
	err := p.ValidateNew()
	if err != nil {
		return err
	}
	_, dberr := config.DB.Exec("INSERT INTO project(uuid, id, name, account_uuid) VALUES(gen_random_uuid(), $1, $2)", p.ID, p.Name)
	// Check if insert was successful
	if dberr != nil {
		return &apierror.ApiError{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	return nil
}

// ToAPI converts a project to an api.Project
func (p *Project) ToAPI() api.Project {
	return api.Project{Id: p.ID, Name: p.Name}
}
