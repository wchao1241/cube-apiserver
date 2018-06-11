package api

import (
	"github.com/cnrancher/cube-apiserver/backend"
	"github.com/cnrancher/cube-apiserver/k8s/pkg/apis/cube/v1alpha1"

	"github.com/rancher/go-rancher/client"
	"k8s.io/api/core/v1"
	"github.com/cnrancher/cube-apiserver/util"
)

var KubeConfigLocation string

var sType *schemaType

type schemaType struct {
	Configmap      string
	Dashboard      string
	Longhorn       string
	BaseInfo       string
	RancherVM      string
	Infrastructure string
}

func GetSchemaType() *schemaType {
	if sType == nil {
		sType = &schemaType{
			Configmap:      "configmap",
			Dashboard:      "dashboard",
			Longhorn:       "longhorn",
			BaseInfo:       "baseinfo",
			RancherVM:      "ranchervm",
			Infrastructure: "infrastructure",
		}
	}

	return sType
}

type Server struct {
	// TODO: HTTP reverse proxy should be here
	c *backend.ClientGenerator
}

type Node struct {
	client.Resource

	Node *v1.Node `json:"node"`
}

type ConfigMap struct {
	client.Resource

	ConfigMap *v1.ConfigMap `json:"configmap"`
}

type BaseInfo struct {
	client.Resource

	BaseInfo []map[string]string `json:"baseinfo"`
}

type Infrastructure struct {
	client.Resource
	Infra *v1alpha1.Infrastructure `json:"infrastructure"`
	Host  string                   `json:"host"`
}

type Cluster struct {
	client.Resource

	Name              string                   `json:"name"`
	Resources         *backend.ClusterResource `json:"resources"`
	ComponentStatuses *v1.ComponentStatusList  `json:"componentStatuses"`
}

type Token struct {
	Name      string `json:"name"`
	LoginName string `json:"loginName"`
	Password  string `json:"password"`
	Token     string `json:"token"`
}

func NewServer(c *backend.ClientGenerator, kubeConfig string) *Server {
	KubeConfigLocation = kubeConfig
	s := &Server{
		c: c,
	}
	return s
}

type Result struct {
	client.Resource

	Message string `json:"message"`
}

func NewSchema() *client.Schemas {
	schemas := &client.Schemas{}

	schemas.AddType("apiVersion", client.Resource{})
	schemas.AddType("schema", client.Schema{})
	schemas.AddType("error", client.ServerApiError{})

	nodeSchema(schemas.AddType("node", Node{}))
	clusterSchema(schemas.AddType("cluster", Cluster{}))

	baseInfoschema(schemas.AddType(GetSchemaType().BaseInfo, BaseInfo{}))
	infraSchema(schemas.AddType(GetSchemaType().Infrastructure, Infrastructure{}))
	return schemas
}

func clusterSchema(cluster *client.Schema) {
	cluster.CollectionMethods = []string{"GET"}
	cluster.ResourceMethods = []string{"GET"}

	cluster.ResourceFields["resources"] = client.Field{
		Type:     "struct",
		Nullable: true,
	}
	cluster.ResourceFields["componentStatuses"] = client.Field{
		Type:     "struct",
		Nullable: true,
	}
}

func nodeSchema(node *client.Schema) {
	node.CollectionMethods = []string{"GET"}
	node.ResourceMethods = []string{"GET"}

	node.ResourceFields["node"] = client.Field{
		Type:     "struct",
		Nullable: true,
	}
}

func baseInfoschema(cm *client.Schema) {
	cm.CollectionMethods = []string{"GET"}
	cm.ResourceMethods = []string{"GET"}

	cm.ResourceFields["resources"] = client.Field{
		Type:     "struct",
		Nullable: true,
	}
	cm.ResourceFields["componentStatuses"] = client.Field{
		Type:     "struct",
		Nullable: true,
	}

}

func infraSchema(infra *client.Schema) {
	infra.CollectionMethods = []string{"GET", "POST"}
	infra.ResourceMethods = []string{"GET", "DELETE"}

	infra.ResourceFields["infrastructure"] = client.Field{
		Type:     "struct",
		Nullable: true,
	}

}

func toInfrastructureCollection(list *v1alpha1.InfrastructureList) *client.GenericCollection {
	data := []interface{}{}
	for _, item := range list.Items {
		data = append(data, toInfrastructureResource(&item, nil, ""))
	}
	return &client.GenericCollection{Data: data, Collection: client.Collection{ResourceType: GetSchemaType().Infrastructure}}
}

func toInfrastructureResource(infra *v1alpha1.Infrastructure, service *v1.Service, ip string) *Infrastructure {
	name := infra.GetName()
	host := ip

	if service != nil {
		port := service.Spec.Ports[0].NodePort
		host = host + ":" + util.Int32ToString(port)
	}
	if name == "" {
		return nil
	}
	return &Infrastructure{
		Resource: client.Resource{
			Id:      string(name),
			Type:    "infrastructure",
			Actions: map[string]string{},
		},
		Host:  host,
		Infra: infra,
	}
}

func toDeleteResource(resType string) *Result {
	msg := "delete success"
	return &Result{
		Resource: client.Resource{
			Id:      string("delete"),
			Type:    resType,
			Actions: map[string]string{},
		},
		Message: msg,
	}
}

func toBaseInfo(baseInfo []map[string]string) *BaseInfo {
	if baseInfo == nil {
		return nil
	}
	return &BaseInfo{
		Resource: client.Resource{
			Id:      string("baseInfo"),
			Type:    GetSchemaType().BaseInfo,
			Actions: map[string]string{},
		},
		BaseInfo: baseInfo,
	}
}

func toNodeResource(node *v1.Node) *Node {
	name := node.GetName()
	if name == "" {
		return nil
	}
	return &Node{
		Resource: client.Resource{
			Id:      string(name),
			Type:    "node",
			Actions: map[string]string{},
		},
		Node: node,
	}
}

func toNodeCollection(nodeList *v1.NodeList) *client.GenericCollection {
	data := []interface{}{}
	for _, node := range nodeList.Items {
		data = append(data, toNodeResource(&node))
	}
	return &client.GenericCollection{Data: data, Collection: client.Collection{ResourceType: "node"}}
}

func toClusterResource(clusterResource *backend.ClusterResource, componentStatuses *v1.ComponentStatusList) *Cluster {
	return &Cluster{
		Resource: client.Resource{
			Id:      "cube",
			Type:    "cluster",
			Actions: map[string]string{},
		},
		Name:              "cube",
		Resources:         clusterResource,
		ComponentStatuses: componentStatuses,
	}
}

func toClusterCollection(clusterResource *backend.ClusterResource, componentStatuses *v1.ComponentStatusList) *client.GenericCollection {
	data := []interface{}{}
	data = append(data, toClusterResource(clusterResource, componentStatuses))
	return &client.GenericCollection{Data: data, Collection: client.Collection{ResourceType: "cluster"}}
}
