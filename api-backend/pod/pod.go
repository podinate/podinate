package pod

import (
	"context"
	"database/sql"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/johncave/podinate/api-backend/apierror"
	"github.com/johncave/podinate/api-backend/config"
	api "github.com/johncave/podinate/api-backend/go"
	lh "github.com/johncave/podinate/api-backend/loghandler"
	"github.com/johncave/podinate/api-backend/project"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ActionCreate   = "pod:create"
	ActionView     = "pod:view"
	ActionUpdate   = "pod:update"
	ActionDelete   = "pod:delete"
	ActionViewLogs = "pod:view"
)

type Pod struct {
	Uuid        string
	ID          string
	Name        string
	Image       string
	Tag         string
	Environment EnvironmentSlice
	Services    ServiceSlice
	Volumes     VolumeSlice
	Status      string // Creating, OK, Down
	Count       int
	Ram         int
	Project     *project.Project

	// TODO - add CPU requests / limits
}

func GetByID(ctx context.Context, theProject *project.Project, id string) (Pod, *apierror.ApiError) {
	p := Pod{}
	dberr := config.DB.QueryRow("SELECT uuid, id, name, image, tag, environment FROM project_pods WHERE id = $1 AND project_uuid = $2", id, theProject.Uuid).Scan(&p.Uuid, &p.ID, &p.Name, &p.Image, &p.Tag, &p.Environment)
	if dberr != nil && dberr != sql.ErrNoRows {
		log.Println("DB error getting pod", dberr)
		return Pod{}, &apierror.ApiError{Code: http.StatusInternalServerError, Message: dberr.Error()}
	}

	if dberr == sql.ErrNoRows {
		lh.Error(ctx, "Pod not found", "project", theProject, "id", id)
		return Pod{}, &apierror.ApiError{Code: http.StatusNotFound, Message: "Pod not found"}
	}

	p.Project = theProject

	// Get the services for the pod
	err := p.loadServices()
	if err != nil {
		return Pod{}, apierror.New(http.StatusInternalServerError, "Error Loading services"+err.Error())
	}
	err = p.loadVolumes()
	if err != nil {
		return Pod{}, apierror.NewWithError(http.StatusInternalServerError, "error loading volumes", err)
	}

	// Get the status of the pod from kubernetes
	ss, err := getKubesStatefulSet(theProject, id)
	if err != nil {
		return Pod{}, apierror.NewWithError(http.StatusInternalServerError, "error getting Kubernetes resources", err)
	}

	out := p

	split := strings.Split(ss.Spec.Template.Spec.Containers[0].Image, ":")
	if len(split) > 1 {
		out.Tag = split[1]
		out.Image = split[0]
	}

	// Get all the pods in the statefulset
	options := metav1.ListOptions{
		LabelSelector: "podinate.com/pod=" + id,
	}
	// var kpods *corev1.PodList

	var kpods *corev1.PodList
	for tries := 0; tries < 5; tries++ {
		// If we create a Pod then immediately try to retrieve it from Kubernetes
		// there may be a race condition IE in tests.
		// Don't ask me about how this wasted an entire Friday evening.
		kpods, err = config.Client.CoreV1().Pods(theProject.GetNamespaceName()).List(ctx, options)
		if err != nil {
			lh.Error(ctx, "Error getting pods from Kubernetes", "error", err)
			return Pod{}, apierror.NewWithError(http.StatusInternalServerError, "Error getting pods from Kubernetes", err)
		}
		//lh.Debug(ctx, "Got all pods from Kubernetes", "pods", kpods, "length", len(kpods.Items), "options", options, "namespace", theProject.GetNamespaceName())
		if len(kpods.Items) > 0 {
			break
		}
		if tries > 5 {
			lh.Error(ctx, "No pods found after 5 tries", "pods", kpods, "length", len(kpods.Items), "options", options, "namespace", theProject.GetNamespaceName())
			return Pod{}, apierror.New(http.StatusInternalServerError, "No pods found after 5 tries")
		}
		time.Sleep(10 * time.Millisecond)
	}

	out.Status = string(kpods.Items[0].Status.Phase)

	//lh.Debug(ctx, "Pod got by ID", "pod", out)

	return out, nil

}

