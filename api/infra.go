package api

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/rancher/go-rancher/api"
	"github.com/cnrancher/cube-apiserver/controller"
	"github.com/cnrancher/cube-apiserver/backend"
	"github.com/gorilla/mux"
)

func (s *Server) InfrastructureList(rw http.ResponseWriter, req *http.Request) error {
	apiContext := api.GetApiContext(req)

	list, err := s.c.InfrastructureList()
	if err != nil || list == nil {
		return errors.Wrap(err, "failed to read dashboard")
	}
	apiContext.Write(toInfrastructureCollection(list))
	return nil
}

func (s *Server) InfrastructureGet(rw http.ResponseWriter, req *http.Request) error {
	apiContext := api.GetApiContext(req)
	kind := mux.Vars(req)["id"]

	db, err := s.c.InfrastructureGet(kind)
	if err != nil || db == nil {
		return errors.Wrap(err, "failed to read dashboard")
	}

	err = s.c.ServiceGet(controller.DashboardNamespace, "kubernetes-dashboard")
	if err != nil {
		return errors.Wrap(err, "failed to get dashboard service")
	}

	ip, err := s.c.NodeIPGet()
	if err != nil {
		return errors.Wrap(err, "failed to get node ip")
	}
	apiContext.Write(toInfrastructureResource(db, backend.Service, ip))
	return nil
}

func (s *Server) InfrastructureCreate(rw http.ResponseWriter, req *http.Request) error {
	//var volume Volume
	//apiContext := api.GetApiContext(req)
	//
	//if err := apiContext.Read(&volume); err != nil {
	//	return err
	//}
	//
	//size, err := util.ConvertSize(volume.Size)
	//if err != nil {
	//	return fmt.Errorf("fail to parse size %v", err)
	//}
	//v, err := s.m.Create(volume.Name, &types.VolumeSpec{
	//	Size:                size,
	//	FromBackup:          volume.FromBackup,
	//	NumberOfReplicas:    volume.NumberOfReplicas,
	//	StaleReplicaTimeout: volume.StaleReplicaTimeout,
	//})
	//if err != nil {
	//	return errors.Wrap(err, "unable to create volume")
	//}




	apiContext := api.GetApiContext(req)
	kind := mux.Vars(req)["id"]

	db, err := s.c.InfrastructureDeploy(kind)
	if err != nil {
		return errors.Wrap(err, "failed to create dashboard")
	}

	err = s.c.ServiceGet(controller.DashboardNamespace, "kubernetes-dashboard")
	if err != nil {
		return errors.Wrap(err, "failed to get dashboard service")
	}

	ip, err := s.c.NodeIPGet()
	if err != nil {
		return errors.Wrap(err, "failed to get node ip")
	}
	apiContext.Write(toInfrastructureResource(db, backend.Service, ip))
	return nil
}

func (s *Server) InfrastructureDelete(rw http.ResponseWriter, req *http.Request) error {
	apiContext := api.GetApiContext(req)
	kind := mux.Vars(req)["id"]

	err := s.c.InfrastructureDelete(kind)
	if err != nil {
		return errors.Wrap(err, "failed to delete dashboard")
	}
	apiContext.Write(toDeleteResource( /*GetSchemaType().Dashboard*/ kind))
	return nil
}
