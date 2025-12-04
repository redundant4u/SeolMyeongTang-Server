package k8s

import (
	"fmt"
	"seolmyeong-tang-server/internal/pkg/logger"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	Clientset *kubernetes.Clientset
}

func NewClient(kubeconfig string) (*Client, error) {
	var config *rest.Config
	var err error

	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("build kubeconfig from %s: %w", kubeconfig, err)
			// panic("Faile to build kubeconfig")
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("create in-cluster config: %w", err)
			// panic("Failed to create in-cluster config")
		}
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("create kubernetes clientset: %w", err)
		// panic("Failed to create clientset")
	}

	logger.Info("Kubernetes client initialized")
	return &Client{Clientset: clientSet}, nil
}
