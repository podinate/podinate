package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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

func createKubesNamespace(name string) error {
	fmt.Println("Create Kubernetes namespace")

	clientset, err := getKubesClient()
	if err != nil {
		log.Printf("error getting kubernetes client: %v\n", err)
		return err
	}

	nsSpec := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	_, err = clientset.CoreV1().
		Namespaces().
		Create(context.Background(), nsSpec, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("error creating kubernetes namespace: %v\n", err)
		return err
	}
	return nil
}

// createKubesDeployment creates a deployment in the specified namespace.
func createKubesDeployment(namespace string, image string, tag string) error {
	fmt.Println("Create Kubernetes deployment")

	clientset, err := getKubesClient()
	if err != nil {
		log.Printf("error getting kubernetes client: %v\n", err)
		return err

	}

	deploymentSpec := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: func(val int32) *int32 { return &val }(3),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:latest",
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
	_, err = clientset.AppsV1().
		Deployments(namespace).
		Create(context.Background(), deploymentSpec, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("error creating deployment: %v\n", err)
		return err
	}
	return nil
}

// getKubesClient returns a Kubernetes clientset.
func getKubesClient() (*kubernetes.Clientset, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting user home dir: %v", err)
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, fmt.Errorf("error getting Kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("error getting Kubernetes clientset: %v", err)
	}

	return clientset, nil
}
