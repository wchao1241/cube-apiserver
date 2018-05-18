package util

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Sirupsen/logrus"
)

func RegisterShutdownChannel(done chan struct{}) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		logrus.Infof("RancherCUBE: receive %v to exit", sig)
		close(done)
	}()
}
