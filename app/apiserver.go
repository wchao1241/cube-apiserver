package app

import (
	"fmt"
	"net/http"

	"github.com/rancher/rancher-cube-apiserver/api"
	"github.com/rancher/rancher-cube-apiserver/util"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/rancher-cube-apiserver/backend"
	"github.com/urfave/cli"
)

const (
	ListenAddress  = "listen-addr"
	ConfigLocation = "kube-config"
)

func APIServerCmd() cli.Command {
	return cli.Command{
		Name: "serve",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  ListenAddress,
				Usage: "Specify apiServer listen address",
			},
			cli.StringFlag{
				Name:  ConfigLocation,
				Usage: "Specify apiServer kubernetes config location",
			},
		},
		Action: func(c *cli.Context) error {
			if err := startAPIServer(c); err != nil {
				logrus.Errorf("RancherCUBE: error starting apiServer: %v", err)
				return err
			}
			return nil
		},
	}
}

func startAPIServer(c *cli.Context) error {
	apiServerListenAddr := c.String(ListenAddress)
	apiServerKubeConfigLocation := c.String(ConfigLocation)
	if "" == apiServerListenAddr {
		return fmt.Errorf("RancherCUBE: require %v", ListenAddress)
	}
	if "" == apiServerKubeConfigLocation {
		return fmt.Errorf("RancherCUBE: require %v", ConfigLocation)
	}

	clientGenerator := backend.NewClientGenerator(apiServerKubeConfigLocation)

	// generate & deploy infrastructure customer resources
	clientGenerator.InfrastructureCRDDeploy()

	done := make(chan struct{})

	server := api.NewServer(clientGenerator)
	router := http.Handler(api.NewRouter(server))

	logrus.Infof("RancherCUBE: listening on %s", apiServerListenAddr)

	go http.ListenAndServe(apiServerListenAddr, router)

	util.RegisterShutdownChannel(done)
	<-done

	return nil
}
