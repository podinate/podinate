package pod

import (
	"context"
	"fmt"
	"log"
	"os"

	api "github.com/johncave/podinate/api-backend/go"
	"github.com/johncave/podinate/api-backend/project"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func callKubes() {
	fmt.Println("Get Kubernetes pods")

	clientset, err := getKubesClient()
	log.Println("Getting pods...")

	pods, err := clientset.CoreV1().
		Pods("kube-system").
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("error getting pods: %v\n", err)
		os.Exit(1)
	}
	for _, pod := range pods.Items {
		fmt.Printf("Pod name: %s\n", pod.Name)
	}

	nsList, err := clientset.CoreV1().
		Namespaces().
		List(context.Background(), metav1.ListOptions{})
	//checkErr(err)
	fmt.Println(err)

	for _, n := range nsList.Items {
		fmt.Printf("Namespace: %s\n", n.Name)
	}

}

func createKubesNamespace(name string) (*corev1.Namespace, error) {
	fmt.Println("Create Kubernetes namespace")

	clientset, err := getKubesClient()
	if err != nil {
		log.Printf("error getting kubernetes client: %v\n", err)
		return nil, err
	}

	nsSpec := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	ns, err := clientset.CoreV1().
		Namespaces().
		Create(context.Background(), nsSpec, metav1.CreateOptions{})
	if errors.IsAlreadyExists(err) {
		// Get the ns instead
		ns, err := clientset.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
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

// getKubesDeployment returns a deployment in the specified namespace.
func getKubesDeployment(theProject project.Project, id string) (*appsv1.Deployment, error) {
	fmt.Println("Get Kubernetes deployment")

	clientset, err := getKubesClient()
	if err != nil {
		log.Printf("error getting kubernetes client: %v\n", err)
		return nil, err
	}

	deployment, err := clientset.AppsV1().
		Deployments(theProject.Account.ID+"-project-"+theProject.ID).
		Get(context.Background(), id, metav1.GetOptions{})

	if err != nil {
		fmt.Printf("error getting deployment: %v\n", err)
		return nil, err
	}
	return deployment, nil
}

// createKubesDeployment creates a deployment in the specified namespace.
func createKubesDeployment(inns *corev1.Namespace, theProject project.Project, requested api.Pod) error {
	fmt.Println("Create Kubernetes deployment")

	clientset, err := getKubesClient()
	if err != nil {
		log.Printf("error getting kubernetes client: %v\n", err)
		return err

	}

	deploymentSpec := getDeploymentSpec(theProject, requested)

	_, err = clientset.AppsV1().
		Deployments(inns.Name).
		Create(context.Background(), deploymentSpec, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("error creating deployment: %v\n", err)
		return err
	}
	return nil
}

// updateKubesDeployment updates a deployment in the specified namespace.
func updateKubesDeployment(thePod Pod, requested api.Pod) error {
	fmt.Println("Update Kubernetes deployment")

	clientset, err := getKubesClient()
	if err != nil {
		log.Printf("error getting kubernetes client: %v\n", err)
		return err
	}

	deploymentSpec := getDeploymentSpec(thePod.Project, thePod.ToAPI())

	_, err = clientset.AppsV1().
		Deployments(thePod.Project.Account.ID+"-project-"+thePod.Project.ID).
		Update(context.Background(), deploymentSpec, metav1.UpdateOptions{})
	if err != nil {
		fmt.Printf("error updating deployment: %v\n", err)
		return err
	}
	return nil
}

// getDeploymentSpec returns a deployment spec for the specified pod.
func getDeploymentSpec(theProject project.Project, requested api.Pod) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: requested.Id,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: func(val int32) *int32 { return &val }(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"podinate": requested.Id,
					"project":  theProject.ID,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"podinate": requested.Id,
						"project":  theProject.ID,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  requested.Id,
							Image: requested.Image + ":" + requested.Tag,
							// TODO: Figure out what to do about ports
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
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

	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("error getting Kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("error getting Kubernetes clientset: %v", err)
	}

	return clientset, nil
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
		Deployments(thePod.Project.Account.ID+"-project-"+thePod.Project.ID).
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
