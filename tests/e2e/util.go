package e2e

import (
	"context"
	"errors"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

const (
	podRunningTimeout = 30 * time.Minute
	testWorkspaceDir  = "/e2e"
	inputDir          = "inputs" // directory to use in test. required.
	//ProvisonerImage           = "172.22.4.104:5000/hypersds-provisioner:test"
	podName            = "hypersds-provisioner"
	podVolumeName      = "config-volume"
	podVolumeMountPath = "/manifest"
	//ProvisonerNamespace       = "default"
	registryName = "regcred"
)

func runProvisionerContainer(client *kubernetes.Clientset, provisionerImage, provisionerNamespace, nodeName string) error {
	configPath := testWorkspaceDir + "/" + inputDir

	pod := newProvisonerPod(configPath, provisionerImage, provisionerNamespace, nodeName)
	if _, err := client.CoreV1().Pods(provisionerNamespace).Create(context.TODO(), pod, metav1.CreateOptions{}); err != nil {
		return err
	}
	return wait.PollImmediate(time.Second, podRunningTimeout, isPodCompleted(client, provisionerNamespace))
}

func isPodCompleted(client *kubernetes.Clientset, provisionerNamespace string) wait.ConditionFunc {
	return func() (bool, error) {
		pod, err := client.CoreV1().Pods(provisionerNamespace).Get(context.TODO(), podName, metav1.GetOptions{})
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

func newProvisonerPod(volumePath, provisionerImage, provisionerNamespace, nodeName string) *corev1.Pod {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: provisionerNamespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            podName,
					Image:           provisionerImage,
					ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
					Args:            []string{},
					VolumeMounts: []corev1.VolumeMount{
						{Name: podVolumeName, MountPath: podVolumeMountPath},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: podVolumeName,
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: volumePath,
						},
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
			ImagePullSecrets: []corev1.LocalObjectReference{
				{Name: registryName},
			},
		},
	}
	if nodeName != "" {
		pod.Spec.NodeName = nodeName
	}
	return pod
}
