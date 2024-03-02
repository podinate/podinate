package router

import (
	"context"
	"io"
	"net/http"

	"strconv"

	"github.com/gorilla/mux"
	"github.com/johncave/podinate/api-backend/apierror"
	api "github.com/johncave/podinate/api-backend/go"
	"github.com/johncave/podinate/api-backend/iam"
	lh "github.com/johncave/podinate/api-backend/loghandler"
	pod "github.com/johncave/podinate/api-backend/pod"
	"github.com/johncave/podinate/api-backend/responder"
)

type PodAPIShim struct {
	service      api.PodApiServicer
	errorHandler api.ErrorHandler
}

func NewPodShimController(s api.PodApiServicer, opts ...api.PodApiOption) api.Router {
	controller := &PodAPIShim{
		service: s,
	}

	return controller
}

func (c *PodAPIShim) Routes() api.Routes {
	return api.Routes{
		{
			Name:        "ProjectProjectIdPodPodIdLogsGet",
			Method:      "GET",
			Pattern:     "/v0/project/{project_id}/pod/{pod_id}/logs",
			HandlerFunc: c.ProjectProjectIdPodPodIdLogsGet,
		},
	}
}

// ProjectProjectIdPodPodIdLogsGet - Get the logs for a pod
func (s *PodAPIShim) ProjectProjectIdPodPodIdLogsGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query := r.URL.Query()
	projectID := params["project_id"]
	podID := params["pod_id"]
	accountID := r.Header.Get("account")
	lines, err := strconv.Atoi(query.Get("lines"))
	if err != nil {
		api.EncodeJSONResponse(responder.Response(400, "Invalid lines parameter "+err.Error()), nil, w)
		return
	}

	follow := false

	if query.Get("follow") != "" {
		follow, err = strconv.ParseBool(query.Get("follow"))
		if err != nil {
			api.EncodeJSONResponse(responder.Response(400, "Invalid follow parameter "+err.Error()), nil, w)
			return
		}
	}

	// Get the account and project
	theAccount, theProject, apiErr := getAccountAndProject(accountID, projectID)
	if err != nil {
		apiErr.EncodeJSONResponse(w)
		return
	}

	p, apiErr := pod.GetByID(r.Context(), theProject, podID)
	if err != nil {
		apiErr.EncodeJSONResponse(w)
		return
	}

	if !iam.RequestorCan(r.Context(), theAccount, p, pod.ActionViewLogs) {
		apiErr := apierror.New(http.StatusForbidden, "You do not have permission to view the logs for this pod")
		apiErr.EncodeJSONResponse(w)
		return
	}

	in, apiErr := p.GetLogsBuffer(r.Context(), lines, follow)
	defer in.Close()
	if apiErr != nil {
		apiErr.EncodeJSONResponse(w)
		return
	}

	n, err := io.Copy(w, in)
	if err != nil {
		lh.Debug(r.Context(), "Error writing logs to response", "error", err, "bytes_written", n, "lines", lines)
		apierror.New(http.StatusInternalServerError, "Error writing logs to response "+err.Error()).EncodeJSONResponse(w)
		return
	}

	//w.Write([]byte("hello"))

}

// ProjectProjectIdPodPodIdLogsGet - Get the logs for a pod
func (s *PodAPIService) ProjectProjectIdPodPodIdLogsGet(ctx context.Context, projectID string, podId string, accountID string, lines int32, follow bool) (api.ImplResponse, error) {

	return api.Response(http.StatusInternalServerError, "This function should never have been called. Pls help."), nil
}
