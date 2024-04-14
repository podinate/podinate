package project

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/johncave/podinate/controller/account"
	"github.com/johncave/podinate/controller/apierror"
	"github.com/johncave/podinate/controller/config"
	api "github.com/johncave/podinate/controller/go"
	"github.com/johncave/podinate/controller/iam"
	lh "github.com/johncave/podinate/controller/loghandler"
	"github.com/johncave/podinate/controller/user"
	"github.com/lib/pq"
	"golang.org/x/exp/maps"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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
func CreateTest() (context.Context, *Project, *user.User, *apierror.ApiError) {
	// Create test account
	newAcc, u, err := account.CreateTest()
	if err != nil {
		return nil, nil, nil, err
	}

	// First off, create a test context
	ctx := iam.TestContext(u)

	// Add initial policies to the account
	superAdminPolicyDocument := `
version: 2023.1
statements:
  - effect: allow
    actions: ["**"]
    resources: ["**"]`
	superAdminPolicy, err := iam.CreatePolicyForAccount(newAcc, "super-administrator", superAdminPolicyDocument, "Default policy created during initial account creation")
	err = superAdminPolicy.AttachToRequestor(u, u)
	if err != nil {
		// We can pass this error directly to the API response
		lh.Log.Fatalw("Error attaching super-administrator policy to initial default account", "error", err)
		return nil, nil, nil, err
	}

	//lh.Log.Debug("Created test account", "account", newAcc)

	rand.Seed(time.Now().UnixNano())
	id := generateRandomString(8)
	name := generateRandomString(10)

	newProj := api.Project{Id: id, Name: name}
	out, err := Create(ctx, newProj, *newAcc)

	return ctx, &out, u, err
}

// DeleteTest deletes a test project
func DeleteTest(ctx context.Context, p *Project, u *user.User) *apierror.ApiError {
	err := p.Account.Delete()
	if err != nil {
		return err
	}

	err = p.Delete()
	if err != nil {
		return err
	}
	err = u.Delete(ctx)
	if err != nil {
		return err
	}
	return nil
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
func Create(ctx context.Context, new api.Project, inAccount account.Account) (Project, *apierror.ApiError) {
	//lh.Log.Debugw("Creating project", "project", new, "account", inAccount)
	out := Project{ID: new.Id, Name: new.Name, Account: inAccount}
	err := out.ValidateNew()
	if err != nil {
		return Project{}, err
	}

	// Create the Kubernetes Namespace
	_, err = out.EnsureNamespace(ctx)

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

func (p *Project) EnsureNamespace(ctx context.Context) (*corev1.Namespace, *apierror.ApiError) {
	fmt.Println("Create Kubernetes namespace")

	// clientset, err := getKubesClient()
	// if err != nil {
	// 	log.Printf("error getting kubernetes client: %v\n", err)
	// 	return nil, err
	// }

	nsSpec := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: p.GetNamespaceName(),
		},
	}
	ns, err := config.Client.CoreV1().
		Namespaces().
		Create(context.Background(), nsSpec, metav1.CreateOptions{})
	if errors.IsAlreadyExists(err) {
		// Get the ns instead
		ns, err := config.Client.CoreV1().Namespaces().Get(ctx, p.GetNamespaceName(), metav1.GetOptions{})
		if err != nil {
			lh.Error(ctx, "Error getting existing kubernetes namespace", err, "project", p, "namespace", nsSpec)
			return ns, apierror.NewWithError(500, "error getting existing kubernetes namespace", err)
		}
		return ns, nil
	}
	if err != nil {
		lh.Error(ctx, "Error creating kubernetes namespace", err, "project", p, "namespace", nsSpec)
		return nil, apierror.NewWithError(500, "error creating kubernetes namespace", err)
	}
	return ns, nil
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
		if err != nil && err.Error() == "namespaces \""+p.GetNamespaceName()+"\" not found" {
			return nil
		}
		log.Printf("error deleting namespace: %v\n", err)
		return apierror.New(500, "error deleting namespace: "+err.Error())
	}
	return nil
}

func (p *Project) GetNamespaceName() string {
	if p.Account.ID == "default" {
		return p.ID
	}
	return p.Account.ID + "-project-" + p.ID
}

// GetLabels returns the labels for the project
func (p *Project) GetLabels() map[string]string {
	out := map[string]string{
		"podinate.com/project": p.ID,
	}

	maps.Copy(out, p.Account.GetLabels())
	return out
}

// GetAnnotations returns the annotations for the project
func (p *Project) GetAnnotations() map[string]string {
	out := map[string]string{}

	maps.Copy(out, p.Account.GetAnnotations())
	return out
}
