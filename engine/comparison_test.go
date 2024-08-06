package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGenerateResourceChange(t *testing.T) {
	cm := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cm",
			Namespace: "default",
		},
		Data: map[string]string{
			"key": "value",
		},
	}
	cm.Kind = "ConfigMap"
	cm.APIVersion = "v1"

	rc, err := getRestHelperForObject(context.Background(), &cm)
	if err != nil {
		t.Errorf("Error getting resource change: %v", err)
	}

	if rc == nil {
		t.Errorf("Resource change is nil")
	}

	out, err := json.Marshal(rc)
	if err != nil {
		t.Errorf("Error marshalling resource change: %v", err)
	}
	fmt.Println(string(out))
}
