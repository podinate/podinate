package router

import (
	"context"
	"encoding/json"
	"fmt"
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
		{
			Name:        "ProjectProjectIdPodPodIdExecPost",
			Method:      "POST",
			Pattern:     "/v0/project/{project_id}/pod/{pod_id}/exec",
			HandlerFunc: c.ProjectProjectIdPodPodIdExecPost,
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
	//defer in.Close()
	if apiErr != nil {
		apiErr.EncodeJSONResponse(w)
		return
	}
	defer in.Close()

	// buf := new(bytes.Buffer)
	// n, err := buf.ReadFrom(in)
	// n, err := io.Copy(w, in)
	// if err != nil {
	// 	lh.Debug(r.Context(), "Error writing logs to response", "error", err, "bytes_written", n, "lines", lines)
	// 	apierror.New(http.StatusInternalServerError, "Error writing logs to response "+err.Error()).EncodeJSONResponse(w)
	// 	return
	// }

	lh.Debug(r.Context(), "Writing logs to response", "lines", lines)
	totalRead := int64(0)
	for {
		n, err := io.CopyN(w, in, 100)
		totalRead += n
		// print bytes read followed by a carriage return
		fmt.Printf("Bytes read: %d\r", totalRead)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Finished copying") // print a newline
				break
			}
			fmt.Println("Error copying:", err)
			return
			// handle error
		}
	}
	// lh.Debug(r.Context(), "Wrote logs to response, closing buffer", "bytes_written", n, "lines", lines)
	// in.Close()
	//w.Write([]byte("hello"))

}

// ProjectProjectIdPodPodIdExecPost - Executes a command in a pod for a project.
// ProjectProjectIdPodPodIdExecPost - Execute a command in a pod
func (s *PodAPIShim) ProjectProjectIdPodPodIdExecPost(w http.ResponseWriter, r *http.Request) {
	// Parameter grabbing logic from the original function
	params := mux.Vars(r)
	projectId := params["project_id"]
	podId := params["pod_id"]
	account := r.Header.Get("account")
	projectProjectIdPodPodIdExecPostRequestParam := api.ProjectProjectIdPodPodIdExecPostRequest{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&projectProjectIdPodPodIdExecPostRequestParam); err != nil {
		apierror.NewWithError(http.StatusBadRequest, "error decoding request body", err).EncodeJSONResponse(w)
		return

	}
	if err := api.AssertProjectProjectIdPodPodIdExecPostRequestRequired(projectProjectIdPodPodIdExecPostRequestParam); err != nil {
		apierror.NewWithError(http.StatusBadRequest, "error validating request body", err).EncodeJSONResponse(w)
		return
	}
	ctx := r.Context()

	// Our Logic - like in handlers
	theAccount, theProject, err := getAccountAndProject(account, projectId)
	if err != nil {
		err.EncodeJSONResponse(w)
		return
	}

	thePod, err := pod.GetByID(ctx, theProject, podId)
	if err != nil {
		err.EncodeJSONResponse(w)
		return
	}

	if !iam.RequestorCan(ctx, theAccount, thePod, pod.ActionExec) {
		apierror.New(http.StatusForbidden, "You do not have permission to execute commands in this pod").EncodeJSONResponse(w)
		return
	}

	result, err := thePod.Exec(ctx, projectProjectIdPodPodIdExecPostRequestParam.Command)
	if err != nil {
		lh.Error(ctx, "Error executing command", "error", err, "command", projectProjectIdPodPodIdExecPostRequestParam.Command, "result", result)
		w.Write([]byte(err.Error()))
		return
	}

	lh.Debug(ctx, "Executed command without error", "result", result)
	w.Write([]byte(result))
}

/*
These functions are just so that we satisfy the interface from the library.
We overrode them in the previous functions so that we had direct control of the reponse.
*/
func (s *PodAPIService) ProjectProjectIdPodPodIdExecPost(ctx context.Context, projectId string, podId string, account string, projectProjectIdPodPodIdExecPostRequest api.ProjectProjectIdPodPodIdExecPostRequest) (api.ImplResponse, error) {

	return api.Response(http.StatusInternalServerError, "This function should never have been called. Pls help."), nil
}
func (s *PodAPIService) ProjectProjectIdPodPodIdLogsGet(ctx context.Context, projectID string, podId string, accountID string, lines int32, follow bool) (api.ImplResponse, error) {
	return api.Response(http.StatusInternalServerError, "This function should never have been called. Pls help."), nil
}
