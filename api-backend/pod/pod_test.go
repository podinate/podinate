package pod

import (
	"os"
	"testing"

	"github.com/johncave/podinate/api-backend/apierror"
	"github.com/johncave/podinate/api-backend/config"
	api "github.com/johncave/podinate/api-backend/go"
	lh "github.com/johncave/podinate/api-backend/loghandler"
	"github.com/johncave/podinate/api-backend/project"
)

var testProject project.Project

func TestCreateThenGet(t *testing.T) {
	// Create a pod
	ctx := lh.TestContext()
	p := api.Pod{
		Id:    "testcreatethenget",
		Name:  "TestCreateThenGet",
		Image: "nginx",
		Tag:   "latest",
	}
	_, err := Create(ctx, testProject, p)
	if err != nil {
		t.Errorf("Error creating pod: %s", err)
	}

	// Get the pod
	p2, err := GetByID(ctx, testProject, p.Id)
	if err != nil {
		t.Errorf("Error getting pod: %s", err)
	}
	if p2.Name != p.Name {
		t.Errorf("Pod names do not match: %s != %s", p2.Name, p.Name)
	}
	if p2.Status != "Running" {
		t.Errorf("Pod status is not Running: %s", p2.Status)
	}

	// Clean up
	err = p2.Delete(ctx)
	if err != nil {
		t.Errorf("Error deleting pod: %s", err)
	}

	// p3, err := GetByID(ctx, testProject, p.Id)
	// if err == nil {
	// 	t.Errorf("Pod should have been deleted: %+v", p3)
	// }
}

func TestMain(m *testing.M) {
	//config.Init()

	// Create a test account
	var err *apierror.ApiError
	testProject, err = project.CreateTest()
	if err != nil {
		lh.Log.Errorw("Error creating test project", "error", err)
		os.Exit(2)
	}

	code := m.Run()
	err = testProject.Account.Delete()
	if err != nil {
		lh.Log.Errorw("Error deleting test account", "error", err)
		os.Exit(3)
	}
	err = testProject.Delete()
	if err != nil {
		lh.Log.Errorw("Error deleting test project", "error", err)
		os.Exit(4)
	}
	config.Cleanup()
	os.Exit(code)
}
