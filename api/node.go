package api

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/rancher/go-rancher/api"
)

func (s *Server) NodeList(w http.ResponseWriter, req *http.Request) error {
	apiContext := api.GetApiContext(req)

	nodes, err := s.c.ClusterNodes()
	if err != nil || nodes == nil {
		return errors.Wrap(err, "fail to read nodes")
	}
	apiContext.Write(toHostCollection(*nodes))
	return nil
}
