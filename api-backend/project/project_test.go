package project

import (
	"math/rand"
	"testing"

	"github.com/johncave/podinate/api-backend/account"
	"github.com/johncave/podinate/api-backend/config"
	api "github.com/johncave/podinate/api-backend/go"
)

func TestBadProjectID(t *testing.T) {
	proj := Project{ID: RandStringBytes(64), Name: RandStringBytes(20)}
	err := proj.ValidateNew()
	if err == nil {
		t.Error("Expected error for too long project ID")
	}
}

func TestBadProjectName(t *testing.T) {
	proj := Project{ID: RandStringBytes(10), Name: RandStringBytes(128)}
	err := proj.ValidateNew()
	if err == nil {
		t.Error("Expected error for too long project name")
	}
}

func TestProjectCreation(t *testing.T) {
	InitDB()
	newAcc, err := account.Create(api.Account{Id: RandStringBytes(10), Name: RandStringBytes(20)})
	if err != nil {
		t.Error("Error encountered creating account")
	}
	apiproj := api.Project{Id: RandStringBytes(10), Name: RandStringBytes(20)}
	theProject, err := Create(apiproj, newAcc)
	if err != nil {
		t.Error("Expected no error creating project"+err.Error(), err.Code)
	}
	if theProject.ID != apiproj.Id {
		t.Error("Expected project ID to match", theProject.ID, " ", apiproj.Id)
	}
	if theProject.Name != apiproj.Name {
		t.Error("Expected project name to match", theProject.Name, " ", apiproj.Name)
	}
	if theProject.ToAPI() != apiproj {
		t.Error("Expected project to match", theProject.ToAPI(), " ", apiproj)
	}

	t.Log("Created project", theProject.ID, " ", theProject.Name, " ", theProject.Uuid, "now testing getting project by ID")
	gotProj, err := GetByID(newAcc, theProject.ID)
	if err != nil {
		t.Error("Expected no error getting project by ID", err.Error())
	}
	if gotProj.ID != theProject.ID {
		t.Error("Expected project ID to match", gotProj.ID, " ", theProject.ID)
	}
	if gotProj.Name != theProject.Name {
		t.Error("Expected project name to match", gotProj.Name, " ", theProject.Name)
	}
	if gotProj.Uuid != theProject.Uuid {
		t.Error("Expected project UUID to match", gotProj.Uuid, " ", theProject.Uuid)
	}

	t.Log("Testing getting projects by account")
	projects, err := GetByAccount(newAcc, 0, 2)
	if err != nil {
		t.Error("Expected no error getting projects by account", err.Error())
	}
	if len(projects) != 1 {
		t.Error("Expected 1 project to be returned")
	}
	if projects[0].ID != theProject.ID {
		t.Error("Expected project ID to match", projects[0].ID, " ", theProject.ID)
	}

	t.Log("Testing deleting project and cleaning up")
	err = gotProj.Delete()
	if err != nil {
		t.Error("Expected no error deleting project", err.Error())
	}
	err = newAcc.Delete()
	if err != nil {
		t.Error("Expected no error deleting account", err.Error())
	}
}

func TestProjectUpdate(t *testing.T) {
	InitDB()
	newAcc, err := account.Create(api.Account{Id: RandStringBytes(10), Name: RandStringBytes(20)})
	if err != nil {
		t.Error("Error encountered creating account")
	}
	apiproj := api.Project{Id: RandStringBytes(10), Name: RandStringBytes(20)}
	theProject, err := Create(apiproj, newAcc)
	if err != nil {
		t.Error("Expected no error creating project"+err.Error(), err.Code)
	}

	update := api.Project{Id: "", Name: RandStringBytes(20)}
	err = theProject.Patch(update)
	if err != nil {
		t.Error("Expected no error updating project", err.Error())
	}
	if theProject.Name != update.Name {
		t.Error("Expected project name to be updated", theProject.Name, " ", update.Name)
	}

	t.Log("Getting the project again to check the update worked")
	gotProj, err := GetByID(newAcc, theProject.ID)
	if err != nil {
		t.Error("Expected no error getting project by ID", err.Error())
	}
	if gotProj.ID != theProject.ID {
		t.Error("Expected project ID to be the same", gotProj.ID, " ", theProject.ID)
	}
	if gotProj.Name != update.Name {
		t.Error("Expected project name to be updared", gotProj.Name, " ", update.Name)
	} else {
		t.Log("Project name was updated correctly")
	}
	if gotProj.Uuid != theProject.Uuid {
		t.Error("Expected project UUID to be the same", gotProj.Uuid, " ", theProject.Uuid)
	}

	t.Log("Testing deleting project and cleaning up")
	err = gotProj.Delete()
	if err != nil {
		t.Error("Expected no error deleting project", err.Error())
	}
	err = newAcc.Delete()
	if err != nil {
		t.Error("Expected no error deleting account", err.Error())
	}

}

// Utility functions for test

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func InitDB() {
	config.Init()
}
