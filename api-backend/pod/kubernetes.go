package pod

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/johncave/podinate/api-backend/apierror"
	"github.com/johncave/podinate/api-backend/config"
	api "github.com/johncave/podinate/api-backend/go"
	lh "github.com/johncave/podinate/api-backend/loghandler"
	"github.com/johncave/podinate/api-backend/project"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

// func callKubes() {
// 	fmt.Println("Get Kubernetes pods")

// 	clientset, err := getKubesClient()
// 	log.Println("Getting pods...")

// 	pods, err := clientset.CoreV1().
// 		Pods("kube-system").
// 		List(context.Background(), metav1.ListOptions{})
// 	if err != nil {
// 		fmt.Printf("error getting pods: %v\n", err)
// 		os.Exit(1)
// 	}
// 	for _, pod := range pods.Items {
// 		fmt.Printf("Pod name: %s\n", pod.Name)
// 	}

// 	nsList, err := clientset.CoreV1().
// 		Namespaces().
// 		List(context.Background(), metav1.ListOptions{})
// 	//checkErr(err)
// 	fmt.Println(err)

// 	for _, n := range nsList.Items {
// 		fmt.Printf("Namespace: %s\n", n.Name)
// 	}

// }

func (p *Pod) ensureNamespace() (*corev1.Namespace, error) {
	fmt.Println("Create Kubernetes namespace")

	clientset, err := getKubesClient()
	if err != nil {
		log.Printf("error getting kubernetes client: %v\n", err)
		return nil, err
	}

	nsSpec := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: p.getNamespaceName(),
		},
	}
	ns, err := clientset.CoreV1().
		Namespaces().
		Create(context.Background(), nsSpec, metav1.CreateOptions{})
	if errors.IsAlreadyExists(err) {
		// Get the ns instead
		ns, err := clientset.CoreV1().Namespaces().Get(context.Background(), p.getNamespaceName(), metav1.GetOptions{})
		if err != nil {
			fmt.Printf("error getting existing kubernetes namespace: %v\n", err)
			return ns, err
		}
		return ns, nil
	}
	if err != nil {
		fmt.Printf("error creating kubernetes namespace: %v\n", err)
		return nil, err
	}
	return ns, nil
}

func (p *Pod) getNamespaceName() string {
	return p.Project.Account.ID + "-project-" + p.Project.ID
}

// getKubesDeployment returns a deployment in the specified namespace.
func getKubesStatefulSet(theProject *project.Project, id string) (*appsv1.StatefulSet, error) {
	fmt.Println("Get Kubernetes deployment")

	clientset, err := getKubesClient()
	if err != nil {
		log.Printf("error getting kubernetes client: %v\n", err)
		return nil, err
	}

	deployment, err := clientset.AppsV1().
		StatefulSets(theProject.Account.ID+"-project-"+theProject.ID).
		Get(context.Background(), id, metav1.GetOptions{})

	if err != nil {
		fmt.Printf("error getting deployment: %v\n", err)
		return nil, err
	}
	return deployment, nil
}

// createKubesDeployment creates a deployment in the specified namespace.
func createKubesDeployment(inns *corev1.Namespace, theProject *project.Project, requested api.Pod) *apierror.ApiError {
	fmt.Println("Create Kubernetes deployment")

	clientset, err := getKubesClient()
	if err != nil {
		log.Printf("error getting kubernetes client: %v\n", err)
		return apierror.New(500, "error getting kubernetes client")

	}

	statefulSet, apierr := getStatefulSetSpec(theProject, requested)
	if apierr != nil {
		lh.Log.Errorw("Error building Kubernetes spec", "error", apierr.Error())
		return apierr
	}

	_, err = clientset.AppsV1().
		StatefulSets(inns.Name).
		Create(context.Background(), statefulSet, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("error creating deployment: %v\n", err)
		return apierror.New(500, "error creating deployment: "+err.Error())
	}
	return nil
}

// updateKubesDeployment updates a deployment in the specified namespace.
func updateKubesDeployment(thePod Pod, requested api.Pod) *apierror.ApiError {
	fmt.Println("Update Kubernetes deployment")

	clientset, err := getKubesClient()
	if err != nil {
		log.Printf("error getting kubernetes client: %v\n", err)
		return apierror.New(500, "error getting kubernetes client")
	}

	statefulSet, apierr := getStatefulSetSpec(thePod.Project, thePod.ToAPI())
	if err != nil {
		lh.Log.Errorw("error getting statefulset spec to update", "error", apierr.Error())
		return apierr
	}

	_, err = clientset.AppsV1().
		StatefulSets(thePod.Project.Account.ID+"-project-"+thePod.Project.ID).
		Update(context.Background(), statefulSet, metav1.UpdateOptions{})
	if err != nil {
		fmt.Printf("error updating deployment: %v\n", err)
		return &apierror.ApiError{Code: 500, Message: "error updating deployment " + err.Error()}
	}

	return nil
}