func GetByProject(ctx context.Context, theProject *project.Project, page int32, limit int32) ([]Pod, *apierror.ApiError) {
	if limit < 1 || limit > 125 {
		limit = 25
	}

	// Get the uuid of all pods in the project

	rows, err := config.DB.Query("SELECT uuid, id FROM project_pods WHERE project_uuid = $1", theProject.Uuid)
	if err != nil {
		return nil, apierror.New(http.StatusInternalServerError, "Could not retrieve pods")
	}
	defer rows.Close()

	// Read all the pods for the project
	pods := make([]Pod, 0)
	for rows.Next() {
		var uuid string
		var id string
		err = rows.Scan(&uuid, &id)
		pod, err := GetByID(ctx, theProject, id)
		if err != nil {
			lh.Log.Errorw("Error getting pod", "error", err)
			return nil, apierror.New(http.StatusInternalServerError, "Could not retrieve pods")
		}
		pods = append(pods, pod)
	}

	lh.Info(ctx, "Pods retrieved", "project", theProject, "count", len(pods), "pods", pods)

	return pods, nil
	// rows, err := config.DB.Query("SELECT uuid, id, name, image, tag, environment FROM project_pods WHERE project_uuid = $1 OFFSET $2 LIMIT $3", theProject.Uuid, page, limit)
	// if err != nil {
	// 	return nil, apierror.New(http.StatusInternalServerError, "Could not retrieve pods")
	// }
	// defer rows.Close()
	// // Read all the pods for the project
	// pods := make([]Pod, 0)
	// for rows.Next() {
	// 	var pod Pod
	// 	err = rows.Scan(&pod.Uuid, &pod.ID, &pod.Name, &pod.Image, &pod.Tag, &pod.Environment)
	// 	if err != nil {
	// 		log.Println("DB error reading pods", err)
	// 		return nil, apierror.New(http.StatusInternalServerError, "Could not retrieve pods")
	// 	}
	// 	pod.Project = theProject

	// 	// Get the services for the pod
	// 	pod.loadServices()

	// 	// Get the status of the pod from kubernetes
	// 	dep, err := getKubesDeployment(theProject, pod.ID)
	// 	if err != nil {
	// 		return nil, apierror.New(http.StatusInternalServerError, err.Error())
	// 	}

	// 	status := "Creating"
	// 	if dep.Status.AvailableReplicas == dep.Status.Replicas {
	// 		status = "OK"
	// 	} else if dep.Status.AvailableReplicas == 0 {
	// 		status = "Down"
	// 	}

	// 	pod.Status = status

	// 	pods = append(pods, pod)
	// }
	// return pods, nil
}

// Create performs the initial registration of a pod in the database and the kubernetes cluster
func Create(ctx context.Context, theProject *project.Project, requestedPod api.Pod) (Pod, *apierror.ApiError) {
	// Check if the pod already exists
	uuid := ""
	dberr := config.DB.QueryRow("SELECT uuid FROM project_pods WHERE id = $1 AND project_uuid = $2", requestedPod.Id, theProject.Uuid).Scan(&uuid)
	// Errors other than no rows is a problem
	// In good state
	// dberr != nil
	// dberr == sql.ErrNoRows
	if dberr != nil && dberr != sql.ErrNoRows {
		log.Println("DB error checking if pod exists", dberr)
		return Pod{}, &apierror.ApiError{Code: http.StatusInternalServerError, Message: dberr.Error()}
	}
	// Handle a conflicting pod existing
	if dberr != sql.ErrNoRows {
		return Pod{}, &apierror.ApiError{Code: http.StatusConflict, Message: "A pod with this ID already exists"}
	}

	// Start creating the pod
	out := Pod{
		Uuid:        uuid, // Note: UUID blank (should be fine?)
		ID:          requestedPod.Id,
		Name:        requestedPod.Name,
		Image:       requestedPod.Image,
		Tag:         requestedPod.Tag,
		Project:     theProject,
		Status:      "Creating",
		Environment: EnvVarFromAPIMany(requestedPod.Environment),
		Services:    servicesFromAPI(requestedPod.Services),
	}

	// Ensure namespace exists
	ns, genericErr := out.ensureNamespace()
	//ns, err := createKubesNamespace(theProject.GetResourceID())
	if genericErr != nil {
		return Pod{}, apierror.New(http.StatusInternalServerError, genericErr.Error())
	}

	err := createKubesDeployment(ns, theProject, requestedPod)
	if err != nil {
		lh.Error(ctx, "Error creating kubes deployment", "err", err)
		return Pod{}, err
	}

	svcErr := out.ensureServices(ctx)
	lh.Debug(ctx, "Ensure services", "err", svcErr)
	if svcErr != nil {
		return Pod{}, svcErr
	}

	// Create the pod in the database
	dberr = config.DB.QueryRow("INSERT INTO project_pods (uuid, id, name, image, tag, project_uuid, environment) VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6) RETURNING uuid", requestedPod.Id, requestedPod.Name, requestedPod.Image, requestedPod.Tag, theProject.Uuid, EnvVarSliceFromAPI(requestedPod.Environment)).Scan(&uuid)
	// Check if insert was successful
	if dberr != nil {
		lh.Error(ctx, "DB error creating pod", dberr)
		return Pod{}, &apierror.ApiError{Code: http.StatusInternalServerError, Message: dberr.Error()}
	}

	// Add the services to the database if it exists
	if len(requestedPod.Services) > 0 {
		for _, service := range requestedPod.Services {
			_, dberr = config.DB.Exec("INSERT INTO pod_services (uuid, pod_uuid, name, port, target_port, protocol, domain_name) VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6)", uuid, service.Name, service.Port, service.TargetPort, service.Protocol, service.DomainName)
			if dberr != nil {
				lh.Error(ctx, "DB error creating pod service", dberr)
				out.Delete(ctx)
				return Pod{}, &apierror.ApiError{Code: http.StatusInternalServerError, Message: dberr.Error()}
			}
		}
	}

	if len(requestedPod.Volumes) > 0 {
		for _, volume := range requestedPod.Volumes {
			_, dberr = config.DB.Exec("INSERT INTO pod_volumes (uuid, pod_uuid, name, mount_path, size) VALUES (gen_random_uuid(), $1, $2, $3, $4)", uuid, volume.Name, volume.MountPath, volume.Size)
			if dberr != nil {
				lh.Error(ctx, "DB error creating pod volume", dberr)
				out.Delete(ctx)
				return Pod{}, &apierror.ApiError{Code: http.StatusInternalServerError, Message: dberr.Error()}
			}
		}
	}

	// genericErr = out.loadServices()
	// if err != nil {
	// 	out.Delete(ctx)
	// 	return Pod{}, apierror.New(http.StatusInternalServerError, err.Error())
	// }
	// genericErr = out.loadVolumes()
	// if err != nil {
	// 	out.Delete(ctx)
	// 	return Pod{}, apierror.New(http.StatusInternalServerError, err.Error())
	// }

	// Everything is created - we can load the pod from the database again
	out, apiErr := GetByID(ctx, theProject, requestedPod.Id)
	if apiErr != nil {
		return Pod{}, apierror.New(http.StatusInternalServerError, apiErr.Error())
	}

	return out, nil
}

