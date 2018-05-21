package app

import (
	"fmt"
	"net/http"

	"github.com/rancher/rancher-cube-apiserver/api"
	"github.com/rancher/rancher-cube-apiserver/util"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	FlagAPIServerImage = "apiserver-image"
)

func APIServerCmd() cli.Command {
	return cli.Command{
		Name: "serve",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  FlagAPIServerImage,
				Usage: "Specify apiServer image",
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
	apiServerImage := c.String(FlagAPIServerImage)
	if "" == apiServerImage {
		return fmt.Errorf("RancherCUBE: require %v", FlagAPIServerImage)
	}

	done := make(chan struct{})

	// TODO: define router...

	server := api.NewServer()
	router := http.Handler(api.NewRouter(server))

	logrus.Infof("RancherCUBE: listening on %s", "127.0.0.1:9500")

	go http.ListenAndServe(":9500", router)

	util.RegisterShutdownChannel(done)
	<-done

	return nil
}
