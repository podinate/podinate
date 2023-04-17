package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func callKubes() {
	fmt.Println("Get Kubernetes pods")

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting user home dir: %v\n", err)
		os.Exit(1)
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Printf("error getting Kubernetes config: %v\n", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		fmt.Printf("error getting Kubernetes clientset: %v\n", err)
		os.Exit(1)
	}

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
