package api

import (
	"github.com/rancher/go-rancher/client"
	"github.com/rancher/rancher-cube-apiserver/backend"
)

type Server struct {
	// TODO: HTTP reverse proxy should be here
	c *backend.ClientGenerator
}

type Host struct {
	client.Resource

	UUID string
	Name string
}

func NewServer(c *backend.ClientGenerator) *Server {
	s := &Server{
		c: c,
	}
	return s
}

func NewSchema() *client.Schemas {
	schemas := &client.Schemas{}

	schemas.AddType("apiVersion", client.Resource{})
	schemas.AddType("schema", client.Schema{})
	schemas.AddType("error", client.ServerApiError{})

	hostSchema(schemas.AddType("host", Host{}))

	return schemas
}

func hostSchema(host *client.Schema) {
	host.CollectionMethods = []string{"GET"}
	host.ResourceMethods = []string{"GET"}
}

func toHostResource(uuid, name string) *Host {
	return &Host{
		Resource: client.Resource{
			Id:      name,
			Type:    "host",
			Actions: map[string]string{},
		},
		Name: name,
		UUID: uuid,
	}
}

func toHostCollection(nodeMap map[string]string) *client.GenericCollection {
	data := []interface{}{}
	for uuid, name := range nodeMap {
		data = append(data, toHostResource(uuid, name))
	}
	return &client.GenericCollection{Data: data}
}