func getStatefulSetSpec(theProject *project.Project, requested api.Pod) (*appsv1.StatefulSet, *apierror.ApiError) {
	out := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: requested.Id,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: func(val int32) *int32 { return &val }(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"podinate.com/pod":     requested.Id,
					"podinate.com/project": theProject.ID,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"podinate.com/pod":     requested.Id,
						"podinate.com/project": theProject.ID,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    requested.Id,
							Image:   requested.Image + ":" + requested.Tag,
							Command: requested.Command,
							Args:    requested.Arguments,
						},
					},
				},
			},
		},
	}

	// // Add command to the pod spec
	// if requested.Command != nil {
	// 	lh.Log.Debugw("Adding command to pod spec", "command", requested.Command)
	// 	out.Spec.Template.Spec.Containers[0].Command = requested.Command
	// } else {
	// 	lh.Log.Debugw("No command to add to pod spec", "command", requested.Command)
	// }

	// // Add arguments to the pod spec
	// if requested.Arguments != nil {

	// Add environment variables to the pod spec
	for _, envVar := range requested.Environment {
		out.Spec.Template.Spec.Containers[0].Env = append(out.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{
			Name:  envVar.Key,
			Value: envVar.Value,
		})
	}

	for _, volume := range requested.Volumes {
		// Add volume mounts
		out.Spec.Template.Spec.Containers[0].VolumeMounts = append(out.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
			Name:      volume.Name,
			MountPath: volume.MountPath,
		})

		// Add volume claim templates
		newPVC := corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      volume.Name,
				Namespace: theProject.Account.ID + "-project-" + theProject.ID,
				Annotations: map[string]string{
					"volumeType": "local",
				},
				Labels: map[string]string{
					"podinate.com/pod":     requested.Id,
					"podinate.com/project": theProject.ID,
				},
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					"ReadWriteOnce",
				},
				//StorageClassName: func(val string) *string { return &val }("local-path"),
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						"storage": func(val string) resource.Quantity { return resource.MustParse(val) }(fmt.Sprintf("%dGi", volume.Size)),
					},
				},
			},
		}

		// Add the storageclass if exists
		if volume.Class != "" {
			// Check given SC exists
			storageClass, err := getStorageClass(volume.Class)
			if err != nil {
				lh.Log.Errorw("error getting storage class", "error", err.Error())
				return nil, apierror.NewWithError(http.StatusBadRequest, "error checking storage class exists", err)
			}
			newPVC.Spec.StorageClassName = &storageClass.Name
		}

		out.Spec.VolumeClaimTemplates = append(out.Spec.VolumeClaimTemplates, newPVC)

	}

	lh.Log.Infow("StatefulSet spec generated", "statefulset", out)

	return out, nil
}

// getKubesClient returns a Kubernetes clientset.
func getKubesClient() (*kubernetes.Clientset, error) {
	// userHomeDir, err := os.UserHomeDir()
	// if err != nil {
	// 	return nil, fmt.Errorf("error getting user home dir: %v", err)
	// }
	// kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	// fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

	// kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)

	// kubeConfig, err := rest.InClusterConfig()
	// if err != nil {
	// 	return nil, fmt.Errorf("error getting Kubernetes config: %v", err)
	// }

	// clientset, err := kubernetes.NewForConfig(kubeConfig)
	// if err != nil {
	// 	return nil, fmt.Errorf("error getting Kubernetes clientset: %v", err)
	// }

	return config.Client, nil
}

// deleteKubesDeployment deletes a deployment in the specified namespace.
func deleteKubesDeployment(thePod Pod) error {
	fmt.Println("Delete Kubernetes deployment")

	clientset, err := getKubesClient()
	if err != nil {
		log.Printf("error getting kubernetes client: %v\n", err)
		return err
	}

	deletePolicy := metav1.DeletePropagationForeground
	err = clientset.AppsV1().
		StatefulSets(thePod.Project.Account.ID+"-project-"+thePod.Project.ID).
		Delete(context.Background(), thePod.ID, metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
	if err != nil {
		fmt.Printf("error deleting deployment: %v\n", err)
		return err
	}
	return nil

	// TODO - Figure out how to clean up unused namespaces
}