func (p *Pod) ToAPI() api.Pod {
	lh.Log.Debug("HELLO WTF", "pod", p)

	out := api.Pod{
		Id:          p.ID,
		Name:        p.Name,
		Image:       p.Image,
		Tag:         p.Tag,
		Status:      p.Status,
		Environment: EnvVarToAPIMany(p.Environment),
		ResourceId:  p.GetResourceID(),
		Services:    ServicesToAPI(p.Services),
		Volumes:     p.Volumes.ToAPI(),
	}

	lh.Log.Debug("Converted pod to API", "pod", p, "apiPod", out)

	return out

}

// Delete removes a pod from the database and the kubernetes cluster
func (p *Pod) Delete(ctx context.Context) *apierror.ApiError {

	// The kubes logic
	err := deleteKubesDeployment(*p)
	if err != nil {
		return apierror.New(http.StatusInternalServerError, err.Error())
	}

	// Delete any services from the database
	_, dberr := config.DB.Exec("DELETE FROM pod_services WHERE pod_uuid = $1", p.Uuid)
	if dberr != nil {
		lh.Error(ctx, "DB error deleting pod services", "err", dberr)
		return &apierror.ApiError{Code: http.StatusInternalServerError, Message: dberr.Error()}
	}

	// Delete the pod from the database
	_, dberr = config.DB.Exec("DELETE FROM project_pods WHERE project_uuid = $1 AND id = $2", p.Project.Uuid, p.ID)
	// Check if delete was successful

	if dberr != nil {
		log.Println("DB error deleting pod", dberr)
		return &apierror.ApiError{Code: http.StatusInternalServerError, Message: dberr.Error()}
	}

	// Delete the services from the database
	_, dberr = config.DB.Exec("DELETE FROM pod_services WHERE pod_uuid = $1", p.Uuid)
	if dberr != nil {
		lh.Error(ctx, "DB error deleting pod services", "err", dberr)
		return &apierror.ApiError{Code: http.StatusInternalServerError, Message: dberr.Error()}
	}

	lh.Info(ctx, "pod deleted", "pod", p)
	//
	return nil
}

