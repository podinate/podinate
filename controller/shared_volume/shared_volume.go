package shared_volume

import (
	"context"
	"database/sql"
	"maps"
	"net/http"
	"strconv"

	"github.com/johncave/podinate/controller/apierror"
	"github.com/johncave/podinate/controller/config"
	api "github.com/johncave/podinate/controller/go"
	"github.com/johncave/podinate/controller/iam"
	lh "github.com/johncave/podinate/controller/loghandler"
	"github.com/johncave/podinate/controller/project"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ActionCreate = "sharedvolume:create"
	ActionView   = "sharedvolume:view"
	ActionUpdate = "sharedvolume:update"
	ActionDelete = "sharedvolume:delete"
)

type SharedVolume struct {
	Uuid    string
	ID      string
	Name    *string
	Project *project.Project
	Size    int
	Class   *string
}

// GetByID returns a shared volume by ID
func GetByID(ctx context.Context, theProject *project.Project, id string) (*SharedVolume, *apierror.ApiError) {
	p := &SharedVolume{}
	dberr := config.DB.QueryRowContext(ctx, "SELECT uuid, id, name, size, class FROM shared_volumes WHERE id = $1 AND project_uuid = $2", id, theProject.Uuid).Scan(&p.Uuid, &p.ID, &p.Name, &p.Size, &p.Class)
	if dberr == sql.ErrNoRows {
		return nil, apierror.New(http.StatusNotFound, "Shared volume not found")

	} else if dberr != nil {
		return nil, apierror.NewWithError(http.StatusInternalServerError, "Error getting shared volume", dberr)
	}

	p.Project = theProject

	if !iam.RequestorCan(ctx, &p.Project.Account, p, ActionView) {
		return nil, apierror.New(http.StatusForbidden, "You do not have permission to view this shared volume")
	}

	return p, nil
}

// GetByProject returns all shared volumes by project
func GetByProject(ctx context.Context, theProject *project.Project) ([]*SharedVolume, *apierror.ApiError) {
	rows, err := config.DB.QueryContext(ctx, "SELECT id, name, size, class FROM shared_volumes WHERE project_uuid = $1", theProject.Uuid)
	if err != nil {
		return nil, apierror.NewWithError(http.StatusInternalServerError, "Error getting shared volumes", err)
	}
	defer rows.Close()

	var volumes []*SharedVolume
	for rows.Next() {
		p := &SharedVolume{}
		err := rows.Scan(&p.ID, &p.Name, &p.Size, &p.Class)
		if err != nil {
			return nil, apierror.NewWithError(http.StatusInternalServerError, "Error scanning shared volume", err)
		}
		p.Project = theProject
		if iam.RequestorCan(ctx, &p.Project.Account, p, ActionView) {
			volumes = append(volumes, p)
		}
	}

	if len(volumes) == 0 {
		if !iam.RequestorCan(ctx, &theProject.Account, theProject, project.ActionView) {
			return nil, apierror.New(http.StatusForbidden, "You do not have permission to view shared volumes in this project")
		}
		return volumes, nil
	}

	return volumes, nil
}

// Create creates a new shared volume
func Create(ctx context.Context, theProject *project.Project, in api.SharedVolume) (*SharedVolume, *apierror.ApiError) {
	new := fromAPI(in)
	new.Project = theProject

	lh.Debug(ctx, "Creating shared volume", "volume", new, "project", theProject)
	if !iam.RequestorCan(ctx, &theProject.Account, new, ActionCreate) {
		return nil, apierror.New(http.StatusForbidden, "You do not have permission to create a shared volume")
	}

	apierr := new.ensureKubernetesResources(ctx)
	if apierr != nil {
		return nil, apierr
	}

	err := config.DB.QueryRowContext(ctx, "INSERT INTO shared_volumes (id, name, size, class, project_uuid) VALUES ($1, $2, $3, $4, $5) returning uuid", in.Id, in.Name, in.Size, in.Class, theProject.Uuid).Scan(&new.Uuid)
	if err != nil {
		return nil, apierror.NewWithError(http.StatusInternalServerError, "Error creating shared volume", err)
	}

	new.Project = theProject

	return new, nil
}

// Update updates a shared volume
func (v *SharedVolume) Update(ctx context.Context, in api.SharedVolume) *apierror.ApiError {
	if !iam.RequestorCan(ctx, &v.Project.Account, v, ActionUpdate) {
		return apierror.New(http.StatusForbidden, "You do not have permission to update this shared volume")
	}

	_, err := config.DB.ExecContext(ctx, "UPDATE shared_volumes SET name = $1, size = $2, class = $3 WHERE id = $4 AND project_uuid = $5", in.Name, in.Size, in.Class, v.ID, v.Project.Uuid)
	if err != nil {
		return apierror.NewWithError(http.StatusInternalServerError, "Error updating shared volume", err)
	}

	out, apierr := GetByID(ctx, v.Project, v.ID)
	if apierr != nil {
		return apierr
	}

	*v = *out
	return nil
}

