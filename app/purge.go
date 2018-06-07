package app

import (
	"time"

	"github.com/cnrancher/cube-apiserver/backend"
	tokenclientset "github.com/cnrancher/cube-apiserver/k8s/pkg/client/clientset/versioned"
	tokenlisters "github.com/cnrancher/cube-apiserver/k8s/pkg/client/listers/cube/v1alpha1"
	"github.com/cnrancher/cube-apiserver/util"

	"github.com/Sirupsen/logrus"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
)

const intervalSeconds int64 = 3600

type purger struct {
	tokenLister tokenlisters.TokenLister
	clientSet   tokenclientset.Clientset
}

func StartPurgeDaemon(clientGenerator *backend.ClientGenerator, done chan struct{}) {
	// obtain references to shared index informers for the Token
	// types.
	tokenInformer := clientGenerator.CubeInformerFactory.Cube().V1alpha1().Tokens()

	p := &purger{
		tokenLister: tokenInformer.Lister(),
		clientSet:   clientGenerator.Infraclientset,
	}
	go wait.JitterUntil(p.purge, time.Duration(intervalSeconds)*time.Second, .1, true, done)
}

func (p *purger) purge() {
	tokens, err := p.tokenLister.Tokens("").List(labels.Everything())
	if err != nil {
		logrus.Errorf("RancherCUBE: error listing tokens during purge: %v", err)
	}

	var count int
	for _, token := range tokens {
		if util.IsExpired(*token) {
			err = p.clientSet.CubeV1alpha1().Tokens(token.ObjectMeta.Namespace).Delete(token.ObjectMeta.Name, &metaV1.DeleteOptions{})
			if err != nil && !k8serrors.IsNotFound(err) {
				logrus.Errorf("RancherCUBE: error while deleting expired token %v: %v", err, token.ObjectMeta.Name)
				continue
			}
			count++
		}
	}
	if count > 0 {
		logrus.Infof("RancherCUBE: purged %v expired tokens", count)
	}
}
