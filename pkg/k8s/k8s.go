package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func BuildClient(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
