package sdk

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/johncave/podinate/lib/api_client"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Pod struct {
	Project *Project
	// The short name (slug/url) of the pod
	ID string `json:"id"`
	// The name of the pod
	Name          string                      `yaml:"name"`
	Image         string                      `yaml:"image"`
	Tag           *string                     `yaml:"tag"`
	Command       []string                    `yaml:"command"`
	Arguments     []string                    `yaml:"arguments"`
	Status        string                      `yaml:"status"`
	CreatedAt     string                      `yaml:"created_at"`
	ResourceID    string                      `yaml:"resource_id"`
	Volumes       VolumeSlice                 `yaml:"volumes"`
	Services      ServiceSlice                `yaml:"services"`
	Environment   EnvironmentSlice            `yaml:"environment"`
	SharedVolumes SharedVolumeAttachmentSlice `yaml:"shared_volumes"`
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

func (p *Pod) Describe() (string, error) {
	out, err := yaml.Marshal(p)
	if err != nil {
		return "", err
	}

	return string(out), nil

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
func (p *Pod) Exec(command []string, interactive bool, tty bool) (io.Reader, error) {
	// req := *api_client.NewProjectProjectIdPodPodIdExecPostRequest(command)
	// _, r, err := C.PodApi.ProjectProjectIdPodPodIdExecPost(context.Background(), p.Project.ID, p.ID).
	// 	Account(viper.GetString("account")).
	// 	Interactive(interactive).
	// 	Tty(tty).
	// 	Command(command).Execute()
	//fmt.Println("resp", resp, "err", err)

	// Have to do the above by hand to be able to send the stdin over http
	u := url.URL{
		Scheme: C.GetConfig().Scheme,
		Host:   C.GetConfig().Host,
		Path:   "/v0/project/" + p.Project.ID + "/pod/" + p.ID + "/exec",
	}
	q := u.Query()
	q.Set("interactive", strconv.FormatBool(interactive))
	q.Set("tty", strconv.FormatBool(tty))
	for _, cmd := range command {
		q.Add("command", cmd)
	}
	u.RawQuery = q.Encode()
	//fmt.Println("u", u.String())
	// stdin := os.Stdin
	// if !interactive {
	// 	stdin = nil
	// }
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", viper.GetString("api_key"))
	req.Header.Set("Account", viper.GetString("account"))

	r, err := C.GetConfig().HTTPClient.Do(req)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}

	return r.Body, err
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
