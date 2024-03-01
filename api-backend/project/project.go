package project

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/johncave/podinate/api-backend/account"
	"github.com/johncave/podinate/api-backend/apierror"
	"github.com/johncave/podinate/api-backend/config"
	api "github.com/johncave/podinate/api-backend/go"
	"github.com/lib/pq"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Project struct {
	Uuid    string          `json:"uuid"`
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Account account.Account `json:"-"`
}

const (
	ActionCreate = "project:create"
	ActionView   = "project:view"
	ActionUpdate = "project:update"
	ActionDelete = "project:delete"
)

// Validate project checks that a user's desired project properties are allowed
func (p *Project) ValidateNew() *apierror.ApiError {
	// check the project id and name are not too long
	if len(p.ID) > 30 {
		return &apierror.ApiError{Code: http.StatusBadRequest, Message: "Project ID too long"}
	}
	if len(p.Name) > 64 {
		return &apierror.ApiError{Code: http.StatusBadRequest, Message: "Project name too long"}
	}

	m, err := regexp.MatchString(`^([a-z]*[0-9]*-*)*$`, p.ID)
	if err != nil {
		return &apierror.ApiError{Code: http.StatusBadRequest, Message: "Error checking project ID " + err.Error()}
	}
	if !m { // ID must be lowercase letters and numbers only
		return apierror.New(http.StatusBadRequest, "Project ID "+p.ID+" invalid, must be lowercase letters, numbers and dashes only")
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
	rows, err := config.DB.Query("SELECT uuid, id, name FROM project WHERE account_uuid = $1 OFFSET $2 LIMIT $3", a.GetUUID(), page, limit)
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

// CreateTest spins up a new project into which some tests can be putted
func CreateTest() (Project, *apierror.ApiError) {
	newAcc, err := account.CreateTest()
	if err != nil {
		return Project{}, err
	}

	//lh.Log.Debug("Created test account", "account", newAcc)

	rand.Seed(time.Now().UnixNano())
	id := generateRandomString(8)
	name := generateRandomString(10)

	newProj := api.Project{Id: id, Name: name}
	return Create(newProj, *newAcc)
}

func generateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// Create creates a new project in the database
func Create(new api.Project, inAccount account.Account) (Project, *apierror.ApiError) {
	//lh.Log.Debugw("Creating project", "project", new, "account", inAccount)
	out := Project{ID: new.Id, Name: new.Name, Account: inAccount}
	err := out.ValidateNew()
	if err != nil {
		return Project{}, err
	}
	//res, dberr := config.DB.Exec("INSERT INTO project(uuid, id, name, account_uuid) VALUES(gen_random_uuid(), $1, $2, $3) RETURNING uuid", new.Id, new.Name, inAccount.Uuid)
	dberr := config.DB.QueryRow("INSERT INTO project(uuid, id, name, account_uuid) VALUES(gen_random_uuid(), $1, $2, $3) RETURNING uuid", new.Id, new.Name, inAccount.GetUUID()).Scan(&out.Uuid)
	// Check if insert was successful
	if dberr != nil && dberr.(*pq.Error).Code.Name() == "unique_violation" {
		return Project{}, &apierror.ApiError{Code: http.StatusConflict, Message: "Project ID already exists"}
	}
	if dberr != nil {
		log.Println("DB error", dberr)
		return Project{}, &apierror.ApiError{Code: http.StatusInternalServerError, Message: dberr.Error()}
	}

	log.Printf("Created project %s / %s", out.ID, out.Name)
	return out, nil
}

// ToAPI converts a project to an api.Project
func (p *Project) ToAPI() api.Project {
	return api.Project{Id: p.ID, Name: p.Name, ResourceId: p.GetResourceID()}
}

func GetByID(a account.Account, id string) (Project, *apierror.ApiError) {
	row := config.DB.QueryRow("SELECT uuid, id, name FROM project WHERE account_uuid = $1 AND id = $2", a.GetUUID(), id)

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
	// Todo: Delete every pod in the project
	_, err := config.DB.Exec("DELETE FROM project WHERE uuid = $1", p.Uuid)
	if err != nil {
		return apierror.New(http.StatusInternalServerError, "Could not delete project")
	}

	// Delete the kubernetes namespace
	apiErr := p.deleteKubeNamespace()
	if err != nil {
		return apiErr
	}
	return nil
}

// Delete Kubernetes namespace if it exists
// func (p *Project) DeleteKubeNamespace() *apierror.ApiError {

func (p Project) GetResourceID() string {
	return p.Account.GetResourceID() + "/project:" + p.ID
}

// deleteKubeNamespace deletes a namespace from the kubernetes cluster
func (p *Project) deleteKubeNamespace() *apierror.ApiError {
	// Delete the namespace
	err := config.Client.CoreV1().Namespaces().Delete(context.Background(), p.GetNamespaceName(), metav1.DeleteOptions{})
	if err != nil {
		log.Printf("error deleting namespace: %v\n", err)
		return apierror.New(500, "error deleting namespace: "+err.Error())
	}
	return nil
}

func (p *Project) GetNamespaceName() string {
	return p.Account.ID + "-project-" + p.ID
}
