package backend

import (
	"github.com/Sirupsen/logrus"
	apics "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ClientGenerator struct {
	clientset    *kubernetes.Clientset
	apiclientset *apics.Clientset
}

func NewClientGenerator(kubeConfig string) *ClientGenerator {
	var config *rest.Config
	var err error

	if kubeConfig == "" {
		config, err = rest.InClusterConfig()
		if err != nil {
			logrus.Fatalf("RancherCUBE: generate config failed: %v", err)
		}
	}

	config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		logrus.Fatalf("RancherCUBE: generate config failed: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Fatalf("RancherCUBE: generate clientset failed: %v", err)
	}

	apiclientset, err := apics.NewForConfig(config)
	if err != nil {
		logrus.Fatalf("RancherCUBE: generate extensions clientset failed: %v", err)
	}

	return &ClientGenerator{
		clientset:    clientset,
		apiclientset: apiclientset,
	}
}
