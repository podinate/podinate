package project

import (
	"log"
	"net/http"

	"github.com/johncave/podinate/api-backend/account"
	"github.com/johncave/podinate/api-backend/apierror"
	"github.com/johncave/podinate/api-backend/config"
	api "github.com/johncave/podinate/api-backend/go"
)

type Project struct {
	Uuid    string          `json:"uuid"`
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Account account.Account `json:"-"`
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
	log.Println("Validated project")
	return nil
}

// // GetByID retrieves a project from the database by its id
// func GetByID(theAccount account.Account, id string) (Project, *apierror.ApiError) {
// 	var p Project
// 	err := config.DB.QueryRow("SELECT uuid, id, name FROM project WHERE account_uuid = $1 AND id = $2", theAccount.Uuid, id).Scan(&p.Uuid, &p.ID, &p.Name)
// 	if err != nil {
// 		return Project{}, &apierror.ApiError{Code: http.StatusNotFound, Message: "Project not found"}
// 	}
// 	return p, nil
// }

// GetProjects returns the projects of the account
func GetByAccount(a account.Account, page int32, limit int32) ([]Project, *apierror.ApiError) {
	//log.Printf("Getting projects for account %s on page %s limit %s", a.Uuid, page, limit)

	if limit < 1 {
		limit = 25
	}
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
		project.Account = a
		projects = append(projects, project)
	}
	return projects, nil

}

// Patch updates an account in the database
func (p *Project) Patch(requested api.Project) *apierror.ApiError {
	// Check which fields are actually being updated
	if requested.Id != "" {
		p.ID = requested.Id
	}
	if requested.Name != "" {
		p.Name = requested.Name
	}
	// Update the database
	_, err := config.DB.Exec("UPDATE project SET name = $1, id = $2 WHERE uuid = $3", p.Name, p.ID, p.Uuid)
	if err != nil {
		log.Println(err)
		return apierror.New(http.StatusInternalServerError, "Could not update project")
	}
	return nil
}

// Create creates a new project in the database
func Create(new api.Project, inAccount account.Account) (Project, *apierror.ApiError) {
	log.Println("Creating project")
	out := Project{ID: new.Id, Name: new.Name}
	err := out.ValidateNew()
	if err != nil {
		return Project{}, err
	}
	//res, dberr := config.DB.Exec("INSERT INTO project(uuid, id, name, account_uuid) VALUES(gen_random_uuid(), $1, $2, $3) RETURNING uuid", new.Id, new.Name, inAccount.Uuid)
	dberr := config.DB.QueryRow("INSERT INTO project(uuid, id, name, account_uuid) VALUES(gen_random_uuid(), $1, $2, $3) RETURNING uuid", new.Id, new.Name, inAccount.Uuid).Scan(&out.Uuid)
	// Check if insert was successful
	if dberr != nil {
		log.Println("DB error", dberr)
		return Project{}, &apierror.ApiError{Code: http.StatusInternalServerError, Message: dberr.Error()}
	}

	log.Printf("Created project %s / %s", out.ID, out.Name)
	return out, nil
}

// ToAPI converts a project to an api.Project
func (p *Project) ToAPI() api.Project {
	return api.Project{Id: p.ID, Name: p.Name}
}

func GetByID(a account.Account, id string) (Project, *apierror.ApiError) {
	row := config.DB.QueryRow("SELECT uuid, id, name FROM project WHERE account_uuid = $1 AND id = $2", a.Uuid, id)

	p := Project{}
	err := row.Scan(&p.Uuid, &p.ID, &p.Name)

	if err != nil {
		return Project{}, apierror.New(http.StatusNotFound, "Could not find project")
	}
	p.Account = a
	return p, nil

}

// Delete deletes a project from the database
func (p *Project) Delete() *apierror.ApiError {
	_, err := config.DB.Exec("DELETE FROM project WHERE uuid = $1", p.Uuid)
	if err != nil {
		return apierror.New(http.StatusInternalServerError, "Could not delete project")
	}
	return nil
}
