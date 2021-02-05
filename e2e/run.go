package e2e

import (
	"context"
	"errors"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

const (
	ProvisonerTimeout         = 30 * time.Minute
	InputDir                  = "example" //temp directory
	ProvisonerImage           = "hypersds-provisoner:canary"
	ProvisonerName            = "hypersds-provisoner"
	ProvisonerVolumeName      = "config-volume"
	ProvisonerVolumeMountPath = "/manifest"
	ProvisonerNamespace       = "default"
)

func runProvisionerContainer(client *kubernetes.Clientset, nodeName string) error {
	currentPath, err := os.Getwd()
	if err != nil {
		return err
	}
	configPath := currentPath + "/" + InputDir
	pod := newProvisonerPod(configPath, nodeName)
	if _, err := client.CoreV1().Pods(ProvisonerNamespace).Create(context.TODO(), pod, metav1.CreateOptions{}); err != nil {
		return err
	}
	return wait.PollImmediate(time.Second, ProvisonerTimeout, isPodCompleted(client))
}

func isPodCompleted(client *kubernetes.Clientset) wait.ConditionFunc {
	return func() (bool, error) {
		pod, err := client.CoreV1().Pods(ProvisonerNamespace).Get(context.TODO(), ProvisonerName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		switch pod.Status.Phase {
		case corev1.PodSucceeded:
			return true, nil
		case corev1.PodFailed:
			return true, errors.New("Pod failed")
		}
		return false, nil
	}
}

func newProvisonerPod(volumePath, nodeName string) *corev1.Pod {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ProvisonerName,
			Namespace: ProvisonerNamespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            ProvisonerName,
					Image:           ProvisonerImage,
					ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
					Args:            []string{},
					VolumeMounts: []corev1.VolumeMount{
						{Name: ProvisonerVolumeName, MountPath: ProvisonerVolumeMountPath},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: ProvisonerVolumeName,
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: volumePath,
						},
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
	if nodeName != "" {
		pod.Spec.NodeName = nodeName
	}
	return pod
}
