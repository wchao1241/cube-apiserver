package backend

import (
	"github.com/rancher/rancher-cube-apiserver/util"
)

func (c *ClientGenerator) ClusterNodes() (*map[string]string, error) {
	nodeMap := make(map[string]string)
	nodeList, err := c.clientset.CoreV1().Nodes().List(util.ListEverything)

	for _, node := range nodeList.Items {
		if node.UID != "" {
			UID := string(node.UID)
			nodeMap[UID] = node.Name
		}
	}

	return &nodeMap, err
}
