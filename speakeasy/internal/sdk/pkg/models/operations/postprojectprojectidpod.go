// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/podinate/terraform-provider-podinate/internal/sdk/pkg/models/shared"
	"net/http"
)

type PostProjectProjectIDPodRequest struct {
	ProjectID string `pathParam:"style=simple,explode=false,name=project_id"`
	// The account to use for the request
	Account string `header:"style=simple,explode=false,name=account"`
	// A JSON object containing the information needed to create a new pod
	Pod shared.Pod `request:"mediaType=application/json"`
}

func (o *PostProjectProjectIDPodRequest) GetProjectID() string {
	if o == nil {
		return ""
	}
	return o.ProjectID
}

func (o *PostProjectProjectIDPodRequest) GetAccount() string {
	if o == nil {
		return ""
	}
	return o.Account
}

func (o *PostProjectProjectIDPodRequest) GetPod() shared.Pod {
	if o == nil {
		return shared.Pod{}
	}
	return o.Pod
}

type PostProjectProjectIDPodResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// The pod was created successfully
	Pod *shared.Pod
	// Request issued incorrectly, for example missing parameters or wrong endpoint
	Error *shared.Error
}

func (o *PostProjectProjectIDPodResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *PostProjectProjectIDPodResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *PostProjectProjectIDPodResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *PostProjectProjectIDPodResponse) GetPod() *shared.Pod {
	if o == nil {
		return nil
	}
	return o.Pod
}

func (o *PostProjectProjectIDPodResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