// Delete deletes a shared volume
func (v *SharedVolume) Delete(ctx context.Context) *apierror.ApiError {
	if !iam.RequestorCan(ctx, &v.Project.Account, v, ActionDelete) {
		return apierror.New(http.StatusForbidden, "You do not have permission to delete this shared volume")
	}

	// Delete from Kubernetes
	err := config.Client.CoreV1().PersistentVolumeClaims(v.Project.GetNamespaceName()).Delete(ctx, v.ID, metav1.DeleteOptions{})
	if err != nil {
		return apierror.NewWithError(http.StatusInternalServerError, "Error deleting shared volume", err)
	}

	_, err = config.DB.ExecContext(ctx, "DELETE FROM shared_volumes WHERE id = $1 AND project_uuid = $2", v.ID, v.Project.Uuid)
	if err != nil {
		return apierror.NewWithError(http.StatusInternalServerError, "Error deleting shared volume", err)
	}

	lh.Info(ctx, "Shared volume deleted", "volume", v)
	return nil
}

// ValidateNewRequested validates a new shared volume
func ValidateNewRequested(ctx context.Context, in api.SharedVolume) *apierror.ApiError {
	if in.Id == "" {
		return apierror.New(http.StatusBadRequest, "ID is required")
	}

	if len(in.Id) > 63 {
		return apierror.New(http.StatusBadRequest, "ID must be 63 characters or less")
	}

	if len(in.Name) > 63 {
		return apierror.New(http.StatusBadRequest, "Name must be 63 characters or less")
	}

	if in.Size == 0 {
		return apierror.New(http.StatusBadRequest, "Size is required")
	}

	if in.Class != "" {

		// Try to get the storageclass
		_, err := config.Client.StorageV1().StorageClasses().Get(ctx, in.Class, metav1.GetOptions{})
		if err != nil {
			return apierror.New(http.StatusBadRequest, "Storage class "+in.Class+" does not exist")
		}
	}

	return nil
}

// ensureKubernetesResources ensures the shared volume exists in kubernetes
func (v *SharedVolume) ensureKubernetesResources(ctx context.Context) *apierror.ApiError {

	pvc, apierr := v.getkubernetesSpec()
	if apierr != nil {
		return apierr
	}

	// First try to get the PVC
	_, err := config.Client.CoreV1().PersistentVolumeClaims(v.Project.GetNamespaceName()).Get(ctx, v.ID, metav1.GetOptions{})
	if err != nil {
		// PVC doesn't exist, create it
		_, err := config.Client.CoreV1().PersistentVolumeClaims(v.Project.GetNamespaceName()).Create(ctx, pvc, metav1.CreateOptions{})
		if err != nil {
			lh.Error(ctx, "Error creating shared volume", "error", err, "volume", v)
			return apierror.NewWithError(http.StatusInternalServerError, "Error creating shared volume", err)
		}
		lh.Info(ctx, "Shared volume created", "volume", v)
	} else {
		// PVC exists, update it
		_, err := config.Client.CoreV1().PersistentVolumeClaims(v.Project.GetNamespaceName()).Update(ctx, pvc, metav1.UpdateOptions{})
		if err != nil {
			lh.Error(ctx, "Error updating shared volume", "error", err, "volume", v)
			return apierror.NewWithError(http.StatusInternalServerError, "Error updating shared volume", err)
		}
		lh.Info(ctx, "Shared volume updated", "volume", v)
	}

	return nil

}

// getkubernetesSpec returns a kubernetes spec for the shared volume
func (v *SharedVolume) getkubernetesSpec() (*apiv1.PersistentVolumeClaim, *apierror.ApiError) {
	out := &apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:        v.ID,
			Namespace:   v.Project.GetNamespaceName(),
			Labels:      v.GetLabels(),
			Annotations: v.GetAnnotations(),
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			AccessModes: []apiv1.PersistentVolumeAccessMode{
				apiv1.ReadWriteOnce,
			},
			Resources: apiv1.VolumeResourceRequirements{
				Requests: apiv1.ResourceList{
					apiv1.ResourceStorage: resource.MustParse(strconv.Itoa(v.Size) + "Gi"),
				},
			},
		},
	}

	if v.Class != nil {
		out.Spec.StorageClassName = v.Class
	}

	return out, nil
}

// fromAPI converts an API SharedVolume to a shared volume
func fromAPI(apiVolume api.SharedVolume) *SharedVolume {
	out := &SharedVolume{
		ID:   apiVolume.Id,
		Size: int(apiVolume.Size),
	}

	if apiVolume.Name != "" {
		out.Name = &apiVolume.Name
	}
	if apiVolume.Class != "" {
		out.Class = &apiVolume.Class
	}
	return out
}

// toAPI converts a shared volume to an API SharedVolume
func (v *SharedVolume) ToAPI() api.SharedVolume {
	out := api.SharedVolume{
		Id:   v.ID,
		Size: int32(v.Size),
	}

	if v.Name != nil {
		out.Name = *v.Name
	}
	if v.Class != nil {
		out.Class = *v.Class
	}
	return out
}

// GetLabels returns the labels for the shared volume
func (v *SharedVolume) GetLabels() map[string]string {
	out := map[string]string{
		"podinate.com/shared_volume_id": v.ID,
	}

	maps.Copy(out, v.Project.GetLabels())
	return out
}

// GetAnnotations returns the annotations for the shared volume
func (v *SharedVolume) GetAnnotations() map[string]string {
	out := map[string]string{}

	maps.Copy(out, v.Project.GetAnnotations())
	return out
}

// GetResourceID returns the resource ID for the shared volume
func (v *SharedVolume) GetResourceID() string {
	//return "penis"
	return v.Project.GetResourceID() + "/sharedvolume:" + v.ID
}
