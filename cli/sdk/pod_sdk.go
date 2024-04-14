package sdk

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/johncave/podinate/lib/api_client"
	"github.com/spf13/viper"
)

type Pod struct {
	Project *Project
	// The short name (slug/url) of the pod
	ID string `json:"id"`
	// The name of the pod
	Name          string                      `json:"name"`
	Image         string                      `json:"image"`
	Tag           *string                     `json:"tag"`
	Command       []string                    `json:"command"`
	Arguments     []string                    `json:"arguments"`
	Status        string                      `json:"status"`
	CreatedAt     string                      `json:"created_at"`
	ResourceID    string                      `json:"resource_id"`
	Volumes       VolumeSlice                 `json:"volumes"`
	Services      ServiceSlice                `json:"services"`
	Environment   EnvironmentSlice            `json:"environment"`
	SharedVolumes SharedVolumeAttachmentSlice `json:"shared_volumes"`
}

// GetPodByID returns a pod by ID from the given project
func (p *Project) GetPodByID(id string) (*Pod, error) {
	resp, r, err := C.PodApi.ProjectProjectIdPodPodIdGet(context.Background(), p.ID, id).Account(viper.GetString("account")).Execute()
	if err := handleAPIError(r, err); err != nil {
		return nil, err
	}
	return getPodFromApi(p, resp), nil

}

// GetPods returns all pods from the given project
func (p *Project) GetPods() ([]*Pod, *SDKError) {
	resp, r, err := C.PodApi.ProjectProjectIdPodGet(context.Background(), p.ID).Account(viper.GetString("account")).Execute()
	if err := handleAPIError(r, err); err != nil {
		return nil, err
	}

	var pods []*Pod

	for _, i := range resp.Items {
		po := i.Pod

		pods = append(pods, getPodFromApi(p, po))
	}

	return pods, nil
}

// Create creates a new pod in the given project
func (p *Project) CreatePod(in Pod) (*Pod, error) {

	req := in.ToAPI()

	resp, r, err := C.PodApi.ProjectProjectIdPodPost(context.Background(), p.ID).Account(viper.GetString("account")).Pod(req).Execute()
	if err := handleAPIError(r, err); err != nil {
		return nil, err
	}

	return getPodFromApi(p, resp), nil

}

// Update updates a pod in the given project
func (p *Pod) Update(in *Pod) error {

	req := in.ToAPI()

	resp, r, err := C.PodApi.ProjectProjectIdPodPodIdPut(context.Background(), p.Project.ID, p.ID).Account(viper.GetString("account")).Pod(req).Execute()
	if err := handleAPIError(r, err); err != nil {
		return err
	}

	p = getPodFromApi(p.Project, resp)
	return nil

}

// getLogs returns the logs for a pod
func (p *Pod) GetLogs(lines int, follow bool) (string, error) {
	resp, _, err := C.PodApi.ProjectProjectIdPodPodIdLogsGet(context.Background(), p.Project.ID, p.ID).Account(viper.GetString("account")).Lines(int32(lines)).Execute()
	return resp, err
}

// getLogsBuffer returns the logs for a pod
func (p *Pod) GetLogsBuffer(lines int, follow bool) (io.ReadCloser, error) {
	//_, r, err := C.PodApi.ProjectProjectIdPodPodIdLogsGet(context.Background(), p.Project.ID, p.ID).Account(viper.GetString("account")).Lines(int32(lines)).Follow(follow).Execute()
	hc := C.GetConfig().HTTPClient
	req, err := http.NewRequest("GET", C.GetConfig().Scheme+"://"+C.GetConfig().Host+"/v0/project/"+p.Project.ID+"/pod/"+p.ID+"/logs?lines="+strconv.Itoa(lines)+"&follow=true", nil)
	req.Header.Set("Authorization", viper.GetString("api_key"))
	req.Header.Set("Account", viper.GetString("account"))
	r, err := hc.Do(req)
	//fmt.Println("r", r, "err", err)
	return r.Body, err
}

// Exec executes a command in the pod
func (p *Pod) Exec(command []string) (string, error) {
	req := *api_client.NewProjectProjectIdPodPodIdExecPostRequest(command)
	resp, _, err := C.PodApi.ProjectProjectIdPodPodIdExecPost(context.Background(), p.Project.ID, p.ID).
		Account(viper.GetString("account")).
		ProjectProjectIdPodPodIdExecPostRequest(req).Execute()
	//fmt.Println("resp", resp, "err", err)
	return resp, err
}

// Delete deletes the pod
func (p *Pod) Delete() error {
	r, err := C.PodApi.ProjectProjectIdPodPodIdDelete(context.Background(), p.Project.ID, p.ID).Account(viper.GetString("account")).Execute()
	return handleAPIError(r, err)
}

func getPodFromApi(p *Project, in *api_client.Pod) *Pod {
	//fmt.Println("in.Id", in.Id, "in", in, "created", in.CreatedAt)
	//fmt.Println("%+V\n", in)
	out := &Pod{
		ID:            in.Id,
		Name:          in.Name,
		Image:         in.Image,
		Tag:           in.Tag,
		Command:       in.Command,
		Arguments:     in.Arguments,
		Status:        *in.Status,
		ResourceID:    *in.ResourceId,
		Project:       p,
		Volumes:       volumesFromAPI(in.Volumes),
		Services:      servicesFromAPI(in.Services),
		Environment:   environmentVariablesFromAPI(in.Environment),
		SharedVolumes: sharedVolumeAttachmentsFromAPI(in.SharedVolumes),
	}

	return out
}

// ToAPI returns a Pod from the API client representation
func (p *Pod) ToAPI() api_client.Pod {
	return api_client.Pod{
		Id:            p.ID,
		Name:          p.Name,
		Image:         p.Image,
		Tag:           p.Tag,
		Command:       p.Command,
		Arguments:     p.Arguments,
		Environment:   environmentVariablesToAPI(p.Environment),
		Volumes:       volumesToAPI(p.Volumes),
		Services:      servicesToAPI(p.Services),
		SharedVolumes: sharedVolumeAttachmentsToAPI(p.SharedVolumes),
	}
}
