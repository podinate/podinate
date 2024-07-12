// Common functionality relating to connecting to our Kubernetes cluster

package kube_client

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var client *kubernetes.Clientset
var RestConfig *rest.Config

func init() {
	// rc, err := GetRestConfig()
	// if err != nil {
	// 	logrus.WithFields(logrus.Fields{
	// 		"error": err,
	// 	}).Fatal("Error getting kubeconfig")
	// }
	// RestConfig = rc

	// c, err := kubernetes.NewForConfig(RestConfig)
	// if err != nil {
	// 	logrus.WithFields(logrus.Fields{
	// 		"error": err,
	// 	}).Fatal("Error creating kube client")
	// }
	// client = c

}

func Client() (*kubernetes.Clientset, error) {

	// Use the cached client if available
	if client != nil {
		return client, nil
	}

	kubeconfig, err := GetRestConfig()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error getting kubeconfig")
		return nil, err
	}

	c, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error creating kube client")
		return nil, err
	}

	client = c

	return client, nil
}

func GetRestConfig() (*rest.Config, error) {
	// Connect to the Kubernetes cluster
	// and return a client
	home, err := os.UserHomeDir()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error(err)
		return nil, err
	}
	kubeConfigPath := filepath.Join(home, ".kube", "config")

	if viper.GetString("kubeconfig") != "" {
		logrus.WithFields(logrus.Fields{
			"kubeconfig": viper.GetString("kubeconfig"),
		}).Debug("Using kubeconfig from flag")
		kubeConfigPath = viper.GetString("kubeconfig")
	}

	configLoadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath}
	configOverrides := &clientcmd.ConfigOverrides{}

	logrus.WithFields(logrus.Fields{
		"kubeconfig": kubeConfigPath,
		"context":    viper.GetString("context"),
	}).Debug("Using kubeconfig")

	if viper.GetString("context") != "" {
		configOverrides = &clientcmd.ConfigOverrides{CurrentContext: viper.GetString("context")}
	}

	//kubeconfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides)
	kubeconfig, err := config.ClientConfig()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":            err,
			"kube config path": kubeConfigPath,
		}).Error(err)
		return nil, err
	}

	return kubeconfig, nil
}
