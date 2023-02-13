package provision

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ensureClientPod ensures a pod and returns the name, with an error when failing.
func ensureClientPod(client *kubernetes.Clientset, name, namespace, image string, command []string, gateway string) error {
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				KubernetesLabelAppName:    name,
				KubernetesLabelSetGateway: gateway,
			},
			Annotations: map[string]string{
				KubernetesAnnotationSetGateway: gateway,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    ContainerName,
					Image:   image,
					Command: command,
				},
			},
		},
	}

	_, err := client.CoreV1().Pods(namespace).Create(
		context.Background(),
		&pod,
		metav1.CreateOptions{},
	)
	if err != nil {
		return err
	}

	return nil
}
