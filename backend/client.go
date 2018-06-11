package backend

import (
	"time"

	infracs "github.com/cnrancher/cube-apiserver/k8s/pkg/client/clientset/versioned"
	cubeinformers "github.com/cnrancher/cube-apiserver/k8s/pkg/client/informers/externalversions"

	"github.com/Sirupsen/logrus"
	apics "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	listerscorev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

var (
	clientGenerator *ClientGenerator
)

type ClientGenerator struct {
	Clientset           kubernetes.Clientset
	Apiclientset        apics.Clientset
	Infraclientset      infracs.Clientset
	InformerFactory     informers.SharedInformerFactory
	CubeInformerFactory cubeinformers.SharedInformerFactory
	serviceLister       listerscorev1.ServiceLister
	serviceSynced       cache.InformerSynced
}

func NewClientGenerator(kubeConfig string) *ClientGenerator {
	if clientGenerator == nil {
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

		infraclientset, err := infracs.NewForConfig(config)
		if err != nil {
			logrus.Fatalf("RancherCUBE: generate infra clientset failed: %v", err)
		}

		informerFactory := informers.NewSharedInformerFactory(clientset, time.Second*30)
		infraInformerFactory := cubeinformers.NewSharedInformerFactory(infraclientset, time.Second*30)
		serviceInformer := informerFactory.Core().V1().Services()

		clientGenerator = &ClientGenerator{
			Clientset:           *clientset,
			Apiclientset:        *apiclientset,
			Infraclientset:      *infraclientset,
			InformerFactory:     informerFactory,
			CubeInformerFactory: infraInformerFactory,
			serviceLister:       serviceInformer.Lister(),
			serviceSynced:       serviceInformer.Informer().HasSynced,
		}
	}

	return clientGenerator
}
