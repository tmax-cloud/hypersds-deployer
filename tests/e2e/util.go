package e2e

import (
	"context"
	"errors"
	"time"

	corev1 "k8s.io/api/core/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

const (
	podRunningTimeout  = 30 * time.Minute
	podRemovingTimeout = 30 * time.Minute
	podName            = "hypersds-provisioner"
	podVolumeName      = "config-volume"
	podVolumeMountPath = "/manifest"
)

func runProvisionerPod(client *kubernetes.Clientset, provisionerNamespace, provisionerImage, volumeHostPath, registryName, nodeName string) error {
	pod := newProvisonerPod(provisionerNamespace, provisionerImage, volumeHostPath, registryName, nodeName)
	if _, err := client.CoreV1().Pods(provisionerNamespace).Create(context.TODO(), pod, metav1.CreateOptions{}); err != nil {
		return err
	}
	return wait.PollImmediate(time.Second, podRunningTimeout, isPodCompleted(client, provisionerNamespace))
}

func removeProvisionerPod(client *kubernetes.Clientset, provisionerNamespace string) error {
	podDeletePolicy := metav1.DeletePropagationForeground
	err := client.CoreV1().Pods(provisionerNamespace).Delete(context.TODO(), podName, metav1.DeleteOptions{
		PropagationPolicy: &podDeletePolicy,
	})
	if err != nil {
		if kubeerrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	return wait.PollImmediate(time.Second, podRemovingTimeout, isPodDeleted(client, provisionerNamespace))
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

func newProvisonerPod(provisionerNamespace, provisionerImage, volumeHostPath, registryName, nodeName string) *corev1.Pod {
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
							Path: volumeHostPath,
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

func isPodDeleted(client *kubernetes.Clientset, provisionerNamespace string) wait.ConditionFunc {
	return func() (bool, error) {
		_, err := client.CoreV1().Pods(provisionerNamespace).Get(context.TODO(), podName, metav1.GetOptions{})
		if err != nil {
			if kubeerrors.IsNotFound(err) {
				return true, nil
			} else {
				return false, err
			}
		}
		return false, nil
	}
}
