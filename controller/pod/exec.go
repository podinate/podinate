package pod

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"strings"

	apierror "github.com/johncave/podinate/controller/apierror"
	"github.com/johncave/podinate/controller/config"
	lh "github.com/johncave/podinate/controller/loghandler"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// Exec executes a command in the pod
func (p *Pod) Exec(ctx context.Context, command []string) (string, *apierror.ApiError) {

	lh.Debug(ctx, "In exec func", "pod", p, "command", strings.Join(command, " "))

	// cmd := []string{
	// 	"sh",
	// 	"-c",
	// 	strings.Join(command, " "),
	// }
	// TODO: Update "-0" to be able to address multiple containers
	req := config.Client.CoreV1().RESTClient().Post().Resource("pods").Name(p.ID + "-0").Namespace(p.Project.GetNamespaceName()).SubResource("exec")
	req.VersionedParams(&v1.PodExecOptions{
		Command: command,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     false,
	}, scheme.ParameterCodec)

	rconfig, err := restclient.InClusterConfig()
	if err != nil {
		return "", apierror.NewWithError(http.StatusInternalServerError, "error getting in-cluster config", err)
	}
	exec, err := remotecommand.NewSPDYExecutor(rconfig, "POST", req.URL())
	if err != nil {
		return "", apierror.NewWithError(http.StatusInternalServerError, "error creating executor to run command", err)
	}

	stdout := new(bytes.Buffer)
	lh.Info(ctx, "Execute command on pod", "pod", p, "command", strings.Join(command, " "))
	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: stdout,
		Stderr: stdout,
	})
	if err != nil {
		// At this point the user messed up, not us ¯\_(ツ)_/¯
		return err.Error(), nil
	}

	lh.Debug(ctx, "Executed command", "stdout", stdout.String())

	return stdout.String(), nil
}
