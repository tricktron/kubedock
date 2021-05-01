package kubernetes

import (
	"fmt"
	"io"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"

	"github.com/joyrex2001/kubedock/internal/container"
)

// ExecContainer will execute given exec object in kubernetes.
func (in *instance) ExecContainer(tainr *container.Container, exec *container.Exec, out io.Writer) error {
	pods, err := in.getPods(tainr)
	if err != nil {
		return err
	}
	if len(pods) == 0 {
		return fmt.Errorf("no matching pod found")
	}

	req := in.cli.CoreV1().RESTClient().Post().Resource("pods").
		Name(pods[0].Name).
		Namespace(pods[0].Namespace).
		SubResource("exec")
	req.VersionedParams(&corev1.PodExecOptions{
		Command: exec.Cmd,
		Stdin:   false,
		Stdout:  exec.Stdout,
		Stderr:  exec.Stderr,
		TTY:     false,
	}, scheme.ParameterCodec)
	ex, err := remotecommand.NewSPDYExecutor(in.cfg, "POST", req.URL())
	if err != nil {
		return err
	}

	opts := remotecommand.StreamOptions{}
	if exec.Stdout {
		opts.Stdout = out
	}
	if exec.Stderr {
		opts.Stderr = out
	}
	return ex.Stream(opts)
}