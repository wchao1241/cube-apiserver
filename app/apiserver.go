package app

import (
	"fmt"

	"github.com/rancher/rancher-cube-apiserver/util"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	FlagAPIServerImage = "apiserver-image"
)

func APIServerCmd() cli.Command {
	return cli.Command{
		Name: "daemon",
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

	util.RegisterShutdownChannel(done)
	<-done

	return nil
}
