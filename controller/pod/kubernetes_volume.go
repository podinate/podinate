package pod

import (
	"context"

	lh "github.com/johncave/podinate/controller/loghandler"
	v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getStorageClass(name string) (*v1.StorageClass, error) {
	// Check if the storage class exists

	clientset, err := getKubesClient()
	if err != nil {
		lh.Log.Errorw("error getting kubernetes client to check for StorageClass existence", "error", err.Error())
		return nil, err
	}

	return clientset.StorageV1().StorageClasses().Get(context.Background(), name, metav1.GetOptions{})

}
