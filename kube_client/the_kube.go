// Common functionality relating to connecting to our Kubernetes cluster

package kube_client

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func Client() (*kubernetes.Clientset, error) {
	kubeconfig, err := RestConfig()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error getting kubeconfig")
		return nil, err
	}

	return kubernetes.NewForConfig(kubeconfig)
}

func RestConfig() (*rest.Config, error) {
	// Connect to the Kubernetes cluster
	// and return a client
	home, err := os.UserHomeDir()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error(err)
		return nil, err
	}
	kubeConfigPath := filepath.Join(home, ".config", "podinate", ".kube", "config")

	kubeconfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":            err,
			"kube config path": kubeConfigPath,
		}).Error(err)
		return nil, err
	}

	return kubeconfig, nil
}
