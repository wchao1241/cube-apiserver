package backend

import (
	"k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/cnrancher/cube-apiserver/util"
	"github.com/cnrancher/cube-apiserver/controller"
)

var (
//dashboardName      = "dashboardName"
//dashboardNamespace = "dashboardNamespace"
//dashboardIcon      = "dashboardIcon"
//dashboardDesc      = "dashboardDesc"
//
//longhornName      = "longhornName"
//longhornNamespace = "longhornNamespace"
//longhornIcon      = "longhornIcon"
//longhornDesc      = "longhornDesc"
//
//rancherVMName      = "rancherVMName"
//rancherVMNamespace = "rancherVMNamespace"
//rancherVMIcon      = "rancherVMIcon"
//rancherVMDesc      = "rancherVMDesc"
)

var (
	ConfigMapName   = "cube-rancher"
	Infrastructures *infrastructures
)

type infrastructure struct {
	name      string
	namespace string
	icon      string
	desc      string
}

type infrastructures struct {
	dashboard *infrastructure
	longhorn  *infrastructure
	rancherVM *infrastructure
}

func initInfrastructures() {
	dashboard := &infrastructure{
		name:      "dashboardName",
		namespace: "dashboardNamespace",
		icon:      "dashboardIcon",
		desc:      "dashboardDesc",
	}

	longhorn := &infrastructure{
		name:      "longhornName",
		namespace: "longhornNamespace",
		icon:      "longhornIcon",
		desc:      "longhornDesc",
	}

	rancherVM := &infrastructure{
		name:      "rancherVMName",
		namespace: "rancherVMNamespace",
		icon:      "rancherVMIcon",
		desc:      "rancherVMDesc",
	}

	Infrastructures = &infrastructures{
		dashboard: dashboard,
		longhorn:  longhorn,
		rancherVM: rancherVM,
	}
}

func (c *ClientGenerator) ConfigMapDeploy() (*v1.ConfigMap, error) {
	initInfrastructures()
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ConfigMapName,
			Namespace: controller.DashboardNamespace,
			Labels: map[string]string{
				"app": ConfigMapName,
			}},

		Data: map[string]string{

			Infrastructures.dashboard.name:      controller.DashboardName,
			Infrastructures.dashboard.namespace: controller.DashboardNamespace,
			Infrastructures.dashboard.icon:      "https://avatars3.githubusercontent.com/u/13629408?s=200&v=4",
			Infrastructures.dashboard.desc:      controller.DashboardDesc,

			Infrastructures.longhorn.name:      controller.LonghornName,
			Infrastructures.longhorn.namespace: controller.LonghornNamespace,
			Infrastructures.longhorn.icon:      "",
			Infrastructures.longhorn.desc:      controller.LanghornDesc,

			Infrastructures.rancherVM.name:      controller.RancherVMName,
			Infrastructures.rancherVM.namespace: controller.RancherVMNamespace,
			Infrastructures.rancherVM.icon:      "",
			Infrastructures.rancherVM.desc:      controller.RancherVMDesc,
		},
	}

	cm, err := c.Clientset.CoreV1().ConfigMaps(controller.DashboardNamespace).Create(configMap)
	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return c.ConfigMapGet(controller.DashboardNamespace, ConfigMapName)
		}
		return nil, err
	}

	return cm, err
}

func (c *ClientGenerator) ConfigMapGet(ns, id string) (*v1.ConfigMap, error) {
	return c.Clientset.CoreV1().ConfigMaps(ns).Get(id, util.GetOptions)
}

func (c *ClientGenerator) ConfigMapList() (*v1.ConfigMapList, error) {
	return c.Clientset.CoreV1().ConfigMaps("").List(util.ListEverything)
}
