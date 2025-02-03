package k8s

import (
	"log"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type KubernetesClientType *kubernetes.Clientset

var kubernetesClient *kubernetes.Clientset

func NewClient() (KubernetesClientType, error) {
	config, err := clientcmd.BuildConfigFromFlags(
		"",
		filepath.Join(homedir.HomeDir(), ".kube", "config"),
	)
	if err != nil {
		return nil, err
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	kubernetesClient=k8sClient;
	log.Println("client created kubernetes")
	return k8sClient, nil
}
