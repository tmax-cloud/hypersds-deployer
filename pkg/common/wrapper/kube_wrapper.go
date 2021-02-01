package wrapper

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubeInterface interface {
	InClusterConfig() (*rest.Config, error)
	NewForConfig(c *rest.Config) (kubernetes.Interface, error)
}

type kubeStruct struct {
}

func (r *kubeStruct) InClusterConfig() (*rest.Config, error) {
	return rest.InClusterConfig()
}
func (r *kubeStruct) NewForConfig(c *rest.Config) (kubernetes.Interface, error) {
	return kubernetes.NewForConfig(c)
}

var KubeWrapper KubeInterface

func init() {
	KubeWrapper = &kubeStruct{}
}