func (p *Pod) Update(ctx context.Context, requested api.Pod) *apierror.ApiError {
	// TODO - Validate the requestedPod

	// Check if the pod already exists
	uuid := ""
	dberr := config.DB.QueryRow("SELECT uuid FROM project_pods WHERE id = $1 AND project_uuid = $2", requested.Id, p.Project.Uuid).Scan(&uuid)
	// Errors other than no rows is a problem
	// In good state
	// dberr != nil
	// dberr == sql.ErrNoRows
	if dberr != nil {
		log.Println("DB error checking if pod exists", dberr)
		return &apierror.ApiError{Code: http.StatusInternalServerError, Message: dberr.Error()}
	}

	// Update the elements of p that can change
	p.Name = requested.Name
	p.Image = requested.Image
	p.Tag = requested.Tag
	p.Status = "Updating"
	p.Environment = EnvVarFromAPIMany(requested.Environment)
	p.Services = servicesFromAPI(requested.Services)

	// Start creating the pod

	err := updateKubesDeployment(*p, requested)
	if err != nil {
		return apierror.New(http.StatusInternalServerError, err.Error())
	}
	err = p.ensureServices(ctx)
	if err != nil {
		return apierror.New(http.StatusInternalServerError, err.Error())
	}

	// Update the pod in the database
	dberr = config.DB.QueryRow("UPDATE project_pods SET name = $1, image = $2, tag = $3, environment = $4 WHERE uuid = $5 RETURNING uuid", requested.Name, requested.Image, requested.Tag, EnvVarSliceFromAPI(requested.Environment), p.Uuid).Scan(&uuid)
	// Check if insert was successful
	if dberr != nil {
		log.Println("DB error creating pod", dberr)
		return &apierror.ApiError{Code: http.StatusInternalServerError, Message: dberr.Error()}
	}

	//Update or create services
	for _, service := range requested.Services {
		// See if the service already exists
		exists, err := p.serviceExists(service.Name)
		if err != nil {
			lh.Error(ctx, "DB error checking if service exists", dberr)
			return &apierror.ApiError{Code: http.StatusInternalServerError, Message: dberr.Error()}
		}
		if exists {
			_, dberr = config.DB.Exec("UPDATE pod_services SET port = $3, target_port = $4, protocol = $5, domain_name = $6 WHERE pod_uuid = $1 AND name = $2", p.Uuid, service.Name, service.Port, service.TargetPort, service.Protocol, service.DomainName)
			if dberr != nil {
				lh.Error(ctx, "DB error updating pod service", dberr)
				return &apierror.ApiError{Code: http.StatusInternalServerError, Message: dberr.Error()}
			}
			continue
		}

		// Finally, try create the service
		_, dberr = config.DB.Exec("INSERT INTO pod_services (uuid, pod_uuid, name, port, target_port, protocol, domain_name) VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6)", p.Uuid, service.Name, service.Port, service.TargetPort, service.Protocol, service.DomainName)
		// If the service already exists, update it

		if dberr != nil {
			lh.Error(ctx, "DB error creating pod service", dberr)
			return &apierror.ApiError{Code: http.StatusInternalServerError, Message: dberr.Error()}
		}
	}

	// Update or create pod volumes
	for _, volume := range requested.Volumes {
		// See if the volume already exists
		exists, err := p.volumeExists(volume.Name)
		if err != nil {
			lh.Error(ctx, "DB error checking if volume exists", dberr)
			return &apierror.ApiError{Code: http.StatusInternalServerError, Message: dberr.Error()}
		}
		if exists {
			_, dberr = config.DB.Exec("UPDATE pod_volumes SET mount_path = $3, size = $4 WHERE pod_uuid = $1 AND name = $2", p.Uuid, volume.Name, volume.MountPath, volume.Size)
			if dberr != nil {
				lh.Error(ctx, "DB error updating pod volume", dberr)
				return &apierror.ApiError{Code: http.StatusInternalServerError, Message: dberr.Error()}
			}
			continue
		} else {
			// Finally, try create the volume
			_, dberr = config.DB.Exec("INSERT INTO pod_volumes (uuid, pod_uuid, name, mount_path, size) VALUES (gen_random_uuid(), $1, $2, $3, $4)", p.Uuid, volume.Name, volume.MountPath, volume.Size)
			if dberr != nil {
				lh.Error(ctx, "DB error creating pod volume", dberr)
				return &apierror.ApiError{Code: http.StatusInternalServerError, Message: dberr.Error()}
			}
		}
	}

	return nil
}

func (p *Pod) GetLogsBuffer(ctx context.Context, lines int, follow bool) (io.ReadCloser, *apierror.ApiError) {
	// Get the logs from kubernetes

	l := int64(lines)
	podLogOpts := corev1.PodLogOptions{
		TailLines: &l,
		Follow:    follow,
	}

	req := config.Client.CoreV1().Pods(p.Project.GetNamespaceName()).GetLogs(p.ID+"-0", &podLogOpts)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return nil, apierror.NewWithError(http.StatusInternalServerError, "error getting logs from kubernetes", err)
	}

	return podLogs, nil
}

func (p Pod) GetResourceID() string {
	return p.Project.GetResourceID() + "/pod:" + p.ID
}

func (p *Pod) getAnnotations() map[string]string {
	return map[string]string{
		"podinate.io/project": p.Project.GetResourceID(),
		"podinate.io/pod":     p.GetResourceID(),
	}

}
