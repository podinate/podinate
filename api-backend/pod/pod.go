package pod

import (
	"net/http"

	"github.com/johncave/podinate/api-backend/apierror"
	api "github.com/johncave/podinate/api-backend/go"
)

type Pod struct {
	Uuid  string
	ID    string
	Name  string
	Image string
	Tag   string
	Count int
	Ram   int
	// TODO - add CPU requests / limits
}

func Create(requestedPod api.Pod) (Pod, *apierror.ApiError) {
	/// Ignore the Kubes logic for now
	err := createKubesNamespace("project-" + requestedPod.Id)
	if err != nil {
		return Pod{}, apierror.New(http.StatusInternalServerError, err.Error())
	}

	// err = createKubesDeployment("project-"+project.Id, project.Image, project.Tag)
	// if err != nil {
	// 	out := api.Model500Error{Code: http.StatusInternalServerError, Message: err.Error()}
	// 	return responder.Response(http.StatusInternalServerError, out), err
	// }

	return Pod{ID: requestedPod.Id}, nil
}
