package api

import (
	"github.com/rancher/go-rancher/client"
)

type Server struct {
	// TODO: HTTP reverse proxy should be here
}

func NewServer() *Server {
	s := &Server{}
	return s
}

func NewSchema() *client.Schemas {
	schemas := &client.Schemas{}

	schemas.AddType("apiVersion", client.Resource{})
	schemas.AddType("schema", client.Schema{})
	schemas.AddType("error", client.ServerApiError{})

	return schemas
}
