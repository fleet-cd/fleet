package kubernetes

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tgs266/fleet/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type ClientSet struct {
	clientSet *kubernetes.Clientset
}

var clientSet *ClientSet

func Connect(cfg *config.Config) {
	fmt.Println("Get Kubernetes pods")

	kubeConfigPath := filepath.Join(cfg.Kubernetes.ConfigPath)
	fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Printf("error getting Kubernetes config: %v\n", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		fmt.Printf("error getting Kubernetes clientset: %v\n", err)
		os.Exit(1)
	}

	clientSet = &ClientSet{
		clientSet: clientset,
	}
}
