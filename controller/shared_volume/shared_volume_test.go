package shared_volume

import (
	"context"
	"fmt"
	"testing"

	"github.com/johncave/podinate/controller/apierror"
	api "github.com/johncave/podinate/controller/go"
	lh "github.com/johncave/podinate/controller/loghandler"
	"github.com/johncave/podinate/controller/project"
	"github.com/johncave/podinate/controller/user"
)

var TestProject *project.Project
var TestContext context.Context

// TestCreateSharedVolume tests we can create a shared volume at all
func TestCreateSharedVolume(t *testing.T) {
	// Create a shared volume
	v := api.SharedVolume{
		Id:   "testsharedvolume",
		Size: 2,
	}
	created, err := Create(TestContext, TestProject, v)
	if err != nil {
		t.Fatalf("Error creating shared volume: %s", err)
	}

	// Check the values match
	if created.ID != v.Id {
		t.Errorf("Shared volume IDs do not match: %s != %s", created.ID, v.Id)
	}

	// Get the shared volume
	v2, err := GetByID(TestContext, TestProject, v.Id)
	if err != nil {
		t.Errorf("Error getting shared volume: %s", err)
	}

	if created.Uuid != v2.Uuid {
		t.Errorf("Shared volume UUIDs do not match: %s != %s", created.Uuid, v2.Uuid)
	}
	// Check values match
	if v2.ID != created.ID {
		t.Errorf("Shared volume IDs do not match: %s != %s", v2.ID, v.Id)
	}
	if v2.Size != created.Size {
		t.Errorf("Shared volume sizes do not match: %d != %d", v2.Size, v.Size)
	}

	// Delete the shared volume
	err = v2.Delete(TestContext)
	if err != nil {
		t.Errorf("Error deleting shared volume: %s", err)
	}

}

// TestEditSharedVolume tests we can edit a shared volume
func TestEditSharedVolume(t *testing.T) {
	// Create a shared volume
	v := api.SharedVolume{
		Id:   "testeditedsharedvolume",
		Name: "Test Edited Shared Volume",
		Size: 2,
	}
	created, err := Create(TestContext, TestProject, v)
	if err != nil {
		t.Fatalf("Error creating shared volume: %s", err)
	}

	created.Update(TestContext, api.SharedVolume{
		Id:   "testeditsharedvolume",
		Name: "Test Edit Shared Volume 2",
		Size: 3,
	})

	// Check the values match
	if created.Size != 3 {
		t.Errorf("Shared volume size does not match: %d != 3", created.Size)
	}
	// newName :=
	if *created.Name != "Test Edit Shared Volume 2" {
		t.Errorf("Shared volume name does not match: %s != Test Edit Shared Volume 2", *created.Name)
	}

	// Delete the shared volume
	err = created.Delete(TestContext)
	if err != nil {
		t.Errorf("Error deleting shared volume: %s", err)
	}

}

// TestCreateSharedVolume tests we can create a shared volume at all
func TestMain(m *testing.M) {
	// Create a test account
	var err *apierror.ApiError
	var u *user.User
	TestContext, TestProject, u, err = project.CreateTest()
	if err != nil {
		lh.Debug(TestContext, "Error creating test project", "error", err)
		lh.Panic(TestContext, "Error creating test project", "error", err)
	}

	fmt.Printf("Context: %v, Project: %v, User: %v\n", TestContext, TestProject, u)

	m.Run()

	err = project.DeleteTest(TestContext, TestProject, u)
	if err != nil {
		lh.Panic(TestContext, "Error deleting test project", "error", err)
	}

}
